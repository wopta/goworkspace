package transaction

import (
	"encoding/json"
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
	var response GetPolicyTransactionsResp

	policyUid := r.Header.Get("policyUid")

	log.Printf("GetPolicyTransactionsFx: %s", policyUid)

	res := GetPolicyTransactions(r.Header.Get("origin"), policyUid)

	response.Transactions = res

	jsonOut, err := json.Marshal(response)

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

func (t Transactions) Len() int      { return len(t) }
func (t Transactions) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Transactions) Less(i, j int) bool {
	firstDate, _ := time.Parse(models.TimeDateOnly, t[i].ScheduleDate)
	secondDate, _ := time.Parse(models.TimeDateOnly, t[j].ScheduleDate)

	return firstDate.Before(secondDate)
}
