package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

func GetPolicyTransactions(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		transactions []models.Transaction
	)

	log.Println("GetPolicyTransactions")

	fireTransactions := lib.GetDatasetByEnv(r.Header.Get("origin"), "transactions")
	policyUID := r.Header.Get("policyUid")

	res := lib.WhereFirestore(fireTransactions, "policyUid", "==", policyUID)

	transactions = models.TransactionToListData(res)

	jsonOut, err := json.Marshal(transactions)

	return string(jsonOut), transactions, err
}
