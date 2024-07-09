package renew

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
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

	paramsMap := extractQueryParams(r)

	if paramsMap["producerUid"] == "" {
		children, err := getNodeChildren(r)
		if len(children) != 0 && err == nil {
			paramsMap["producerUid"] = children
		}
	}

	queryBuilder := NewBigQueryQueryBuilder(lib.RenewPolicyViewCollection, "rp", nil)
	query, queryParams := queryBuilder.BuildQuery(paramsMap)
	if query == "" {
		log.Print("error generating query")
		return "", nil, errors.New("error generating query")
	}

	log.Printf("query: %s\nqueryParams: %+v", query, queryParams)

	policies, err := lib.QueryParametrizedRowsBigQuery[PolicyInfo](query, queryParams)
	if err != nil {
		log.Printf("error executing query: %s", err.Error())
		return "", nil, err
	}

	resp := &GetRenewPolicies{
		RenewPolicies: policies,
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), policies, err

}

func extractQueryParams(r *http.Request) map[string]string {
	inputParams := r.URL.Query()

	paramsMap := make(map[string]string)
	for key, values := range inputParams {
		paramsMap[key] = values[0]
	}
	return paramsMap
}

func getNodeChildren(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		return "", err
	}
	if authToken.IsNetworkNode {
		childrenList := make([]string, 0)
		children, err := network.GetNodeChildren(authToken.UserID)
		if err != nil {
			return "", err
		}
		for _, child := range children {
			childrenList = append(childrenList, child.NodeUid)
		}
		return strings.Join(childrenList, ", "), nil
	}
	return "", nil
}
