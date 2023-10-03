package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
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

	policy, err := plc.GetPolicy(policyUid, origin)
	lib.CheckError(err)

	userUid := authToken.UserID

	isAgentOrAgency := strings.EqualFold(authToken.Role, models.UserRoleAgent) || strings.EqualFold(authToken.Role, models.UserRoleAgency)
	if isAgentOrAgency && policy.ProducerUid != userUid {
		log.Printf("[GetPolicyTransactionsFx] policy %s is not included in %s %s portfolio", policyUid, authToken.Role, userUid)
		return "", response, fmt.Errorf("%s %s unauthorized for policy %s", authToken.Role, userUid, policyUid)
	}

	transactions := transaction.GetPolicyTransactions(r.Header.Get("origin"), policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	log.Printf("[GetPolicyTransactionsFx] response: %s", string(jsonOut))

	return string(jsonOut), response, err
}
