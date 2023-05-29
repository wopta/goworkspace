package enrich

import (
	"encoding/json"
	"log"
	"net/http"

	lib "github.com/wopta/goworkspace/lib"
)

func Works(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		works  []byte
		result []Work
	)
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println("Work")
	w.Header().Set("Content-Type", "application/json")
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

		//log.Println(v)
		//log.Println(k)
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
				class = "x"
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
	//log.Println(string(b))
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
