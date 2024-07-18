package enrich

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
)

func NaicsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		naics  []byte
		result []Naics
	)

	log.SetPrefix("[NaicsFx] ")
	defer log.SetPrefix("")

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
	b, err := json.Marshal(result)
	lib.CheckError(err)

	log.Println("Handler end -------------------------------------------------")

	return "{\"naics\":" + string(b) + "}", nil, nil
}

type Naics struct {
	Category      string `json:"category"`
	Detail        string `json:"detail "`
	Code          string `json:"code "`
	IsQbeSellable bool   `json:"isQbeSellable"`
}
