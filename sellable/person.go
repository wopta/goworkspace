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
		policy models.Policy
		err    error
	)

	log.Println("Person Sellable")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	policy = Person(req)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	return string(policyJson), policy, nil
}

func Person(body []byte) models.Policy {
	var (
		policy models.Policy
		err    error
	)
	const (
		rulesFileName = "person.json"
	)

	quotingInputData := getRulesInputData(&policy, err, body)

	fx := new(models.Fx)

	rulesFile := lib.GetRulesFile(rulesFileName)
	_, ruleOut := lib.RulesFromJsonV2(fx, rulesFile, initRuleOut(), quotingInputData, []byte(getQuotingData()))

	ruleOut.(*RuleOut).ToPolicy(&policy)

	return policy
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

func getPersonProduct() (models.Product, error) {
	product, err := prd.GetName("persona", "v1")
	return product, err
}

type RuleOut struct {
	Guarantees map[string]*models.Guarante         `json:"guarantees"`
	OfferPrice map[string]map[string]*models.Price `json:"offerPrice"`
}

func (r *RuleOut) ToPolicy(policy *models.Policy) {
	policy.OffersPrices = r.OfferPrice
	guarantees := make([]models.Guarante, 0)
	for _, guarantee := range r.Guarantees {
		guarantees = append(guarantees, *guarantee)
	}
	policy.Assets[0].Guarantees = guarantees
}

func initRuleOut() *RuleOut {
	var guarantees = make(map[string]*models.Guarante)
	offerPrice := make(map[string]map[string]*models.Price)

	product, err := getPersonProduct()
	lib.CheckError(err)

	for guaranteeKey, guarantee := range product.Companies[0].GuaranteesMap {
		guarantees[guaranteeKey] = &models.Guarante{
			CompanyName:                guarantee.CompanyName,
			Slug:                       guarantee.Slug,
			Deductible:                 guarantee.Deductible,
			Tax:                        guarantee.Tax,
			SumInsuredLimitOfIndemnity: guarantee.SumInsuredLimitOfIndemnity,
			Offer:                      guarantee.Offer,
		}
	}

	for offerKey, _ := range product.Offers {
		offerPrice[offerKey] = map[string]*models.Price{
			monthly: {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
			yearly: {
				Net:      0.0,
				Tax:      0.0,
				Gross:    0.0,
				Delta:    0.0,
				Discount: 0.0,
			},
		}
	}

	return &RuleOut{
		Guarantees: guarantees,
		OfferPrice: offerPrice,
	}
}

func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
