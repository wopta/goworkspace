package question

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func GapStatementsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
	)

	log.Println("GAP Statements")

	b, err := io.ReadAll(r.Body)
	lib.CheckError(err)
	err = json.Unmarshal(b, &policy)

	statements := GapStatements(policy)

	jsonOut, err := json.Marshal(statements)

	return string(jsonOut), statements, err
}

func GapStatements(policy models.Policy) []models.Statement {
	const (
		rulesFilename = "gap_statements.json"
	)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	fx := new(models.Fx)

	statements := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	rulesFile := lib.GetRulesFile(rulesFilename)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, statements, policyJson, nil)

	st := make([]models.Statement, 0)
	for _, statement := range ruleOutput.(*Statements).Statements {
		st = append(st, *statement)
	}

	return st
}
