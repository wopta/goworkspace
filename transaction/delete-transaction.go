package transaction

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func DeleteTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[DeleteTransactionFx] ")
	log.Println("Handler start -----------------------------------------------")

	transactionUid := r.Header.Get("transactionUid")
	origin := r.Header.Get("Origin")

	transaction := GetTransactionByUid(transactionUid, origin)
	if transaction == nil {
		errMessage := fmt.Sprintf("could not find transaction with uid '%s'", transactionUid)
		log.Println(errMessage)
		return "", "", fmt.Errorf(errMessage)
	}

	err := deleteTransaction(transaction, origin)
	if err != nil {
		log.Printf("could not delete transaction '%s'", transactionUid)
		return "", "", err
	} else {
		log.Printf("transaction '%s' successfully deleted", transactionUid)
	}

	models.CreateAuditLog(r, "")

	log.Println("Handler end -------------------------------------------------")

	return "", "", nil
}

func deleteTransaction(transaction *models.Transaction, origin string) error {
	log.Printf("deleting transaction '%s'...", transaction.Uid)

	transaction.IsDelete = true
	transaction.Status = models.TransactionStatusDeleted
	transaction.StatusHistory = append(transaction.StatusHistory, transaction.Status)
	transaction.UpdateDate = time.Now().UTC()

	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	log.Println("saving transaction to firestore...")
	err := lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
	if err != nil {
		log.Printf("error saving transaction to firestore: %s", err.Error())
		return err
	}

	log.Println("saving transaction to bigquery...")
	transaction.BigQuerySave(origin)

	return nil
}
