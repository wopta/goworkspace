package renew

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/network"
	md "github.com/wopta/goworkspace/policy/models"
	"github.com/wopta/goworkspace/policy/query-builder"
)

type GetRenewPoliciesResp struct {
	RenewPolicies []md.PolicyInfo `json:"policies"`
}

func GetRenewPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
	)

	log.SetPrefix("[GetRenewPoliciesFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	paramsMap := extractQueryParams(r)
	if len(paramsMap) == 0 {
		return "", nil, errors.New("no query params")
	}

	if paramsMap["producerUid"] == "" {
		children, err := getNodeChildren(r)
		if len(children) != 0 && err == nil {
			paramsMap["producerUid"] = children
		}
	}

	queryBuilder := query_builder.NewQueryBuilder("policy")
	query, queryParams := queryBuilder.BuildQuery(paramsMap)
	if query == "" {
		log.Print("error generating query")
		return "", nil, errors.New("error generating query")
	}

	log.Printf("query: %s\nqueryParams: %+v", query, queryParams)

	policies, err := lib.QueryParametrizedRowsBigQuery[md.PolicyInfo](query, queryParams)
	if err != nil {
		log.Printf("error executing query: %s", err.Error())
		return "", nil, err
	}

	resp := &GetRenewPoliciesResp{
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
