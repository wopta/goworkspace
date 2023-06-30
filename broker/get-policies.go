package broker

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"time"
)

type GetPoliciesReq struct {
	Queries []models.Query `json:"queries,omitempty"`
	Limit   int            `json:"limit"`
	Page    int            `json:"page"`
}

type GetPoliciesResp struct {
	Policies []models.Policy `json:"policies"`
}

func GetPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetPoliciesReq
		response   GetPoliciesResp
		limitValue = 10
	)
	log.Println("GetPolicies")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")

	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0),
	}

	for index, q := range req.Queries {
		log.Printf("query %d/%d field: \"%s\" op: \"%s\" value: \"%v\"", index+1, len(req.Queries), q.Field, q.Op, q.Value)
		value := q.Value
		if q.Type == "dateTime" {
			value, _ = time.Parse(time.RFC3339, value.(string))
		}
		fireQueries.Queries = append(fireQueries.Queries, lib.Firequery{
			Field:      q.Field,
			Operator:   q.Op,
			QueryValue: value,
		})
	}

	docsnap, err := fireQueries.FirestoreWhereLimitFields(firePolicy, limitValue)
	if err != nil {
		return "", nil, err
	}

	response.Policies = models.PolicyToListData(docsnap)
	log.Printf("GetPolicies: found %d policies", len(response.Policies))

	jsonOut, err := json.Marshal(response)
	return string(jsonOut), response, err
}
