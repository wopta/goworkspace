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

func Life(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		rulesFile []byte
		err       error
	)
	const (
		rulesFileName = "life.json"
	)

	log.Println("Life")

	fx := new(models.Fx)

	in, err := getInputData(lib.ErrorByte(io.ReadAll(r.Body)))
	lib.CheckError(err)

	rulesFile = lib.GetRulesFile(rulesFile, rulesFileName)
	product, err := prd.GetProduct("life", "v1")
	lib.CheckError(err)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	productJson, product, err := prd.ReplaceDatesInProduct(ruleOutput.(models.Product), 69)

	return productJson, product, nil
}

func getInputData(b []byte) ([]byte, error) {
	var policy models.Policy
	err := json.Unmarshal(b, &policy)
	if err != nil {
		return nil, err
	}

	age, err := policy.CalculateAge()
	if err != nil {
		return nil, err
	}

	out := make(map[string]int)
	out["age"] = age

	output, err := json.Marshal(out)

	return output, err
}
