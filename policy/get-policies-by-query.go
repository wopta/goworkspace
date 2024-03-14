package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type GetPoliciesReq struct {
	Queries []models.Query `json:"queries,omitempty"`
	Limit   int            `json:"limit"`
	Page    int            `json:"page"`
}

type GetPoliciesResp struct {
	Policies []PolicyInfo `json:"policies"`
}

func GetPoliciesByQueryFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req      GetPoliciesReq
		response GetPoliciesResp
		policies []models.Policy
	)
	log.Println("[GetPoliciesByQueryFx] Handler start ----------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[GetPoliciesByQueryFx] Request: %s", string(body))

	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	policies, err = getPoliciesByQuery(req.Queries, req.Limit)
	if err != nil {
		log.Println("[GetPoliciesByQueryFx] query error: ", err.Error())
		return "", nil, errors.New("query error")
	}
	log.Printf("[GetPoliciesByQueryFx]: found %d policies", len(response.Policies))

	for _, p := range policies {
		response.Policies = append(response.Policies, policyToPolicyInfo(p, ""))
	}

	response.Policies = lib.SliceMap(policies, func(p models.Policy) PolicyInfo {
		return policyToPolicyInfo(p, "")
	})

	jsonOut, err := json.Marshal(response)

	log.Println("[GetPoliciesByQueryFx] Response: ", string(jsonOut))
	log.Println("[GetPoliciesByQueryFx] Handler end ----------------------------------------")
	return string(jsonOut), response, err
}

func getPoliciesByQuery(queries []models.Query, limit int) (policies []models.Policy, err error) {
	if len(queries) == 0 {
		err = fmt.Errorf("no query specified")
		return
	}

	var limitValue = 10
	if limit != 0 {
		limitValue = limit
	}

	return GetPoliciesByQueriesBigQuery(models.WoptaDataset, models.PoliciesViewCollection, queries, limitValue)
}
