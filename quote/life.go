package quote

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
	"io"
	"math"
	"modernc.org/mathutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	e := json.Unmarshal(req, &data)

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	res, e := Life(authToken.Role, data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}

func Life(role string, data models.Policy) (models.Policy, error) {
	var err error
	contractorAge, err := data.CalculateContractorAge()

	b := lib.GetFilesByEnv("quote/life_matrix.csv")
	df := lib.CsvToDataframe(b)
	var selectRow []string

	_, ruleProduct, err := sellable.Life(role, data)
	lib.CheckError(err)

	addDefaultGuarantees(data, *ruleProduct)

	switch role {
	case models.UserRoleAll, models.UserRoleCustomer:
		deathGuarantee, err := data.ExtractGuarantee("death")
		lib.CheckError(err)
		fmt.Println("[Life] setting sumInsuredLimitOfIndeminity")
		calculateSumInsuredLimitOfIndemnity(data.Assets, deathGuarantee.Value.SumInsuredLimitOfIndemnity)
		fmt.Println("[Life] setting guarantees duration")
		calculateGuaranteeDuration(data.Assets, contractorAge, deathGuarantee.Value.Duration.Year)
	case models.UserRoleAgent, models.UserRoleAgency:
		guaranteesMap := data.GuaranteesToMap()
		if guaranteesMap["death"].IsSelected {
			guaranteesMap["permanent-disability"].Value.SumInsuredLimitOfIndemnity = math.Max(
				guaranteesMap["permanent-disability"].Value.SumInsuredLimitOfIndemnity,
				guaranteesMap["death"].Value.SumInsuredLimitOfIndemnity)

			guaranteesMap["serious-ill"].Value.SumInsuredLimitOfIndemnity = math.Min(guaranteesMap["serious-ill"].
				Value.SumInsuredLimitOfIndemnity, guaranteesMap["death"].Value.SumInsuredLimitOfIndemnity*0.5)
		} else if guaranteesMap["permanent-disability"].IsSelected {
			guaranteesMap["serious-ill"].Value.SumInsuredLimitOfIndemnity = math.Min(guaranteesMap["serious-ill"].
				Value.SumInsuredLimitOfIndemnity, guaranteesMap["permanent-disability"].Value.
				SumInsuredLimitOfIndemnity*0.5)
		}

		guaranteesList := make([]models.Guarante, 0)
		for _, guarantee := range guaranteesMap {
			guaranteesList = append(guaranteesList, guarantee)
		}

		data.Assets[0].Guarantees = guaranteesList
	}

	updatePolicyStartEndDate(&data)

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

	for _, asset := range data.Assets {
		for _, guarantee := range asset.Guarantees {
			base, baseTax := getMultipliersIndex(guarantee.Slug)

			offset := getOffset(guarantee.Value.Duration.Year)

			baseFloat, taxFloat := getMultipliers(selectRow, offset, base, baseTax)

			calculateGuaranteePrices(guarantee, baseFloat, taxFloat, *ruleProduct)

			if guarantee.IsSelected && guarantee.IsSellable {
				calculateOfferPrices(data, guarantee)
			}
		}

	}

	monthlyToBeRemoved := !ruleProduct.Companies[0].IsMonthlyPaymentAvailable ||
		data.OffersPrices["default"]["monthly"].Gross < ruleProduct.Companies[0].MinimumMonthlyPrice
	if monthlyToBeRemoved {
		delete(data.OffersPrices["default"], "monthly")
	}

	roundOfferPrices(data.OffersPrices)

	sort.Slice(data.Assets[0].Guarantees, func(i, j int) bool {
		return data.Assets[0].Guarantees[i].Order < data.Assets[0].Guarantees[j].Order
	})

	return data, err
}

func addDefaultGuarantees(data models.Policy, product models.Product) {
	guaranteeList := make([]models.Guarante, 0)

	for _, guarantee := range data.Assets[0].Guarantees {
		product.Companies[0].GuaranteesMap[guarantee.Slug].Value = guarantee.Value
		product.Companies[0].GuaranteesMap[guarantee.Slug].IsSelected = guarantee.IsSelected
	}

	for _, guarantee := range product.Companies[0].GuaranteesMap {
		if guarantee.Value == nil {
			guarantee.Value = guarantee.Offer["default"]
		}
		guaranteeList = append(guaranteeList, *guarantee)
	}

	data.Assets[0].Guarantees = guaranteeList
}

