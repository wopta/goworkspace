package enrich

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
)

func naicsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		naics  []byte
		result []Naics
	)

	log.AddPrefix("NaicsFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	naics = lib.GetFilesByEnv("enrich/naics.csv")

	df := lib.CsvToDataframe(naics)

	for k, v := range df.Records() {
		var (
			isQbeSellable bool
		)

		if k > 0 {

			if v[3] == "OK" {

				isQbeSellable = true
			}

			sub := Naics{
				Category:      v[0],
				Detail:        v[1],
				Code:          v[2],
				IsQbeSellable: isQbeSellable,
			}
			result = append(result, sub)
		}
	}

	a := make(map[string][]Naics)
	for _, v := range result {
		a[v.Category] = append(a[v.Category], v)

	}
	b, err := json.Marshal(a)
	lib.CheckError(err)
	log.Println("Handler end -------------------------------------------------")

	return "{\"naicsMap\":" + string(b) + "}", nil, err
}

type Naics struct {
	Category      string `json:"category"`
	Detail        string `json:"detail"`
	Code          string `json:"code"`
	IsQbeSellable bool   `json:"isQbeSellable"`
}
