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

// DEPRECATED
func PersonaFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		product *models.Product
		err     error
	)

	log.Println("[PersonaFx] handler start ----------------------")

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Printf("[PersonaFx] body: %s", string(body))

	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[PersonaFx] error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	product = Persona(*policy, authToken.GetChannelByRoleV2())

	productJson, err := json.Marshal(product)
	lib.CheckError(err)

	log.Printf("[PersonaFx] response: %s", string(productJson))

	log.Println("[PersonaFx] handler end ------------------------")

	return string(productJson), product, nil
}

func Persona(policy models.Policy, channel string) *models.Product {
	log.Println("[Persona] function start -------------------------------------")

	log.Println("[Persona] loading rules input data")

	quotingInputData := getPersonaRulesInputData(policy)

	log.Println("[Persona] loading product file")

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, channel)
	if product == nil {
		log.Printf("[Persona] no product found")
		return nil
	}

	fx := new(models.Fx)

	log.Println("[Persona] loading rules file")

	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, rulesFilename)
	_, ruleOut := lib.RulesFromJsonV2(fx, rulesFile, product, quotingInputData, []byte(getQuotingData()))

	log.Println("[Persona] function end -------------------------------------")

	return ruleOut.(*models.Product)
}

func getPersonaRulesInputData(policy models.Policy) []byte {
	log.Println("[getPersonaRulesInputData] function start ------------------")

	age, err := policy.CalculateContractorAge()
	if err != nil {
		log.Printf("[getPersonaRulesInputData] error getting contractor age: %s", err.Error())
		return nil
	}

	log.Printf("[getPersonaRulesInputData] contractor age: %d", age)

	policy.QuoteQuestions["age"] = age
	policy.QuoteQuestions["work"] = policy.Contractor.Work
	policy.QuoteQuestions["workType"] = policy.Contractor.WorkType
	policy.QuoteQuestions["class"] = policy.Contractor.RiskClass

	result, err := json.Marshal(policy.QuoteQuestions)
	if err != nil {
		log.Printf("[getPersonaRulesInputData] error marshaling policy quote questions: %s", err.Error())
		return nil
	}

	log.Printf("[getPersonaRulesInputData] response: %s", result)

	log.Println("[getPersonaRulesInputData] function end --------------------")

	return result
}

func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
