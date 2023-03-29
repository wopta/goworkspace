package rules

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	"io"
	"log"
	"net/http"
)

func Life(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy    models.Policy
		rulesFile []byte
		err       error
	)
	const (
		rulesFileName = "life.json"
	)

	log.Println("Life")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	quotingInputData := getRulesInputData(&policy, err, req)

	rulesFile = getRulesFile(rulesFile, rulesFileName)
	product, err := getLifeProduct()
	lib.CheckError(err)
	_, ruleOutput := rulesFromJson(rulesFile, product, quotingInputData, nil)

	productJson, product, err := prd.ReplaceDatesInProduct(ruleOutput.(models.Product), 69)
	lib.CheckError(err)

	return productJson, product, err
}

func getLifeProduct() (models.Product, error) {
	product, err := prd.GetName("life", "v2")
	return product, err
}
