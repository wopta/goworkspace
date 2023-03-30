package quote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	e := json.Unmarshal([]byte(req), &data)
	res, e := Life(data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}
func Life(data models.Policy) (models.Policy, error) {
	var e error
	birthDate, e := time.Parse("2006-01-02T15:04:05Z", data.Contractor.BirthDate)
	lib.CheckError(e)
	year := time.Now().Year() - birthDate.Year()

	b := lib.GetFilesByEnv("quote/life_matrix.csv")
	df := lib.CsvToDataframe(b)
	var selectRow []string

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
		for _, guarance := range asset.Guarantees {
			var base int
			var baseTax int

			switch guarance.Slug {
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
			switch guarance.Value.Duration.Year {
			case 5:
				base = base * 1
				baseTax = baseTax * 1

			case 10:
				base = base * 2
				baseTax = baseTax * 2
			case 15:
				base = base * 3
				baseTax = baseTax * 3
			case 20:
				base = base * 4
				baseTax = baseTax * 4
			}
			basefloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[base], "%", "", 1), ",", ".", 1), 64)
			taxfloat, _ := strconv.ParseFloat(strings.Replace(strings.Replace(selectRow[baseTax], "%", "", 1), ",", ".", 1), 64)
			basefloat = basefloat / 100
			taxfloat = taxfloat / 100

			guarance.Offer["default"].PremiumNetYearly = guarance.Value.SumInsuredLimitOfIndemnity * basefloat
			guarance.Offer["default"].PremiumGrossYearly = guarance.Value.SumInsuredLimitOfIndemnity * taxfloat
			guarance.Offer["default"].PremiumNetMonthly = guarance.Value.SumInsuredLimitOfIndemnity * basefloat / 12
			guarance.Offer["default"].PremiumGrossMonthly = guarance.Value.SumInsuredLimitOfIndemnity * taxfloat / 12

			data.OffersPrices["default"]["yearly"].Gross = data.OffersPrices["default"]["yearly"].Gross + guarance.Offer["default"].PremiumGrossYearly
			data.OffersPrices["default"]["yearly"].Net = data.OffersPrices["default"]["yearly"].Net + guarance.Offer["default"].PremiumNetYearly
			data.OffersPrices["default"]["monthly"].Gross = data.OffersPrices["default"]["monthly"].Gross + guarance.Offer["default"].PremiumGrossMonthly
			data.OffersPrices["default"]["monthly"].Net = data.OffersPrices["default"]["monthly"].Net + guarance.Offer["default"].PremiumNetMonthly

		}

	}

	return data, e
}
