package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
	qb "github.com/wopta/goworkspace/policy/query-builder/pkg"
)

type getPortfolioResp struct {
	Policies []policyInfo `json:"policies"`
}

type policyInfo struct {
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
	HasMandate     bool           `json:"hasMandate" bigquery:"hasMandate"`
	ContractorType string         `json:"contractorType" bigquery:"contractorType"`
}

func GetPortfolioFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var err error

	log.SetPrefix("[GetPortfolioFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	portfolioType := chi.URLParam(r, "type")

	paramsMap := extractQueryParams(r)
	if len(paramsMap) == 0 {
		return "", nil, errors.New("no query params")
	}

	log.Printf("input params: %v", paramsMap)

	if paramsMap["producerUid"] == "" {
		children, err := getNodeChildren(r)
		if err != nil {
			return "", nil, err
		}
		if len(children) != 0 {
			paramsMap["producerUid"] = children
		}
	}

	queryBuilder := qb.NewQueryBuilder(portfolioType)
	if queryBuilder == nil {
		return "", nil, errors.New("error initializing query builder")
	}
	query, queryParams, err := queryBuilder.Build(paramsMap)
	if err != nil {
		return "", nil, fmt.Errorf("error generating query: %v", err)
	}

	log.Printf("query: %s\nqueryParams: %+v", query, queryParams)

	policies, err := lib.QueryParametrizedRowsBigQuery[policyInfo](query, queryParams)
	if err != nil {
		log.Printf("error executing query: %s", err.Error())
		return "", nil, err
	}

	resp := &getPortfolioResp{
		Policies: policies,
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
