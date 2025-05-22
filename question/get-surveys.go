package question

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type Surveys struct {
	Surveys []*models.Survey
	Text    string
}

func GetSurveys(policy *models.Policy) ([]models.Survey, error) {
	log.AddPrefix("GetSurveys")
	defer log.PopPrefix()
	log.Println("function start ---------------------")

	policyJson, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling policy: %s", err.Error())
		return nil, err
	}

	fx := new(models.Fx)
	ruleSurveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
		Text:    "",
	}

	log.Println("loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, surveys)
	data := loadExternalData(policy.Name, policy.ProductVersion)

	log.Println("executing rules")

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, ruleSurveys, policyJson, data)

	result := make([]models.Survey, 0)
	for _, survey := range ruleOutput.(*Surveys).Surveys {
		result = append(result, *survey)
	}

	log.Println("function end ------------------")

	return result, err

}
