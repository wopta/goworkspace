package sellable

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
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

	return Life(authToken.Role, policy)
}

func Life(role string, policy models.Policy) (string, *models.Product, error) {
	var (
		err error
	)
	const (
		rulesFileName = "life.json"
	)

	in, err := getInputData(policy)
	if err != nil {
		return "", &models.Product{}, err
	}
	rulesFile := lib.GetRulesFile(rulesFileName)
	product, err := prd.GetProduct(policy.Name, policy.ProductVersion, role)
	if err != nil {
		return "", &models.Product{}, err
	}

	fx := new(models.Fx)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	product = ruleOutput.(*models.Product)
	jsonOut, err := json.Marshal(product)

	return string(jsonOut), product, err
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
