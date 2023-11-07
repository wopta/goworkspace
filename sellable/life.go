package sellable

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/network"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
)

// DEPRECATED
func LifeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		warrant *models.Warrant
		err     error
	)

	log.Println("[LifeFx] handler start ------------------------------------- ")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[LifeFx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return "", nil, err
	}

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	log.Println("[LifeFx] loading network node")
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("[LifeFx] calling vendibility rules function")

	product, err := Life(policy, authToken.GetChannelByRoleV2(), networkNode, warrant)
	if err != nil {
		log.Printf("[LifeFx] vednibility rules error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := product.Marshal()

	models.CreateAuditLog(r, string(body))

	log.Printf("[LifeFx] response: %s", string(jsonOut))
	log.Println("[LifeFx] handler end ----------------------------------------")

	return string(jsonOut), product, err
}

func Life(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) (*models.Product, error) {
	var (
		err     error
		product *models.Product
	)

	log.Println("[Life] function start -----------")

	log.Println("[Life] loading input data")

	in, err := getInputData(*policy)
	if err != nil {
		log.Printf("[Life] error getting input data: %s", err.Error())
		return nil, err
	}

	log.Println("[Life] loading vendibility rules file")
	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)

	log.Println("[Life] loading product")
	product = prd.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		return nil, fmt.Errorf("no product found")
	}

	log.Println("[Life] applying vendibility rules")

	fx := new(models.Fx)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, product, in, nil)

	product = ruleOutput.(*models.Product)

	log.Println("[Life] function end ----------")

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
