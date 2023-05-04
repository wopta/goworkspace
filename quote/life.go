package quote

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"modernc.org/mathutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	e := json.Unmarshal(req, &data)
	res, e := Life(data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}
func Life(data models.Policy) (models.Policy, error) {
	var err error
	contractorAge, err := data.CalculateContractorAge()

	b := lib.GetFilesByEnv("quote/life_matrix.csv")
	df := lib.CsvToDataframe(b)
	var selectRow []string

	//TODO: this should not be here, only for version 1
	deathGuarantee, err := extractGuarantee(data.Assets[0].Guarantees, "death")
	lib.CheckError(err)
	//TODO: this should not be here, only for version 1
	calculateSumInsuredLimitOfIndemnity(data.Assets[0].Guarantees, deathGuarantee.Value.SumInsuredLimitOfIndemnity)

	calculateGuaranteeDuration(data.Assets, contractorAge, deathGuarantee.Value.Duration.Year)

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

	for _, guarantee := range data.Assets[0].Guarantees {
		base, baseTax := getMultipliersIndex(guarantee.Slug)

		offset := getOffset(guarantee.Value.Duration.Year)

		baseFloat, taxFloat := getMultipliers(selectRow, offset, base, baseTax)

		calculateGuaranteePrices(guarantee, baseFloat, taxFloat)

		calculateOfferPrices(data, guarantee)
	}

	filterByMinimumPrice(data)

	roundOfferPrices(data)

	return data, err
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
	policy.StartDate = time.Now().UTC()
	maxDuration := 0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Value.Duration.Year > maxDuration {
			maxDuration = guarantee.Value.Duration.Year
		}
	}
	policy.EndDate = policy.StartDate.AddDate(maxDuration, 0, 0)
}

func roundOfferPrices(policy models.Policy) {
	for offerKey, offerValue := range policy.OffersPrices {
		for paymentKey, _ := range offerValue {
			policy.OffersPrices[offerKey][paymentKey].Net = lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Net, 2)
			policy.OffersPrices[offerKey][paymentKey].Tax = lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Tax, 2)
			policy.OffersPrices[offerKey][paymentKey].Gross = lib.RoundFloat(policy.OffersPrices[offerKey][paymentKey].Gross, 2)
		}
	}
}

func filterByMinimumPrice(policy models.Policy) {
	product, err := prd.GetProduct("life", "V1")
	lib.CheckError(err)

	for guaranteeIndex, guarantee := range policy.Assets[0].Guarantees {
		hasNotMinimumYearlyPrice := guarantee.Value.PremiumGrossYearly < product.Companies[0].GuaranteesMap[guarantee.Slug].Config.MinimumGrossYearly
		if hasNotMinimumYearlyPrice {
			policy.Assets[0].Guarantees[guaranteeIndex].IsSellable = false
			policy.OffersPrices["default"]["monthly"].Net -= guarantee.Value.PremiumNetMonthly
			policy.OffersPrices["default"]["monthly"].Gross -= guarantee.Value.PremiumGrossMonthly
			policy.OffersPrices["default"]["monthly"].Tax -= guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly
			policy.OffersPrices["default"]["yearly"].Net -= guarantee.Value.PremiumNetYearly
			policy.OffersPrices["default"]["yearly"].Gross -= guarantee.Value.PremiumGrossYearly
			policy.OffersPrices["default"]["yearly"].Tax -= guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly
		}
	}

	if policy.OffersPrices["default"]["monthly"].Gross < product.Companies[0].MinimumMonthlyPrice {
		delete(policy.OffersPrices["default"], "monthly")
	}

}

func getGuaranteeSubtitle(assets []models.Asset) {
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			if guarantee.Slug != "death" {
				assets[assetIndex].Guarantees[guaranteeIndex].Subtitle = fmt.Sprintf("Durata: %d anni - Capitale: %s€",
					guarantee.Value.Duration.Year, humanize.FormatFloat("#.###,", guarantee.Value.SumInsuredLimitOfIndemnity))
			} else {
				assets[assetIndex].Guarantees[guaranteeIndex].Subtitle = "per qualsiasi causa"
			}
		}
	}
}

func calculateOfferPrices(data models.Policy, guarantee models.Guarante) {
	data.OffersPrices["default"]["yearly"].Gross += guarantee.Value.PremiumGrossYearly
	data.OffersPrices["default"]["yearly"].Net += guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["yearly"].Tax += guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["monthly"].Gross += guarantee.Value.PremiumGrossMonthly
	data.OffersPrices["default"]["monthly"].Net += guarantee.Value.PremiumNetMonthly
	data.OffersPrices["default"]["monthly"].Tax += guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly
}

func calculateGuaranteePrices(guarantee models.Guarante, baseFloat float64, taxFloat float64) {
	if guarantee.Slug != "temporary-disability" {
		guarantee.Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat, 2)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat, 2)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*lib.RoundFloat(baseFloat/12, 8), 2)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*lib.RoundFloat(taxFloat/12, 8), 2)
	} else {
		guarantee.Value.PremiumNetYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*baseFloat*12, 2)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*taxFloat*12, 2)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*lib.RoundFloat(baseFloat/12, 8), 2)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloat(guarantee.Value.SumInsuredLimitOfIndemnity*lib.RoundFloat(baseFloat/12, 8), 2)
	}

	guarantee.Value.PremiumTaxAmountYearly = lib.RoundFloat(guarantee.Value.PremiumGrossYearly-guarantee.Value.PremiumNetYearly, 2)
	guarantee.Value.PremiumTaxAmountMonthly = lib.RoundFloat(guarantee.Value.PremiumGrossMonthly-guarantee.Value.PremiumNetMonthly, 2)
}

func getMultipliers(selectRow []string, offset int, base int, baseTax int) (float64, float64) {
	baseFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base+offset], "%", "", 1), ",", ".", 1), 64)
	taxFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax+offset], "%", "", 1), ",", ".", 1), 64)
	baseFloat = baseFloat / 100
	taxFloat = taxFloat / 100
	return baseFloat, taxFloat
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

func extractGuarantee(guarantees []models.Guarante, guaranteeSlug string) (models.Guarante, error) {
	for _, guarantee := range guarantees {
		if guarantee.Slug == guaranteeSlug {
			return guarantee, nil
		}
	}
	return models.Guarante{}, fmt.Errorf("%s guarantee not found", guaranteeSlug)
}

func calculateSumInsuredLimitOfIndemnity(guarantees []models.Guarante, deathSumInsuredLimitOfIndemnity float64) {
	for _, guarantee := range guarantees {
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
