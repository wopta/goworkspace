package question

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

// DEPRECATED
func GetQuestionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		out    interface{}
		policy models.Policy
	)
	log.Println("[GetQuestionsFx]")

	questionType := r.Header.Get("questionType")
	log.Println("[GetQuestionFx] questionType " + questionType)

	body, err := io.ReadAll(r.Body)
	lib.CheckError(err)
	err = json.Unmarshal(body, &policy)

	switch questionType {
	case statements:
		log.Printf("[GetQuestionFx] loading statements for %s product", policy.Name)
		out = GetStatements(policy)
	case surveys:
		log.Printf("[GetQuestionFx] loading surveys for %s product", policy.Name)
		out = GetSurveys(policy)
	default:
		return "", nil, fmt.Errorf("questionType %s not allowed", questionType)
	}

	jsonOut, err := json.Marshal(out)

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, out, err
}

// DEPRECATED
func GetStatements(policy models.Policy) []models.Statement {
	const (
		rulesFilenameSuffix = "_statements.json"
	)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	fx := new(models.Fx)
	statementsValue := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	rulesFile := lib.GetRulesFile(policy.Name + rulesFilenameSuffix)
	data := loadExternalData(policy.Name)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, statementsValue, policyJson, data)

	out := make([]models.Statement, 0)
	for _, statement := range ruleOutput.(*Statements).Statements {
		out = append(out, *statement)
	}

	return out
}

// DEPRECATED
func GetSurveys(policy models.Policy) []models.Survey {
	const (
		rulesFilenameSuffix = "_surveys.json"
	)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	fx := new(models.Fx)
	surveysValue := &Surveys{
		Surveys: make([]*models.Survey, 0),
		Text:    "",
	}

	rulesFile := lib.GetRulesFile(policy.Name + rulesFilenameSuffix)
	data := loadExternalData(policy.Name)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, surveysValue, policyJson, data)

	out := make([]models.Survey, 0)
	for _, survey := range ruleOutput.(*Surveys).Surveys {
		out = append(out, *survey)
	}

	return out
}
