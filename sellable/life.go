package sellable

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
)

func LifeHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy models.Policy
		err    error
	)

	log.Println("Life")

	err = json.Unmarshal(lib.ErrorByte(io.ReadAll(r.Body)), &policy)
	if err != nil {
		return "", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	product, productJson, err := Life(authToken.Role, policy)

	return productJson, product, err
}

func Life(role string, policy models.Policy) (models.Product, string, error) {
	var (
		err error
	)
	const (
		rulesFileName = "life.json"
	)

	in, err := getInputData(policy)
	if err != nil {
		return models.Product{}, "", err
	}
	rulesFile := lib.GetRulesFile(rulesFileName)
	product, err := prd.GetProduct("life", "v1", role)
	if err != nil {
		return models.Product{}, "", err
	}

	fx := new(models.Fx)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	productJson, product, err := prd.ReplaceDatesInProduct(ruleOutput.(models.Product), 70, 55)
	return product, productJson, err
}

func getInputData(policy models.Policy) ([]byte, error) {
	age, err := policy.CalculateContractorAge()
	if err != nil {
		return nil, err
	}

	out := make(map[string]int)
	out["age"] = age

	output, err := json.Marshal(out)

	return output, err
}
