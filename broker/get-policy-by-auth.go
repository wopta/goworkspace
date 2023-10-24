package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

func GetPoliciesByAuthFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		req        GetPoliciesReq
		resp       GetPoliciesResp
		fieldName  = "producerUid"
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

	if req.Limit != 0 {
		limitValue = req.Limit
	}

	req.Queries = append(req.Queries, models.Query{
		Field: fieldName,
		Op:    "==",
		Value: authToken.UserID,
	})

	resp.Policies, err = plc.GetPoliciesByQueries(origin, req.Queries, limitValue)
	if err != nil {
		log.Println("[GetPoliciesByAuth] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPoliciesByAuth]: found %d policies", len(resp.Policies))

	jsonOut, err := json.Marshal(resp)
	return string(jsonOut), resp, err
}
