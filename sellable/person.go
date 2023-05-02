package sellable

/*



INPUT:
{
	"age": 74,
	"work": "select work list da elenco professioni",
	"workType": "dipendente" ,
	"coverageType": "24 ore" ,
	"childrenScool":true,
	"issue1500": 1,
	"riskInLifeIs":1,
	"class":2

}
{
	age:75
	work: "select work list da elenco professioni"
	worktype: dipendente / autonomo / non lavoratore
	coverageType 24 ore / tempo libero / professionale
	childrenScool:true
	issue1500:"si, senza problemi 1 || si, ma dovrei rinunciare a qualcosa 2|| no, non ci riscirei facilmente 3"
	riskInLifeIs:da evitare; 1 da accettare; 2 da gestire 3
	class

}
*/
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

func Person(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy    models.Policy
		rulesFile []byte
		e         error
	)
	const (
		rulesFileName = "person.json"
	)

	log.Println("Person")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	quotingInputData := getRulesInputData(&policy, e, req)

	fx := new(models.Fx)

	rulesFile = lib.GetRulesFile(rulesFile, rulesFileName)
	_, ruleOut := lib.RulesFromJsonV2(fx, rulesFile, initRuleOut(), quotingInputData, []byte(getQuotingData()))

	ruleOut.(*RuleOut).ToPolicy(&policy)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	return string(policyJson), policy, nil
}

func getRulesInputData(policy *models.Policy, e error, req []byte) []byte {
	*policy, e = models.UnmarshalPolicy(req)
	lib.CheckError(e)

	age, e := policy.CalculateAge()
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
