package transaction

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction/renew"
)

var transactionRoutes []lib.Route = []lib.Route{
	{
		Route:  "/policy/v1/{policyUid}",
		Fn:     GetTransactionsByPolicyUidFx,
		Method: http.MethodGet,
		Roles: []string{
			models.UserRoleAdmin,
			models.UserRoleManager,
			models.UserRoleAgency,
			models.UserRoleAgent,
		},
		Entitlement: "transaction.get.transactions",
	},
	{
		Route:       "/restore/v1/{transactionUid}",
		Fn:          RestoreTransactionFx,
		Method:      http.MethodPost,
		Roles:       []string{models.UserRoleAdmin},
		Entitlement: "transaction.restore.transaction",
	},
	{
		Route:  "/policy/renew/v1/{policyUid}",
		Fn:     renew.GetRenewTransactionsByPolicyUidFx,
		Method: http.MethodGet,
		Roles: []string{
			models.UserRoleAdmin,
			models.UserRoleManager,
			models.UserRoleAgency,
			models.UserRoleAgent,
		},
		Entitlement: "transaction.get.transactions.renew",
	},
}

func init() {
	log.Println("INIT Transaction")
	functions.HTTP("Transaction", Transaction)
}

func Transaction(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("transaction", transactionRoutes)
	router.ServeHTTP(w, r)
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

func GetTransactionToBePaid(policyUid, providerId, scheduleDate, collection string) (models.Transaction, error) {
	var (
		transactions []models.Transaction
		err          error
	)

	transactions, err = getTransactionByPolicyUidAndProviderId(policyUid, providerId, collection)
	if err != nil {
		log.Printf("[GetPolicyFirstTransaction] ERROR By ProviderId %s", err.Error())
		return models.Transaction{}, err
	}

	if len(transactions) == 0 {
		transactions, err = getTransactionByPolicyUidAndScheduleDate(policyUid, scheduleDate, collection)
		if err != nil {
			log.Printf("[GetPolicyFirstTransaction] ERROR By ScheduleDate %s", err.Error())
			return models.Transaction{}, err
		}
	}

	transaction := transactions[0]

	return transaction, nil
}

func getTransactionByPolicyUidAndProviderId(policyUid, providerId, collection string) ([]models.Transaction, error) {
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "policyUid",
				Operator:   "==",
				QueryValue: policyUid,
			},
			{
				Field:      "providerId",
				Operator:   "==",
				QueryValue: providerId,
			},
			{
				Field:      "isDelete",
				Operator:   "==",
				QueryValue: false,
			},
		},
	}

	query, err := q.FirestoreWherefields(collection)
	if err != nil {
		log.Printf("[getTransactionByPolicyUidAndProviderId] ERROR %s", err.Error())
		return nil, err
	}
	return models.TransactionToListData(query), nil
}

func getTransactionByPolicyUidAndScheduleDate(policyUid, scheduleDate, collection string) ([]models.Transaction, error) {
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
			{
				Field:      "isDelete",
				Operator:   "==",
				QueryValue: false,
			},
		},
	}
	query, err := q.FirestoreWherefields(collection)
	if err != nil {
		log.Printf("[getTransactionByPolicyUidAndScheduleDate] ERROR %s", err.Error())
		return nil, err
	}
	return models.TransactionToListData(query), nil
}

func Pay(transaction *models.Transaction, origin, paymentMethod string) error {
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	transaction.IsDelete = false
	transaction.IsPay = true
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)
	transaction.PayDate = time.Now().UTC()
	transaction.TransactionDate = transaction.PayDate
	transaction.UpdateDate = transaction.PayDate
	transaction.PaymentMethod = paymentMethod

	return lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
}
