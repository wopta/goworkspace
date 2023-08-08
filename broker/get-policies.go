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
		resp       GetPoliciesResp
		limitValue = 10
	)
	log.Println("[GetPolicies]")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	resp.Policies, err = getPolicies(origin, req.Queries, limitValue)
	if err != nil {
		log.Println("[GetPolicies] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPolicies]: found %d policies", len(resp.Policies))

	jsonOut, err := json.Marshal(resp)
	return string(jsonOut), resp, err
}

func getPolicies(origin string, queries []models.Query, limitValue int) ([]models.Policy, error) {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0),
	}

	for index, q := range queries {
		log.Printf("query %d/%d field: \"%s\" op: \"%s\" value: \"%v\"", index+1, len(queries), q.Field, q.Op, q.Value)
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

	docSnap, err := fireQueries.FirestoreWhereLimitFields(firePolicy, limitValue)
	return models.PolicyToListData(docSnap), err
}
