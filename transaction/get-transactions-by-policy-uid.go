package transaction

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"sort"
)

func GetTransactionsByPolicyUidFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		response GetPolicyTransactionsResp
	)

	log.Println("GetPolicyTransactions")

	fireTransactions := lib.GetDatasetByEnv(r.Header.Get("origin"), "transactions")
	policyUID := r.Header.Get("policyUid")

	res := lib.WhereFirestore(fireTransactions, "policyUid", "==", policyUID)

	response.Transactions = models.TransactionToListData(res)

	sort.Sort(response.Transactions)

	jsonOut, err := json.Marshal(response)

	return string(jsonOut), response, err
}

type GetPolicyTransactionsResp struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions []models.Transaction

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Transactions) Less(i, j int) bool { return t[i].CreationDate.Before(t[j].CreationDate) }
