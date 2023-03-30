package rules

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
)

func Life(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		rulesFile []byte
		err       error
	)
	const (
		rulesFileName = "life.json"
	)

	log.Println("Life")
	policyJson := lib.ErrorByte(io.ReadAll(r.Body))

	rulesFile = getRulesFile(rulesFile, rulesFileName)
	product, err := prd.GetName("life", "v1")
	if err != nil {
		return "", nil, err
	}

	_, ruleOutput := rulesFromJson(rulesFile, product, policyJson, nil)

	productJson, product, err := prd.ReplaceDatesInProduct(ruleOutput.(models.Product), 69)

	return productJson, product, err
}

func getInputData(policy *models.Policy, e error, req []byte) []byte {
	*policy, e = models.UnmarshalPolicy(req)
	lib.CheckError(e)

	age, e := calculateAge(policy.Contractor.BirthDate)
	lib.CheckError(e)
	tmpMap := make(map[string]int)

	tmpMap["age"] = age

	request, e := json.Marshal(tmpMap)
	lib.CheckError(e)
	return request
}
