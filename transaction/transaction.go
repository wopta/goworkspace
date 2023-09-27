package transaction

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Transaction")
	functions.HTTP("Transaction", Transaction)
}

func Transaction(w http.ResponseWriter, r *http.Request) {
	log.Println("Transaction")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "policy/v1/:policyUid",
				Handler: GetTransactionsByPolicyUidFx, // Broker.GetPolicyTransactions,
				Method:  http.MethodGet,
				Roles: []string{
					models.UserRoleAdmin,
					models.UserRoleManager,
					models.UserRoleAgency,
					models.UserRoleAgent,
				},
			},
		},
	}
	route.Router(w, r)
}

func SetPolicyFirstTransactionPaid(policyUid string, scheduleDate string, origin string) {
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "policyUid",
				Operator:   "==",
				QueryValue: policyUid,
			},
			{
				Field:      "scheduleDate",
				Operator:   "==",
				QueryValue: scheduleDate,
			},
		},
	}
	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")
	query, _ := q.FirestoreWherefields(fireTransactions)
	transactions := models.TransactionToListData(query)
	transaction := transactions[0]
	tr, _ := json.Marshal(transaction)
	log.Println("SetPolicyFirstTransactionPaid::payment "+policyUid+" ", string(tr))
	transaction.IsPay = true
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
	transaction.PayDate = time.Now().UTC()
	transaction.TransactionDate = time.Now().UTC()
	lib.SetFirestore(fireTransactions, transaction.Uid, transaction)
	transaction.BigQuerySave(origin)
}

func GetTransactionByPolicyUidAndScheduleDate(policyUid string, scheduleDate string, origin string) (models.Transaction, error) {
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "policyUid",
				Operator:   "==",
				QueryValue: policyUid,
			},
			{
				Field:      "scheduleDate",
				Operator:   "==",
				QueryValue: scheduleDate,
			},
		},
	}
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	query, err := q.FirestoreWherefields(fireTransactions)
	if err != nil {
		log.Printf("[GetPolicyFirstTransaction] ERROR %s", err.Error())
		return models.Transaction{}, err
	}
	transactions := models.TransactionToListData(query)
	transaction := transactions[0]

	return transaction, nil
}

func Pay(transaction *models.Transaction, origin, paymentMethod string) error {
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	transaction.IsPay = true
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
	transaction.PayDate = time.Now().UTC()
	transaction.TransactionDate = transaction.PayDate
	transaction.PaymentMethod = paymentMethod

	return lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
}
