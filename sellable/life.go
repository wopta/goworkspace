package sellable

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
)

// DEPRECATED
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

	return Life(authToken.GetChannelByRole(), policy)
}

// DEPRECATED
func Life(channel string, policy models.Policy) (string, *models.Product, error) {
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
	product, err := prd.GetProduct(policy.Name, policy.ProductVersion, channel)
	if err != nil {
		return "", &models.Product{}, err
	}

	fx := new(models.Fx)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	product = ruleOutput.(*models.Product)
	jsonOut, err := json.Marshal(product)

	return string(jsonOut), product, err
}

func LifeV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy *models.Policy
		err    error
	)

	log.Println("[LifeV2Fx] handler start ----------- ")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[LifeV2Fx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return "", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	log.Println("[LifeV2Fx] calling vendibility rules function")

	product, err := LifeV2(policy, authToken.GetChannelByRoleV2())
	if err != nil {
		log.Printf("[LifeV2Fx] vednibility rules error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := product.Marshal()

	return string(jsonOut), product, err
}

func LifeV2(policy *models.Policy, channel string) (*models.Product, error) {
	var (
		err     error
		product *models.Product
	)
	const (
		sellableFileName = "sellable"
	)

	log.Println("[LifeV2] function start -----------")

	log.Println("[LifeV2] loading input data")

	in, err := getInputData(*policy)
	if err != nil {
		log.Printf("[LifeV2] error getting input data: %s", err.Error())
		return nil, err
	}

	log.Println("[LifeV2] loading vendibility rules file")
	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, sellableFileName)

	// TODO: controllare questo caso
	/*
		corretto prendere ultima versione attiva e non versione del prodotto all'interno della policy?
		potrebbero esserci problemi per la chiamata di quote all'emit (tentativo emissione di una proposta
		creata con versione vecchia di un prodotto)
	*/

	log.Println("[LifeV2] loading product")
	product = prd.GetDefaultProduct(policy.Name, channel)
	if product == nil {
		return nil, fmt.Errorf("no product found")
	}

	log.Println("[LifeV2] applying vendibility rules")

	fx := new(models.Fx)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	product = ruleOutput.(*models.Product)

	log.Println("[LifeV2] function end ----------")

	return product, nil
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
