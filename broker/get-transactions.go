package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	plcUtils "gitlab.dev.wopta.it/goworkspace/policy/utils"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func getPolicyTransactionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var response GetPolicyTransactionsResp

	log.AddPrefix("GetPolicyTransactionsFx")
	defer log.PopPrefix()
	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "policyUid")

	log.Printf("policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		log.ErrorF("error getting authToken: %s", err.Error())
		return "", nil, err
	}

	policy, err := plc.GetPolicy(policyUid)
	if err != nil {
		log.ErrorF("error fetching policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}

	userUid := authToken.UserID

	if !plcUtils.CanBeAccessedBy(authToken.Role, policy.ProducerUid, authToken.UserID) {
		log.Printf("policy %s is not included in %s %s portfolio", policyUid, authToken.Role, userUid)
		return "", response, fmt.Errorf("%s %s unauthorized for policy %s", authToken.Role, userUid, policyUid)
	}

	transactions := transaction.GetPolicyTransactions(policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), response, err
}
