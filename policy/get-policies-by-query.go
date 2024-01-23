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
	Policies []models.Policy `json:"policies"`
}

func GetPoliciesByQueryFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetPoliciesReq
		response   GetPoliciesResp
	)
	log.Println("[GetPoliciesByQueryFx] Handler start ----------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[GetPoliciesByQueryFx] Request: %s", string(body))

	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	response.Policies, err = getPoliciesByQuery(req.Queries, req.Limit)
	if err != nil {
		log.Println("[GetPoliciesByQueryFx] query error: ", err.Error())
		return "", nil, errors.New("query error")
	}
	log.Printf("[GetPoliciesByQueryFx]: found %d policies", len(response.Policies))

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
