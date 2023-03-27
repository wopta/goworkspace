package quote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
		if row[0] == string(year) {
			selectRow = row

		}
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
			basefloat, _ := strconv.ParseFloat(selectRow[base], 64)
			taxfloat, _ := strconv.ParseFloat(selectRow[base], 64)
			guarance.PriceNett = guarance.Value.LimitOfIndemnity * basefloat
			guarance.PriceGross = guarance.Value.LimitOfIndemnity * taxfloat

		}

	}

	return data, e
}
