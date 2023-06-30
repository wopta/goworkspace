package test

/*

 */
import (
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var (
	signatureID = 0
)

func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("Test")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	creationDateFrom := time.Now()
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "creationDate",
				Operator:   "<",
				QueryValue: creationDateFrom,
			},
		},
	}
	fireTransactions := "transactions"
	query, _ := q.FirestoreWherefields(fireTransactions)
	transactions := models.TransactionToListData(query)
	for i, transaction := range transactions {
		transaction.BigPayDate = lib.GetBigQueryNullDateTime(transaction.PayDate)
		transaction.BigTransactionDate = lib.GetBigQueryNullDateTime(transaction.TransactionDate)
		transaction.BigCreationDate = civil.DateTimeOf(transaction.CreationDate)
		transaction.BigStatusHistory = strings.Join(transaction.StatusHistory, ",")
		log.Println(i)
		log.Println(" Transaction save BigQuery: " + transaction.Uid)
		err := lib.InsertRowsBigQuery("wopta", fireTransactions, transaction)
		if err != nil {
			log.Println("ERROR Transaction "+transaction.Uid+" save BigQuery: ", err)

		}
	}
}