func calculateSumInsuredLimitOfIndemnity(assets []models.Asset, deathSumInsuredLimitOfIndemnity float64) {
	for _, asset := range assets {
		for _, guarantee := range asset.Guarantees {
			switch guarantee.Slug {
			case "permanent-disability":
				guarantee.Value.SumInsuredLimitOfIndemnity = deathSumInsuredLimitOfIndemnity
			case "temporary-disability":
				guarantee.Value.SumInsuredLimitOfIndemnity = (deathSumInsuredLimitOfIndemnity / 100) * 1
			case "serious-ill":
				if deathSumInsuredLimitOfIndemnity > 100000 {
					guarantee.Value.SumInsuredLimitOfIndemnity = 10000
				} else {
					guarantee.Value.SumInsuredLimitOfIndemnity = 5000
				}
			}
		}
	}
}

func calculateGuaranteeDuration(assets []models.Asset, contractorAge int, deathDuration int) {
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			switch guarantee.Slug {
			case "permanent-disability":
				assets[assetIndex].Guarantees[guaranteeIndex].Value.Duration.Year = deathDuration
			case "temporary-disability", "serious-ill":
				assets[assetIndex].Guarantees[guaranteeIndex].Value.Duration.Year = mathutil.Min(deathDuration, 10)
			}
		}
	}
}

func updatePolicyStartEndDate(policy *models.Policy) {
	if policy.StartDate.IsZero() {
		policy.StartDate = time.Now().UTC()
	}
	maxDuration := 0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Value.Duration.Year > maxDuration {
			maxDuration = guarantee.Value.Duration.Year
		}
	}
	policy.EndDate = policy.StartDate.AddDate(maxDuration, 0, 0)
}

func getGuaranteeSubtitle(assets []models.Asset) {
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
	case "death":
		base = 1
		baseTax = 2
	case "permanent-disability":
		base = 3
		baseTax = 4
	case "temporary-disability":
		base = 5
		baseTax = 6
	case "serious-ill":
		base = 7
		baseTax = 8
	}
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
	return offset
}

func getMultipliers(selectRow []string, offset int, base int, baseTax int) (float64, float64) {
	baseFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base+offset], "%", "", 1), ",", ".", 1), 64)
	taxFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax+offset], "%", "", 1), ",", ".", 1), 64)
	baseFloat = baseFloat / 100
	taxFloat = taxFloat / 100
	return baseFloat, taxFloat
}

func calculateGuaranteePrices(guarantee models.Guarante, baseFloat, taxFloat float64, product models.Product) {
	if guarantee.Slug != "temporary-disability" {
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

	hasNotMinimumYearlyPrice := guarantee.Value.PremiumGrossYearly < product.Companies[0].GuaranteesMap[guarantee.Slug].Config.MinimumGrossYearly
	if hasNotMinimumYearlyPrice {
		guarantee.Value.PremiumGrossYearly = 10
		if guarantee.Slug == "death" {
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
	data.OffersPrices["default"]["yearly"].Gross += guarantee.Value.PremiumGrossYearly
	data.OffersPrices["default"]["yearly"].Net += guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["yearly"].Tax += guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["monthly"].Gross += guarantee.Value.PremiumGrossMonthly
	data.OffersPrices["default"]["monthly"].Net += guarantee.Value.PremiumNetMonthly
	data.OffersPrices["default"]["monthly"].Tax += guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly
}

func roundOfferPrices(offersPrices map[string]map[string]*models.Price) {
	for offerKey, offerValue := range offersPrices {
		for paymentKey, _ := range offerValue {
			offersPrices[offerKey][paymentKey].Net = lib.RoundFloat(offersPrices[offerKey][paymentKey].Net, 2)
			offersPrices[offerKey][paymentKey].Tax = lib.RoundFloat(offersPrices[offerKey][paymentKey].Tax, 2)
			offersPrices[offerKey][paymentKey].Gross = lib.RoundFloat(offersPrices[offerKey][paymentKey].Gross, 2)
		}
	}
}
