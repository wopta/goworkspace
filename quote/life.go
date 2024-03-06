package quote

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/network"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
	"modernc.org/mathutil"
)

const (
	deathGuarantee               = "death"
	permanentDisabilityGuarantee = "permanent-disability"
	temporaryDisabilityGuarantee = "temporary-disability"
	seriousIllGuarantee          = "serious-ill"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		data    models.Policy
		warrant *models.Warrant
	)

	log.Println("[LifeFx] handler start ----------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Println("[LifeFx] body: ", string(req))

	err := json.Unmarshal(req, &data)
	if err != nil {
		log.Printf("[LifeFx] error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	data.Normalize()

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("[LifeFx] error getting authToken from idToken: %s", err.Error())
		return "", nil, err
	}

	flow := authToken.GetChannelByRoleV2()

	log.Println("[LifeFx] loading network node")
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		if warrant != nil {
			flow = warrant.GetFlowName(data.Name)
		}
	}

	log.Println("[LifeFx] start quoting")

	result, err := Life(data, authToken.GetChannelByRoleV2(), networkNode, warrant, flow)
	jsonOut, err := json.Marshal(result)

	log.Printf("[LifeFx] response: %s", string(jsonOut))

	log.Println("[LifeFx] handler end ---------------------------------------")

	return string(jsonOut), result, err

}

func Life(data models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, flow string) (models.Policy, error) {
	var err error

	log.Println("[Life] function start --------------------------------------")

	contractorAge, err := data.CalculateContractorAge()

	log.Printf("[Life] contractor age: %d", contractorAge)

	b := lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/taxes.csv", models.ProductsFolder,
		data.Name, data.ProductVersion))
	df := lib.CsvToDataframe(b)
	var selectRow []string

	log.Printf("[Life] call sellable")
	ruleProduct, err := sellable.Life(&data, channel, networkNode, warrant)
	if err != nil {
		log.Printf("[Life] error in sellable: %s", err.Error())
		return models.Policy{}, err
	}

	log.Printf("[Life] loading available rates for flow %s", flow)

	availableRates := getAvailableRates(ruleProduct, flow)

	log.Printf("[Life] available rates: %s", availableRates)

	log.Printf("[Life] add default guarantees")

	addDefaultGuarantees(data, *ruleProduct)

	switch data.ProductVersion {
	case models.ProductV1:
		death, err := data.ExtractGuarantee(deathGuarantee)
		lib.CheckError(err)

		if channel == models.ECommerceChannel {
			log.Println("[Life] e-commerce flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnity(data.Assets, death.Value.SumInsuredLimitOfIndemnity)
			log.Println("[Life] setting guarantees duration")
			calculateGuaranteeDuration(data.Assets, death.Value.Duration.Year)
		}
	case models.ProductV2:
		if channel == models.ECommerceChannel {
			death, err := data.ExtractGuarantee(deathGuarantee)
			lib.CheckError(err)
			log.Println("[Life] e-commerce flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnity(data.Assets, death.Value.SumInsuredLimitOfIndemnity)
			log.Println("[Life] setting guarantees duration")
			calculateGuaranteeDuration(data.Assets, death.Value.Duration.Year)
		} else {
			log.Println("[Life] mga, network flow")
			log.Println("[Life] setting sumInsuredLimitOfIndeminity")
			calculateSumInsuredLimitOfIndemnityV2(&data)
		}
	}

	log.Println("[Life] updating policy start and end date")

	updatePolicyStartEndDate(&data)

	log.Println("[Life] set guarantees subtitle")

	getGuaranteeSubtitle(data.Assets)

	for _, row := range df.Records() {
		if row[0] == strconv.Itoa(contractorAge) {
			selectRow = row
			break
		}
	}

	data.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly":  &models.Price{},
			"monthly": &models.Price{},
		},
	}

	log.Println("[Life] calculate guarantees and offers prices")

	for assetIndex, asset := range data.Assets {
		for guaranteeIndex, _ := range asset.Guarantees {
			guarantee := &data.Assets[assetIndex].Guarantees[guaranteeIndex]
			base, baseTax := getMultipliersIndex(guarantee.Slug)

			offset := getOffset(guarantee.Value.Duration.Year)

			baseFloat, taxFloat := getMultipliers(selectRow, offset, base, baseTax)

			calculateGuaranteePrices(guarantee, baseFloat, taxFloat, *ruleProduct)

			if guarantee.IsSelected && guarantee.IsSellable {
				calculateOfferPrices(data, *guarantee)
			}
		}

	}

	log.Println("[Life] check monthly limit")

	monthlyToBeRemoved := !ruleProduct.Companies[0].IsMonthlyPaymentAvailable ||
		data.OffersPrices["default"]["monthly"].Gross < ruleProduct.Companies[0].MinimumMonthlyPrice
	if monthlyToBeRemoved {
		log.Println("[Life] monthly payment disabled")
		delete(data.OffersPrices["default"], "monthly")
	}

	log.Println("[Life] filtering available rates")

	removeOfferRate(&data, availableRates)

	log.Println("[Life] round offers prices")

	roundOfferPrices(data.OffersPrices)

	log.Println("[Life] sort guarantees list")

	sort.Slice(data.Assets[0].Guarantees, func(i, j int) bool {
		return data.Assets[0].Guarantees[i].Order < data.Assets[0].Guarantees[j].Order
	})

	log.Println("[Life] function end --------------------------------------")

	return data, err
}

