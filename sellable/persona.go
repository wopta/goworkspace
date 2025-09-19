package sellable

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

// DEPRECATED
func personaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		product *models.Product
		warrant *models.Warrant
		err     error
	)
	log.AddPrefix("PersonaFx")
	defer log.PopPrefix()

	log.Println("handler start ----------------------")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf(" body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.ErrorF("error unmarshaling body: %s", err.Error())
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

	product = Persona(*policy, authToken.GetChannelByRoleV2(), networkNode, warrant)

	productJson, err := json.Marshal(product)
	lib.CheckError(err)

	log.Printf("response: %s", string(productJson))

	log.Println("handler end ------------------------")

	return string(productJson), product, nil
}

func Persona(policy models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) *models.Product {
	log.AddPrefix("Persona")
	defer log.PopPrefix()

	log.Println("function start -------------------------------------")

	log.Println("loading rules input data")

	quotingInputData := getPersonaRulesInputData(policy)

	log.Println("loading product file")

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, channel, networkNode, warrant)
	if product == nil {
		log.Printf("no product found")
		return nil
	}

	fx := new(models.Fx)

	log.Println("loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)
	_, ruleOut := lib.RulesFromJsonV2(fx, rulesFile, product, quotingInputData, []byte(getQuotingData()))

	log.Println("function end -------------------------------------")

	return ruleOut.(*models.Product)
}

func getPersonaRulesInputData(policy models.Policy) []byte {
	log.AddPrefix("getPersonaRulesInputData")
	defer log.PopPrefix()

	log.Println("function start ------------------")

	age, err := policy.CalculateContractorAge()
	if err != nil {
		log.ErrorF(" error getting contractor age: %s", err.Error())
		return nil
	}

	log.Printf("contractor age: %d", age)

	policy.QuoteQuestions["age"] = age
	policy.QuoteQuestions["work"] = policy.Contractor.Work
	policy.QuoteQuestions["workType"] = policy.Contractor.WorkType
	policy.QuoteQuestions["class"] = policy.Contractor.RiskClass

	result, err := json.Marshal(policy.QuoteQuestions)
	if err != nil {
		log.ErrorF("error marshaling policy quote questions: %s", err.Error())
		return nil
	}

	log.Printf("response: %s", result)

	log.Println("function end --------------------")

	return result
}

func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
