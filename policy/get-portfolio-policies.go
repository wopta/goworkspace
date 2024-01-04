package policy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPortfolioPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request    GetPoliciesReq
		response   GetPoliciesResp
		fieldName  = "producerUid"
		limitValue = 25
	)

	log.Println("[GetPortfolioPoliciesFx]")

	idToken := r.Header.Get("Authorization")

	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &request)
	lib.CheckError(err)

	if request.Limit != 0 {
		limitValue = request.Limit
	}

	for _, q := range request.Queries {
		if q.Field == fieldName {
			return "", nil, fmt.Errorf("field name is not allowed: %s", fieldName)
		}
	}

	request.Queries = append(request.Queries, models.Query{
		Field: fieldName,
		Op:    "==",
		Value: authToken.UserID,
	})

	response.Policies, err = GetPoliciesByQueriesBigQuery(models.WoptaDataset, models.PoliciesViewCollection, request.Queries, limitValue)
	if err != nil {
		log.Println("[GetPortfolioPoliciesFx] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPortfolioPoliciesFx]: found %d policies", len(response.Policies))

	jsonOut, err := json.Marshal(response)
	return string(jsonOut), response, err
}
