package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/civil"
	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/network"
	qb "gitlab.dev.wopta.it/goworkspace/policy/query-builder/pkg"
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
	Consultancy    float64        `json:"consultancy" bigquery:"consultancy"`
	Total          float64        `json:"total" bigquery:"total"`
}

func getPortfolioFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var err error

	log.AddPrefix("GetPortfolioFx")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	portfolioType := chi.URLParam(r, "type")

	paramsMap := extractQueryParams(r)
	if len(paramsMap) == 0 {
		return "", nil, errors.New("no query params")
	}

	log.Printf("input params: %v", paramsMap)

	if err := populateProducerUidParam(r, paramsMap); err != nil {
		log.ErrorF("error populating producerUid param")
		return "", nil, err
	}

	log.Printf("populated params: %v", paramsMap)

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
		log.ErrorF("error executing query: %s", err.Error())
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

func populateProducerUidParam(r *http.Request, paramsMap map[string]string) error {
	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		return err
	}

	if authToken.Role == lib.UserRoleAdmin && paramsMap["producerUid"] == "" {
		return nil
	}

	producerUid := paramsMap["producerUid"]
	if authToken.IsNetworkNode && producerUid == "" {
		producerUid = authToken.UserID
	}
	if authToken.IsNetworkNode && producerUid != "" && !network.IsParentOf(authToken.UserID, producerUid) {
		return fmt.Errorf("cannot access policy from node %s", producerUid)
	}

	childrenUids, err := getNodeChildren(producerUid)
	if err != nil {
		return err
	}
	paramsMap["producerUid"] = childrenUids

	return nil
}

func getNodeChildren(producerUid string) (string, error) {
	childrenMap := make(map[string]string)
	// We add the producer to the map because the partnership nodes are not
	// present in the GetNodeChildren table, returning an empty list.
	// The map is used to keep only unique uids for the query
	childrenMap[producerUid] = producerUid
	children, err := network.GetNodeChildren(producerUid)
	if err != nil {
		return "", err
	}
	for _, child := range children {
		childrenMap[child.NodeUid] = child.NodeUid
	}

	childrenList := make([]string, 0)
	for nodeUid := range childrenMap {
		childrenList = append(childrenList, nodeUid)
	}
	return strings.Join(childrenList, ", "), nil
}
