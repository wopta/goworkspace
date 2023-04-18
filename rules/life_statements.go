package rules

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func LifeStatements(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy    models.Policy
		grule     []byte
		rulesFile []byte
	)
	const (
		rulesFilename = "life_statements.json"
	)

	log.Println("Life Statements")

	b, err := io.ReadAll(r.Body)
	lib.CheckError(err)
	err = json.Unmarshal(b, &policy)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	statements := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	rulesFile = getRulesFile(grule, rulesFilename)

	_, ruleOutput := rulesFromJson(rulesFile, statements, policyJson, nil)

	ruleOutputJson, err := json.Marshal(ruleOutput)
	lib.CheckError(err)

	return string(ruleOutputJson), ruleOutput, nil
}

type Statements struct {
	Statements []*models.Statement `json:"statements"`
	Text       string              `json:"text,omitempty"`
}
