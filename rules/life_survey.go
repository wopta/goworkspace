package rules

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

func LifeSurvey(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		grule []byte
	)
	const (
		rulesFileName = "life_survey.json"
	)

	log.Println("Life Survey")

	surveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
	}

	rulesFile := getRulesFile(grule, rulesFileName)

	_, ruleOutput := rulesFromJson(rulesFile, surveys, nil, nil)

	ruleOutputJson, err := json.Marshal(ruleOutput)
	lib.CheckError(err)

	return string(ruleOutputJson), ruleOutput, nil
}
