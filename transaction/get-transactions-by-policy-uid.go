package transaction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func GetTransactionsByPolicyUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GetTransactionsByPolicyUidFx] Handler start ---------------------")

	var response GetPolicyTransactionsResp

	policyUid := r.Header.Get("policyUid")

	log.Printf("[GetTransactionsByPolicyUidFx] policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	userUid := authToken.UserID

	switch authToken.Role {
	case models.UserRoleAgent:
		if !models.IsPolicyInAgentPortfolio(userUid, policyUid) {
			log.Printf("[GetTransactionsByPolicyUidFx] policy %s is not included in agent %s", policyUid, userUid)
			return "", response, fmt.Errorf("agent %s unauthorized for policy %s", userUid, policyUid)
		}
	case models.UserRoleAgency:
		if !models.IsPolicyInAgencyPortfolio(userUid, policyUid) {
			log.Printf("[GetTransactionsByPolicyUidFx] policy %s is not included in agency %s", policyUid, userUid)
			return "", response, fmt.Errorf("agency %s unauthorized for policy %s", userUid, policyUid)
		}
	}

	transactions := GetPolicyTransactions(r.Header.Get("origin"), policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	log.Printf("[GetTransactionsByPolicyUidFx] response: %s", string(jsonOut))

	return string(jsonOut), response, err
}

func GetPolicyTransactions(origin string, policyUid string) []models.Transaction {
	var transactions Transactions

	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")

	res := lib.WhereFirestore(fireTransactions, "policyUid", "==", policyUid)

	transactions = models.TransactionToListData(res)

	sort.Sort(transactions)

	return transactions
}

func GetPolicyActiveTransactions(origin, policyUid string) []models.Transaction {
	var transactions Transactions

	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")

	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "policyUid",
				Operator:   "==",
				QueryValue: policyUid,
			},
			{
				Field:      "isDelete",
				Operator:   "==",
				QueryValue: false,
			},
		},
	}
	docsnap, err := q.FirestoreWherefields(fireTransactions)
	if err != nil {
		log.Printf("[GetPolicyActiveTransactions] query error: %s", err.Error())
		return transactions
	}

	transactions = models.TransactionToListData(docsnap)

	sort.Sort(transactions)

	return transactions
}

func (t Transactions) Len() int      { return len(t) }
func (t Transactions) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Transactions) Less(i, j int) bool {
	firstDate, _ := time.Parse(models.TimeDateOnly, t[i].ScheduleDate)
	secondDate, _ := time.Parse(models.TimeDateOnly, t[j].ScheduleDate)

	return firstDate.Before(secondDate)
}
