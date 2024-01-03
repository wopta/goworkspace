package policy

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
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
		limitValue = 10
	)
	log.Println("[GetPoliciesByQueryFx]")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	response.Policies, err = GetPoliciesByQueriesBigQuery(models.WoptaDataset, models.PoliciesViewCollection, req.Queries, limitValue)
	if err != nil {
		log.Println("[GetPoliciesByQueryFx] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPoliciesByQueryFx]: found %d policies", len(response.Policies))

	jsonOut, err := json.Marshal(response)
	return string(jsonOut), response, err
}
