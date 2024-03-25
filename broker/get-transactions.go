package broker

import (
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
	"log"
	"net/http"
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

	policyUid := r.Header.Get("policyUid")

	log.Printf("policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
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

	if authToken.Role != models.UserRoleAdmin && policy.ProducerUid != userUid && !network.IsParentOf(authToken.UserID, policy.ProducerUid) {
		log.Printf("[GetPolicyTransactionsFx] policy %s is not included in %s %s portfolio", policyUid, authToken.Role, userUid)
		return "", response, fmt.Errorf("%s %s unauthorized for policy %s", authToken.Role, userUid, policyUid)
	}

	transactions := transaction.GetPolicyTransactions(r.Header.Get("origin"), policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	log.Printf("[GetPolicyTransactionsFx] response: %s", string(jsonOut))

	return string(jsonOut), response, err
}
