package quote

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
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
	res, e := life(data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}
func life(data models.Policy) (models.Policy, error) {
	var err error
	birthDate, err := time.Parse("2006-01-02T15:04:05Z", data.Contractor.BirthDate)
	lib.CheckError(err)
	year := time.Now().Year() - birthDate.Year()

	b := lib.GetFilesByEnv("quote/life_matrix.csv")
	df := lib.CsvToDataframe(b)
	var selectRow []string

	//TODO: this should not be here, only for version 1
	deathSumInsuredLimitOfIndemnity := getDeathSumInsuredLimitOfIndemnity(data.Assets)
	//TODO: this should not be here, only for version 1
	calculateSumInsuredLimitOfIndemnity(data.Assets, deathSumInsuredLimitOfIndemnity)

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
			var base int
			var baseTax int
			var offset int

			base, baseTax = getMultipliersIndex(guarantee.Slug, base, baseTax)

			offset = getOffset(guarantee.Value.Duration.Year, offset)

			baseFloat, taxFloat := getMultipliers(selectRow, offset, base, baseTax)

			calculateGuaranteePrices(guarantee, baseFloat, taxFloat)

			calculateOfferPrices(data, guarantee)
		}

	}

	return data, err
}

func getGuaranteeSubtitle(assets []models.Asset) {
	for assetIndex, asset := range assets {
		for guaranteeIndex, guarantee := range asset.Guarantees {
			assets[assetIndex].Guarantees[guaranteeIndex].Subtitle = fmt.Sprintf("Durata: %d anni - Capitale: %s€",
				guarantee.Value.Duration.Year, humanize.FormatFloat("#.###,##", guarantee.Value.SumInsuredLimitOfIndemnity))
		}
	}
}

func calculateOfferPrices(data models.Policy, guarantee models.Guarante) {
	data.OffersPrices["default"]["yearly"].Gross += guarantee.Value.PremiumGrossYearly
	data.OffersPrices["default"]["yearly"].Net += guarantee.Value.PremiumNetYearly
	data.OffersPrices["default"]["monthly"].Gross += guarantee.Value.PremiumGrossMonthly
	data.OffersPrices["default"]["monthly"].Net += guarantee.Value.PremiumNetMonthly
}

func calculateGuaranteePrices(guarantee models.Guarante, baseFloat float64, taxFloat float64) {
	guarantee.Value.PremiumNetYearly = guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat
	guarantee.Value.PremiumGrossYearly = guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat
	guarantee.Value.PremiumNetMonthly = guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat / 12
	guarantee.Value.PremiumGrossMonthly = guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat / 12
}

func getMultipliers(selectRow []string, offset int, base int, baseTax int) (float64, float64) {
	baseFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base+offset], "%", "", 1), ",", ".", 1), 64)
	taxFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax+offset], "%", "", 1), ",", ".", 1), 64)
	baseFloat = baseFloat / 100
	taxFloat = taxFloat / 100
	return baseFloat, taxFloat
}

func getMultipliersIndex(guaranteeSlug string, base int, baseTax int) (int, int) {
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

func getOffset(duration int, offset int) int {
	switch duration {
	case 5:
		offset = 8 * 0
	case 10:
		offset = 8 * 1
	case 15:
		offset = 8 * 2
	case 20:
		offset = 8 * 3
	}
	return offset
}

func getDeathSumInsuredLimitOfIndemnity(assets []models.Asset) float64 {
	deathSumInsuredLimitOfIndemnity := 0.0
	for _, asset := range assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == "death" {
				deathSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
				break
			}
		}
	}
	return deathSumInsuredLimitOfIndemnity
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
