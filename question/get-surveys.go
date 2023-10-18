package question

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
)

type Surveys struct {
	Surveys []*models.Survey
	Text    string
}

func GetSurveys(policy *models.Policy) ([]models.Survey, error) {
	log.Println("[GetSurveys] function start ---------------------")

	policyJson, err := policy.Marshal()
	if err != nil {
		log.Printf("[GetSurveys] error marshaling policy: %s", err.Error())
		return nil, err
	}

	fx := new(models.Fx)
	ruleSurveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
		Text:    "",
	}

	log.Println("[GetSurveys] loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, surveys)
	data := loadExternalData(policy.Name, policy.ProductVersion)

	log.Println("[GetSurveys] executing rules")

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, ruleSurveys, policyJson, data)

	result := make([]models.Survey, 0)
	for _, survey := range ruleOutput.(*Surveys).Surveys {
		result = append(result, *survey)
	}

	log.Println("[GetSurveys] function end ------------------")

	return result, err

}
