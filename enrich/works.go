package enrich

import (
	"encoding/json"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
)

func worksFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		works  []byte
		result []Work
	)

	log.AddPrefix("WorksFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	works = lib.GetFilesByEnv("enrich/works.csv")

	df := lib.CsvToDataframe(works)

	for k, v := range df.Records() {
		var (
			selfEmployed bool
			employed     bool
			unemployed   bool
			workType     string
			class        string
		)

		if k > 0 {

			if v[2] == "x" {
				workType = "autonomo"
				selfEmployed = true
			}
			if v[3] == "x" {
				workType = "dipendente"
				employed = true
			}
			if v[4] == "x" {
				workType = "disoccupato"
				unemployed = true
			}

			class = v[1]
			if class == "S.E." {
				class = "1"
			}

			sub := Work{
				Work:           v[0],
				Class:          class,
				WorkType:       workType,
				IsSelfEmployed: selfEmployed,
				IsEmployed:     employed,
				IsUnemployed:   unemployed,
			}
			result = append(result, sub)
		}
	}
	b, err := json.Marshal(result)
	lib.CheckError(err)

	log.Println("Handler end -------------------------------------------------")

	return "{\"works\":" + string(b) + "}", nil, nil
}

type Work struct {
	Work           string `json:"work"`
	WorkType       string `json:"workType"`
	Class          string `json:"class"`
	IsSelfEmployed bool   `json:"isSelfEmployed"`
	IsEmployed     bool   `json:"isEmployed"`
	IsUnemployed   bool   `json:"isUnemployed"`
}
