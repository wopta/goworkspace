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

const (
	monthly = "monthly"
	yearly  = "yearly"
)

func PersonHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		product *models.Product
		err     error
	)

	log.Println("Person Sellable")

	origin := r.Header.Get("origin")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	product = Person(origin, req)

	productJson, err := json.Marshal(product)
	lib.CheckError(err)

	return string(productJson), product, nil
}

func Person(origin string, body []byte) *models.Product {
	var (
		policy models.Policy
		err    error
	)
	const (
		rulesFileName = "person.json"
	)

	quotingInputData := getRulesInputData(&policy, err, body)
	product, err := getPersonProduct(origin)
	lib.CheckError(err)

	fx := new(models.Fx)

	rulesFile := lib.GetRulesFile(rulesFileName)
	_, ruleOut := lib.RulesFromJsonV2(fx, rulesFile, &product, quotingInputData, []byte(getQuotingData()))

	return ruleOut.(*models.Product)
}

func getRulesInputData(policy *models.Policy, e error, req []byte) []byte {
	*policy, e = models.UnmarshalPolicy(req)
	lib.CheckError(e)

	age, e := policy.CalculateContractorAge()
	lib.CheckError(e)
	policy.QuoteQuestions["age"] = age
	policy.QuoteQuestions["work"] = policy.Contractor.Work
	policy.QuoteQuestions["workType"] = policy.Contractor.WorkType
	policy.QuoteQuestions["class"] = policy.Contractor.RiskClass

	request, e := json.Marshal(policy.QuoteQuestions)
	lib.CheckError(e)
	return request
}

func getPersonProduct(origin string) (models.Product, error) {
	product, err := prd.GetName(origin, "persona", "v1")
	return product, err
}

func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
