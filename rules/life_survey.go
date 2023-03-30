package rules

import (
	"encoding/json"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

func LifeSurvey(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		grule     []byte
		questions []*models.Statement
	)
	const (
		rulesFileName = "life_survey.json"
	)

	log.Println("Life Survey")

	statements := &Statements{
		Statements: questions,
	}

	rulesFile := getRulesFile(grule, rulesFileName)

	_, statementsOut := rulesFromJson(rulesFile, statements, nil, nil)

	statementsJson, err := json.Marshal(statementsOut)

	return string(statementsJson), statementsOut, err
}
