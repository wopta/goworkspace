package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	plcUtils "github.com/wopta/goworkspace/policy/utils"
	"github.com/wopta/goworkspace/transaction"
)

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func GetPolicyTransactionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetPolicyTransactionsResp

	log.SetPrefix("[GetPolicyTransactionsFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "policyUid")

	log.Printf("policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.Printf("error getting authToken: %s", err.Error())
		return "", nil, err
	}

	policy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("error fetching policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}

	userUid := authToken.UserID

	if !plcUtils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, authToken.UserID) {
		log.Printf("policy %s is not included in %s %s portfolio", policyUid, authToken.Role, userUid)
		return "", response, fmt.Errorf("%s %s unauthorized for policy %s", authToken.Role, userUid, policyUid)
	}

	transactions := transaction.GetPolicyTransactions(r.Header.Get("Origin"), policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), response, err
}
