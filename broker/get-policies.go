package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
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
	log.SetPrefix("[GetPoliciesFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err := json.Unmarshal(body, &req)
	lib.CheckError(err)

	if len(req.Queries) == 0 {
		return "", nil, fmt.Errorf("no query specified")
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	resp.Policies, err = plc.GetPoliciesByQueries(origin, req.Queries, limitValue)
	if err != nil {
		log.Println("query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("found %d policies", len(resp.Policies))

	jsonOut, err := json.Marshal(resp)

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), resp, err
}
