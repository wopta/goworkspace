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

	lib "github.com/wopta/goworkspace/lib"
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
	birthDate, err := time.Parse("2006-01-02T15:04:05Z", data.Contractor.BirthDate)
	lib.CheckError(err)
	year := time.Now().Year() - birthDate.Year()

	b := lib.GetFilesByEnv("quote/life_matrix.csv")
	df := lib.CsvToDataframe(b)
	var selectRow []string

	//TODO: this should not be here, only for version 1
	deathGuarantee, err := getDeathGuarantee(data.Assets)
	lib.CheckError(err)
	//TODO: this should not be here, only for version 1
	calculateSumInsuredLimitOfIndemnity(data.Assets, deathGuarantee.Value.SumInsuredLimitOfIndemnity)

	calculateGuaranteeDuration(data.Assets, deathGuarantee.Value.Duration.Year)

	getGuaranteeSubtitle(data.Assets)

	for _, row := range df.Records() {
		if row[0] == strconv.Itoa(year) {
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

			calculateGuaranteePrices(guarantee, baseFloat, taxFloat)

			calculateOfferPrices(data, guarantee)
		}

	}

	filterByMinimumPrice(data.Assets, data.OffersPrices)

	roundOfferPrices(data.OffersPrices)

	return data, err
}

func calculateGuaranteeDuration(assets []models.Asset, deathDuration int) {
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

func roundOfferPrices(offersPrices map[string]map[string]*models.Price) {
	for offerKey, offerValue := range offersPrices {
		for paymentKey, _ := range offerValue {
			offersPrices[offerKey][paymentKey].Net = lib.RoundFloatTwoDecimals(offersPrices[offerKey][paymentKey].Net)
			offersPrices[offerKey][paymentKey].Tax = lib.RoundFloatTwoDecimals(offersPrices[offerKey][paymentKey].Tax)
			offersPrices[offerKey][paymentKey].Gross = lib.RoundFloatTwoDecimals(offersPrices[offerKey][paymentKey].Gross)
		}
	}
}

func filterByMinimumPrice(assets []models.Asset, offersPrices map[string]map[string]*models.Price) {
	product, err := prd.GetProduct("life", "V1")
	lib.CheckError(err)

	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			hasNotMinimumYearlyPrice := guarantee.Value.PremiumGrossYearly < product.Companies[0].GuaranteesMap[guarantee.Slug].Config.MinimumGrossYearly
			if hasNotMinimumYearlyPrice {
				assets[assetIndex].Guarantees[guaranteeIndex].IsSellable = false
				offersPrices["default"]["monthly"].Net -= guarantee.Value.PremiumNetMonthly
				offersPrices["default"]["monthly"].Gross -= guarantee.Value.PremiumGrossMonthly
				offersPrices["default"]["monthly"].Tax -= guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly
				offersPrices["default"]["yearly"].Net -= guarantee.Value.PremiumNetYearly
				offersPrices["default"]["yearly"].Gross -= guarantee.Value.PremiumGrossYearly
				offersPrices["default"]["yearly"].Tax -= guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly
			}
		}
	}

	if offersPrices["default"]["monthly"].Gross < product.Companies[0].MinimumMonthlyPrice {
		delete(offersPrices["default"], "monthly")
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
		guarantee.Value.PremiumNetYearly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat / 12)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat / 12)
	} else {
		guarantee.Value.PremiumNetYearly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat * 12)
		guarantee.Value.PremiumGrossYearly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat * 12)

		guarantee.Value.PremiumNetMonthly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat)
		guarantee.Value.PremiumGrossMonthly = lib.RoundFloatTwoDecimals(guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat)
	}

	guarantee.Value.PremiumTaxAmountYearly = lib.RoundFloatTwoDecimals(guarantee.Value.PremiumGrossYearly - guarantee.Value.PremiumNetYearly)
	guarantee.Value.PremiumTaxAmountMonthly = lib.RoundFloatTwoDecimals(guarantee.Value.PremiumGrossMonthly - guarantee.Value.PremiumNetMonthly)
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

func getDeathGuarantee(assets []models.Asset) (models.Guarante, error) {
	for _, asset := range assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == "death" {
				return guarantee, nil
			}
		}
	}
	return models.Guarante{}, fmt.Errorf("death guarantee not found")
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
