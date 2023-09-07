package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
)

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func GetPolicyTransactionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetPolicyTransactionsFx] Handler start ---------------------")

	var response GetPolicyTransactionsResp

	policyUid := r.Header.Get("policyUid")

	log.Printf("[GetPolicyTransactionsFx] policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	userUid := authToken.UserID

	switch authToken.Role {
	case models.UserRoleAgent:
		if !models.IsPolicyInAgentPortfolio(userUid, policyUid) {
			log.Printf("[GetPolicyTransactionsFx] policy %s is not included in agent %s", policyUid, userUid)
			return "", response, fmt.Errorf("agent %s unauthorized for policy %s", userUid, policyUid)
		}
	case models.UserRoleAgency:
		if !models.IsPolicyInAgencyPortfolio(userUid, policyUid) {
			log.Printf("[GetPolicyTransactionsFx] policy %s is not included in agency %s", policyUid, userUid)
			return "", response, fmt.Errorf("agency %s unauthorized for policy %s", userUid, policyUid)
		}
	}

	transactions := transaction.GetPolicyTransactions(r.Header.Get("origin"), policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	return string(jsonOut), response, err
}
