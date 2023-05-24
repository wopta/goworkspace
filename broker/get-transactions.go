package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"sort"
)

func GetPolicyTransactions(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		transactions Transactions
	)

	log.Println("GetPolicyTransactions")

	fireTransactions := lib.GetDatasetByEnv(r.Header.Get("origin"), "transactions")
	policyUID := r.Header.Get("policyUid")

	res := lib.WhereFirestore(fireTransactions, "policyUid", "==", policyUID)

	transactions = models.TransactionToListData(res)

	sort.Sort(transactions)

	jsonOut, err := json.Marshal(transactions)

	return string(jsonOut), transactions, err
}

type Transactions []models.Transaction

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Transactions) Less(i, j int) bool { return t[i].CreationDate.Before(t[j].CreationDate) }
