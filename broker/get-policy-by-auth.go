package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func GetPoliciesByAuthFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetPoliciesReq
		resp       GetPoliciesResp
		fieldName  string
		limitValue = 25
	)

	log.Println("[GetPoliciesByAuth]")

	origin := r.Header.Get("Origin")
	idToken := r.Header.Get("Authorization")

	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &req)
	lib.CheckError(err)

	if authToken.Role == models.UserRoleAgent {
		fieldName = "agentUid"
	} else if authToken.Role == models.UserRoleAgency {
		fieldName = "agencyUid"
	}

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	req.Queries = append(req.Queries, models.Query{
		Field: fieldName,
		Op:    "==",
		Value: authToken.UserID,
	})

	resp.Policies, err = getPolicies(origin, req.Queries, limitValue)
	if err != nil {
		log.Println("[GetPoliciesByAuth] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPoliciesByAuth]: found %d policies", len(resp.Policies))

	jsonOut, err := json.Marshal(resp)
	return string(jsonOut), resp, err
}
