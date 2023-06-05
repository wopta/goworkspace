package question

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func PersonSurvey(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	const (
		rulesFileName = "person_survey.json"
	)

	log.Println("Person Survey")

	body, err := io.ReadAll(r.Body)
	lib.CheckError(err)

	fx := new(models.Fx)

	surveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
	}

	rulesFile := lib.GetRulesFile(rulesFileName)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, surveys, body, nil)

	ruleOutputJson, err := json.Marshal(ruleOutput)

	return string(ruleOutputJson), ruleOutput, err
}