func calculateSumInsuredLimitOfIndemnityV2(data *models.Policy) {
	guaranteesMap := data.GuaranteesToMap()

	log.Println("[Life] setting sumInsuredLimitOfIndeminity")
	if guaranteesMap[deathGuarantee].IsSelected {
		guaranteesMap[permanentDisabilityGuarantee].Value.SumInsuredLimitOfIndemnity =
			math.Max(guaranteesMap[permanentDisabilityGuarantee].Value.SumInsuredLimitOfIndemnity,
				guaranteesMap[deathGuarantee].Value.SumInsuredLimitOfIndemnity)

		minSumInsuredLimitOfIndemnity := math.Min(guaranteesMap[deathGuarantee].Value.SumInsuredLimitOfIndemnity,
			guaranteesMap[permanentDisabilityGuarantee].Value.SumInsuredLimitOfIndemnity)

		guaranteesMap[seriousIllGuarantee].Value.SumInsuredLimitOfIndemnity = math.Min(0.5*minSumInsuredLimitOfIndemnity,
			guaranteesMap[seriousIllGuarantee].Value.SumInsuredLimitOfIndemnity)

	} else if guaranteesMap[permanentDisabilityGuarantee].IsSelected {
		minSumInsuredLimitOfIndemnity := guaranteesMap[permanentDisabilityGuarantee].Value.SumInsuredLimitOfIndemnity

		guaranteesMap[seriousIllGuarantee].Value.SumInsuredLimitOfIndemnity = math.Min(0.5*minSumInsuredLimitOfIndemnity,
			guaranteesMap[seriousIllGuarantee].Value.SumInsuredLimitOfIndemnity)
	}

	guaranteesList := make([]models.Guarante, 0)
	for _, guarantee := range guaranteesMap {
		guarantee.Value.SumInsuredLimitOfIndemnity = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity, 0)
		guaranteesList = append(guaranteesList, guarantee)
	}

	data.Assets[0].Guarantees = guaranteesList
}

func addDefaultGuarantees(data models.Policy, product models.Product) {
	guaranteeList := make([]models.Guarante, 0)

	log.Println("[Life] adding default guarantees")

	for _, guarantee := range data.Assets[0].Guarantees {
		product.Companies[0].GuaranteesMap[guarantee.Slug].Value = guarantee.Value
		product.Companies[0].GuaranteesMap[guarantee.Slug].IsSelected = product.Companies[0].GuaranteesMap[guarantee.Slug].IsMandatory || guarantee.IsSelected
	}

	for _, guarantee := range product.Companies[0].GuaranteesMap {
		if guarantee.Value == nil {
			guarantee.Value = guarantee.Offer["default"]
			guarantee.IsSelected = guarantee.IsMandatory || getGuaranteeIsSelected(data, guarantee)
		}
		guaranteeList = append(guaranteeList, *guarantee)
	}

	data.Assets[0].Guarantees = guaranteeList
	log.Println("[Life] added default guarantees")
}

func getGuaranteeIsSelected(data models.Policy, guarantee *models.Guarante) bool {
	isSelected := false
	policyGuarantee, err := data.ExtractGuarantee(guarantee.Slug)
	if err == nil {
		isSelected = policyGuarantee.IsSelected
	}
	return isSelected
}

func calculateSumInsuredLimitOfIndemnity(assets []models.Asset, deathSumInsuredLimitOfIndemnity float64) {
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			switch guarantee.Slug {
			case permanentDisabilityGuarantee:
				assets[assetIndex].Guarantees[guaranteeIndex].Value.SumInsuredLimitOfIndemnity = deathSumInsuredLimitOfIndemnity
			case temporaryDisabilityGuarantee:
				assets[assetIndex].Guarantees[guaranteeIndex].Value.SumInsuredLimitOfIndemnity = (deathSumInsuredLimitOfIndemnity / 100) * 1
			case seriousIllGuarantee:
				if deathSumInsuredLimitOfIndemnity > 100000 {
					assets[assetIndex].Guarantees[guaranteeIndex].Value.SumInsuredLimitOfIndemnity = 10000
				} else {
					assets[assetIndex].Guarantees[guaranteeIndex].Value.SumInsuredLimitOfIndemnity = 5000
				}
			}
		}
	}
}

func calculateGuaranteeDuration(assets []models.Asset, deathDuration int) {
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			switch guarantee.Slug {
			case permanentDisabilityGuarantee:
				assets[assetIndex].Guarantees[guaranteeIndex].Value.Duration.Year = deathDuration
			case temporaryDisabilityGuarantee, seriousIllGuarantee:
				assets[assetIndex].Guarantees[guaranteeIndex].Value.Duration.Year = mathutil.Min(deathDuration, 10)
			}
		}
	}
}

