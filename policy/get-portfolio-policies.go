package policy

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPortfolioPoliciesFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request  GetPoliciesReq
		response GetPoliciesResp
	)

	log.Println("[GetPortfolioPoliciesFx] Handler start ----------------------------------------")

	idToken := r.Header.Get("Authorization")

	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[GetPortfolioPoliciesFx] Request: %s", string(body))

	err = json.Unmarshal(body, &request)
	lib.CheckError(err)

	response.Policies, err = getPortfolioPolicies(authToken.UserID, request.Queries, request.Limit)
	if err != nil {
		log.Println("[GetPortfolioPoliciesFx] query error: ", err.Error())
		return "", nil, err
	}
	log.Printf("[GetPortfolioPoliciesFx]: found %d policies", len(response.Policies))

	jsonOut, err := json.Marshal(response)

	log.Println("[GetPortfolioPoliciesFx] Response: ", string(jsonOut))
	log.Println("[GetPortfolioPoliciesFx] Handler end ----------------------------------------")

	return string(jsonOut), response, err
}

func getPortfolioPolicies(producerUid string, requestQueries []models.Query, limit int) (policies []models.Policy, err error) {
	if len(requestQueries) == 0 {
		err = errors.New("no query specified")
		return
	}

	var (
		fieldName  = "producerUid"
		limitValue = 25
		queries    []models.Query
	)
	if limit != 0 {
		limitValue = limit
	}

	for _, q := range requestQueries {
		if q.Field == fieldName {
			log.Println("[getPortfolioPolicies] WARNING query with following field will be ignored: ", fieldName)
			continue
		} else {
			queries = append(queries, q)
		}
	}

	queries = append(queries, models.Query{
		Field: fieldName,
		Op:    "==",
		Value: producerUid,
	})

	return GetPoliciesByQueriesBigQuery(models.WoptaDataset, models.PoliciesViewCollection, queries, limitValue)
}
