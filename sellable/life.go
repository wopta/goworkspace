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
		policy    models.Policy
		rulesFile []byte
		err       error
	)
	const (
		rulesFileName = "life.json"
	)

	log.Println("Life")

	fx := new(models.Fx)

	err = json.Unmarshal(lib.ErrorByte(io.ReadAll(r.Body)), &policy)
	if err != nil {
		return "", nil, err
	}

	product, productJson, err := Life(err, policy, rulesFile, rulesFileName, fx)

	return productJson, product, err
}

func Life(err error, policy models.Policy, rulesFile []byte, rulesFileName string, fx *models.Fx) (models.Product, string, error) {
	in, err := getInputData(policy)
	if err != nil {
		return models.Product{}, "", err
	}
	rulesFile = lib.GetRulesFile(rulesFile, rulesFileName)
	product, err := prd.GetProduct("life", "v1")
	if err != nil {
		return models.Product{}, "", err
	}

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	productJson, product, err := prd.ReplaceDatesInProduct(ruleOutput.(models.Product), 69)
	return product, productJson, err
}

func getInputData(policy models.Policy) ([]byte, error) {
	age, err := policy.CalculateAge()
	if err != nil {
		return nil, err
	}

	out := make(map[string]int)
	out["age"] = age

	output, err := json.Marshal(out)

	return output, err
}
