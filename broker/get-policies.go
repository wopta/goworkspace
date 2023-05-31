package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func GetPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req GetPoliciesReq
	)
	log.Println("GetPolicies")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")

	fireQueries := lib.Firequeries{
		Queries: make([]lib.Firequery, 0),
	}

	for _, q := range req.Queries {
		fireQueries.Queries = append(fireQueries.Queries, lib.Firequery{
			Field:      q.Field,
			Operator:   q.Op,
			QueryValue: q.Value,
		})
	}

	docsnap, err := fireQueries.FirestoreWhereLimitFields(firePolicy, req.Limit)
	if err != nil {
		return "", nil, err
	}

	policies := models.PolicyToListData(docsnap)

	jsonOut, err := json.Marshal(policies)
	return string(jsonOut), policies, err
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
