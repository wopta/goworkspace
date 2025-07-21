package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func GetTransactionsByPolicyUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("GetTransactionsByPolicyUidFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	var response GetPolicyTransactionsResp

	policyUid := chi.URLParam(r, "policyUid")

	log.Printf("policyUid %s", policyUid)

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	lib.CheckError(err)

	userUid := authToken.UserID

	switch authToken.Role {
	case models.UserRoleAgent:
		if !models.IsPolicyInAgentPortfolio(userUid, policyUid) {
			log.Printf("policy %s is not included in agent %s", policyUid, userUid)
			return "", response, fmt.Errorf("agent %s unauthorized for policy %s", userUid, policyUid)
		}
	case models.UserRoleAgency:
		if !models.IsPolicyInAgencyPortfolio(userUid, policyUid) {
			log.Printf("policy %s is not included in agency %s", policyUid, userUid)
			return "", response, fmt.Errorf("agency %s unauthorized for policy %s", userUid, policyUid)
		}
	}

	transactions := GetPolicyTransactions(policyUid)

	response.Transactions = transactions

	jsonOut, err := json.Marshal(response)

	return string(jsonOut), response, err
}

func GetPolicyTransactions(policyUid string) []models.Transaction {
	var transactions Transactions

	fireTransactions := models.TransactionsCollection

	res := lib.WhereFirestore(fireTransactions, "policyUid", "==", policyUid)

	transactions = models.TransactionToListData(res)

	sort.Sort(transactions)

	return transactions
}

func GetPolicyValidTransactions(policyUid string, isPaid *bool) []models.Transaction {
	var transactions Transactions

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

	if isPaid != nil {
		q.Queries = append(q.Queries, lib.Firequery{
			Field:      "isPay",
			Operator:   "==",
			QueryValue: *isPaid,
		})
	}

	docsnap, err := q.FirestoreWherefields(models.TransactionsCollection)
	if err != nil {
		log.ErrorF("query error: %s", err.Error())
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
