package question

import (
	"encoding/json"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type Statements struct {
	Statements []*models.Statement
	Text       string
}
type inputRuleSettings struct {
	*models.Policy
	IncludeExternalCompanyStatements bool `json:"includeExternalCompanyStatements"`
}

func GetStatements(policy *models.Policy, includeCompanyStatements bool) ([]models.Statement, error) {

	log.AddPrefix("GetStatements")
	defer log.PopPrefix()
	log.Println("function start ----------------")
	inputRuleStruct := inputRuleSettings{
		Policy:                           policy,
		IncludeExternalCompanyStatements: includeCompanyStatements,
	}
	inputRuleStr, err := json.Marshal(inputRuleStruct)
	if err != nil {
		log.ErrorF("error marshaling policy: %s", err.Error())
		return nil, err
	}

	fx := new(models.Fx)
	rulesStatements := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	log.Println("loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, statements)
	data := loadExternalData(policy.Name, policy.ProductVersion)

	log.Println("executing rules")

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, rulesStatements, inputRuleStr, data)

	result := make([]models.Statement, 0)
	for _, statement := range ruleOutput.(*Statements).Statements {
		result = append(result, *statement)
	}

	log.Println("function end ----------------")

	return result, err
}