func updatePolicyStartEndDate(policy *models.Policy) {
	if policy.StartDate.IsZero() {
		policy.StartDate = time.Now().UTC()
	}
	policy.StartDate = lib.SetDateToStartOfDay(policy.StartDate)
	maxDuration := 0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Value.Duration.Year > maxDuration {
			maxDuration = guarantee.Value.Duration.Year
		}
	}
	policy.EndDate = policy.StartDate.AddDate(maxDuration, 0, 0)
}

func getGuaranteeSubtitle(assets []models.Asset) {
	log.Println("[Life] setting guarantees subtitles")
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			assets[assetIndex].Guarantees[guaranteeIndex].Subtitle = fmt.Sprintf("Durata: %d anni - "+
				"Capitale: %s€", guarantee.Value.Duration.Year, humanize.FormatFloat("#.###,",
				guarantee.Value.SumInsuredLimitOfIndemnity))
		}
	}
}

func getMultipliersIndex(guaranteeSlug string) (int, int) {
	var base, baseTax int
	switch guaranteeSlug {
	case deathGuarantee:
		base = 1
		baseTax = 2
	case permanentDisabilityGuarantee:
		base = 3
		baseTax = 4
	case temporaryDisabilityGuarantee:
		base = 5
		baseTax = 6
	case seriousIllGuarantee:
		base = 7
		baseTax = 8
	}
	log.Printf("[Life] guarantee multipliers indexes base: %d baseTax: %d", base, baseTax)
	return base, baseTax
}

func getOffset(duration int) int {
	var offset int
	switch duration {
	case 5:
		offset = 8 * 1
	case 10:
		offset = 8 * 2
	case 15:
		offset = 8 * 3
	case 20:
		offset = 8*3 + 4
	}
	log.Printf("[Life] offset: %d", offset)
	return offset
}

func getMultipliers(selectRow []string, offset int, base int, baseTax int) (float64, float64) {
	baseFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base+offset], "%", "", 1), ",", ".", 1), 64)
	taxFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax+offset], "%", "", 1), ",", ".", 1), 64)
	baseFloat = baseFloat / 100
	taxFloat = taxFloat / 100
	log.Printf("[Life] guarantee multipliers baseFloat: %f taxFloat: %f", baseFloat, taxFloat)
	return baseFloat, taxFloat
}

func calculateGuaranteePrices(guarantee *models.Guarante, baseFloat, taxFloat float64, product models.Product) {
	if guarantee.Slug != temporaryDisabilityGuarantee {
		guarantee.Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat, 2)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat, 2)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat/12, 2)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat/12, 2)
	} else {
		guarantee.Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat*12, 2)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat*12, 2)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat, 2)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat, 2)
	}

	hasZeroYearlyPrice := guarantee.Value.PremiumGrossYearly == 0
	hasNotMinimumYearlyPrice := guarantee.Value.PremiumGrossYearly < product.Companies[0].GuaranteesMap[guarantee.Slug].Config.MinimumGrossYearly

	if hasZeroYearlyPrice {
		guarantee.IsSelected = false
		guarantee.IsSellable = false
		return
	} else if hasNotMinimumYearlyPrice {
		guarantee.Value.PremiumGrossYearly = 10
		if guarantee.Slug == deathGuarantee {
			guarantee.Value.PremiumNetYearly = 10
		} else {
			guarantee.Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly/(1+(2.5/100)), 2)
		}

		guarantee.Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly/12, 2)
		guarantee.Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.PremiumNetYearly/12, 2)
	}

	guarantee.Value.PremiumTaxAmountYearly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly-guarantee.Value.PremiumNetYearly, 2)
	guarantee.Value.PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Value.PremiumGrossMonthly-guarantee.Value.PremiumNetMonthly, 2)
}

func calculateOfferPrices(data models.Policy, guarantee models.Guarante) {
	log.Println("[Life] calculate offer prices")
	data.OffersPrices["default"]["yearly"].Gross += guarantee.Value.PremiumGrossYearly
	data.OffersPrices["default"]["yearly"].Net += guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["yearly"].Tax += guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["monthly"].Gross += guarantee.Value.PremiumGrossMonthly
	data.OffersPrices["default"]["monthly"].Net += guarantee.Value.PremiumNetMonthly
	data.OffersPrices["default"]["monthly"].Tax += guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly
}

func roundOfferPrices(offersPrices map[string]map[string]*models.Price) {
	log.Println("[Life] round offer prices")
	for offerKey, offerValue := range offersPrices {
		for paymentKey, _ := range offerValue {
			offersPrices[offerKey][paymentKey].Net = lib.RoundFloat(offersPrices[offerKey][paymentKey].Net, 2)
			offersPrices[offerKey][paymentKey].Tax = lib.RoundFloat(offersPrices[offerKey][paymentKey].Tax, 2)
			offersPrices[offerKey][paymentKey].Gross = lib.RoundFloat(offersPrices[offerKey][paymentKey].Gross, 2)
		}
	}
}
