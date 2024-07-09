package renew

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
)

type PolicyInfo struct {
	Uid            string         `json:"uid" bigquery:"uid"`
	ProductName    string         `json:"productName" bigquery:"productName"`
	CodeCompany    string         `json:"codeCompany" bigquery:"codeCompany"`
	ProposalNumber int            `json:"proposalNumber" bigquery:"proposalNumber"`
	NameDesc       string         `json:"nameDesc" bigquery:"nameDesc"`
	Status         string         `json:"status" bigquery:"status"`
	Contractor     string         `json:"contractor" bigquery:"contractor"`
	Price          float64        `json:"price" bigquery:"price"`
	PriceMonthly   float64        `json:"priceMonthly" bigquery:"priceMonthly"`
	Producer       string         `json:"producer" bigquery:"producer"`
	ProducerCode   string         `json:"producerCode" bigquery:"producerCode"`
	StartDate      civil.DateTime `json:"startDate" bigquery:"startDate"`
	EndDate        civil.DateTime `json:"endDate" bigquery:"endDate"`
	PaymentSplit   string         `json:"paymentSplit" bigquery:"paymentSplit"`
}

type GetRenewPolicies struct {
	RenewPolicies []PolicyInfo `json:"renewPolicies"`
}

func GetRenewPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	inputParams := r.URL.Query()

	paramsMap := make(map[string]string)
	for key, values := range inputParams {
		paramsMap[key] = values[0]
	}

	queryBuilder := NewBigQueryQueryBuilder(lib.RenewPolicyViewCollection, "rp", nil)
	query, queryParams := queryBuilder.BuildQuery(paramsMap)

	policies, err := lib.QueryParametrizedRowsBigQuery[PolicyInfo](query, queryParams)
	if err != nil {

		return "", nil, err
	}

	resp := &GetRenewPolicies{
		RenewPolicies: policies,
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), policies, err

}
