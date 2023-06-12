package broker

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

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
		fireQueries.Queries = append(fireQueries.Queries, lib.Firequery{
			Field:      q.Field,
			Operator:   q.Op,
			QueryValue: q.Value,
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

type GetPoliciesResp struct {
	Policies []models.Policy `json:"policies"`
}

type GetPoliciesReq struct {
	Queries []struct {
		Field string      `json:"field"`
		Op    string      `json:"op"`
		Value interface{} `json:"value"`
	} `json:"queries,omitempty"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
