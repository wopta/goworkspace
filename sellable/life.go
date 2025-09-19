package sellable

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/network"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

// DEPRECATED
func lifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		warrant *models.Warrant
		err     error
	)
	log.AddPrefix("LifeFx")
	defer log.PopPrefix()

	log.Println("handler start ----------- ")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return "", nil, err
	}

	policy.Normalize()

	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	log.Println("loading network node")
	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("calling vendibility rules function")

	product, err := Life(policy, authToken.GetChannelByRoleV2(), networkNode, warrant)
	if err != nil {
		log.ErrorF("vednibility rules error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := product.Marshal()

	log.Printf("response: %s", string(jsonOut))
	log.Println("handler end -------------------")

	return string(jsonOut), product, err
}

func Life(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Product, error) {
	var (
		err     error
		product *models.Product
	)
	log.AddPrefix("Life")
	defer log.PopPrefix()

	log.Println("function start -----------")

	log.Println("loading input data")

	in, err := getInputData(*policy)
	if err != nil {
		log.ErrorF("error getting input data: %s", err.Error())
		return nil, err
	}

	log.Println("loading vendibility rules file")
	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)

	log.Println("loading product")
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		return nil, fmt.Errorf("no product found")
	}

	log.Println("applying vendibility rules")

	fx := new(models.Fx)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	product = ruleOutput.(*models.Product)

	log.Println("function end ----------")

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
