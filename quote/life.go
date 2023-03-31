package quote

import (
	"encoding/json"
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
	e := json.Unmarshal([]byte(req), &data)
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

	deathSumInsuredLimitOfIndemnity := 0.0
	for _, asset := range data.Assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == "death" {
				deathSumInsuredLimitOfIndemnity = guarantee.SumInsuredLimitOfIndemnity
			}
		}
	}

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

			switch guarantee.Slug {
			case "death":
				base = 1
				baseTax = 2
			case "permanent-disability":
				base = 3
				baseTax = 4
				guarantee.SumInsuredLimitOfIndemnity = deathSumInsuredLimitOfIndemnity
			case "temporary-disability":
				base = 5
				baseTax = 6
				guarantee.SumInsuredLimitOfIndemnity = (deathSumInsuredLimitOfIndemnity / 100) * 1
			case "serious-ill":
				base = 7
				baseTax = 8
				if deathSumInsuredLimitOfIndemnity > 100000 {
					guarantee.SumInsuredLimitOfIndemnity = 10000
				} else {
					guarantee.SumInsuredLimitOfIndemnity = 5000
				}
			}

			switch guarantee.Value.Duration.Year {
			case 5:
				offset = 8 * 0
			case 10:
				offset = 8 * 1
			case 15:
				offset = 8 * 2
			case 20:
				offset = 8 * 3
			}

			baseFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base+offset], "%", "", 1), ",", ".", 1), 64)
			taxFloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax+offset], "%", "", 1), ",", ".", 1), 64)
			baseFloat = baseFloat / 100
			taxFloat = taxFloat / 100

			guarantee.Offer["default"].PremiumNetYearly = guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat
			guarantee.Offer["default"].PremiumGrossYearly = guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat
			guarantee.Offer["default"].PremiumNetMonthly = guarantee.Value.SumInsuredLimitOfIndemnity * baseFloat / 12
			guarantee.Offer["default"].PremiumGrossMonthly = guarantee.Value.SumInsuredLimitOfIndemnity * taxFloat / 12

			data.OffersPrices["default"]["yearly"].Gross += guarantee.Offer["default"].PremiumGrossYearly
			data.OffersPrices["default"]["yearly"].Net += guarantee.Offer["default"].PremiumNetYearly
			data.OffersPrices["default"]["monthly"].Gross += guarantee.Offer["default"].PremiumGrossMonthly
			data.OffersPrices["default"]["monthly"].Net += guarantee.Offer["default"].PremiumNetMonthly
		}

	}

	return data, err
}
