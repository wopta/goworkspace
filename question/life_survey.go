package question

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/rules"
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

	fx := new(models.Fx)

	surveys := &rules.Surveys{
		Surveys: make([]*models.Survey, 0),
	}

	rulesFile := lib.GetRulesFile(grule, rulesFileName)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, surveys, nil, nil)

	ruleOutputJson, err := json.Marshal(ruleOutput)
	lib.CheckError(err)

	return string(ruleOutputJson), ruleOutput, nil
}
