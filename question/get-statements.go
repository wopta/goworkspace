package question

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
)

type Statements struct {
	Statements []*models.Statement
	Text       string
}

func GetStatementsV2(policy *models.Policy) ([]models.Statement, error) {
	log.Println("[GetStatementsV2] function start ----------------")

	policyJson, err := policy.Marshal()
	if err != nil {
		log.Printf("[GetStatementsV2] error marshaling policy: %s", err.Error())
		return nil, err
	}

	fx := new(models.Fx)
	rulesStatements := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	log.Println("[GetStatementsV2] loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, statements)
	data := loadExternalData(policy.Name)

	log.Println("[GetStatementsV2] executing rules")

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, rulesStatements, policyJson, data)

	result := make([]models.Statement, 0)
	for _, statement := range ruleOutput.(*Statements).Statements {
		result = append(result, *statement)
	}

	log.Println("[GetStatementsV2] function end ----------------")

	return result, err
}
