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

func GetSurveysV2(policy *models.Policy) ([]models.Survey, error) {
	log.Println("[GetSurveysV2] function start ---------------------")

	policyJson, err := policy.Marshal()
	if err != nil {
		log.Printf("[GetSurveysV2] error marshaling policy: %s", err.Error())
		return nil, err
	}

	fx := new(models.Fx)
	ruleSurveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
		Text:    "",
	}

	log.Println("[GetSurveysV2] loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, surveys)
	data := loadExternalData(policy.Name)

	log.Println("[GetSurveysV2] executing rules")

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, ruleSurveys, policyJson, data)

	result := make([]models.Survey, 0)
	for _, survey := range ruleOutput.(*Surveys).Surveys {
		result = append(result, *survey)
	}

	log.Println("[GetSurveysV2] function end ------------------")

	return result, err

}
