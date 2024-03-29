package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func DeleteTransaction(transaction *models.Transaction, origin, note string) error {
	log.Printf("deleting transaction '%s'...", transaction.Uid)

	now := time.Now().UTC()

	transaction.IsDelete = true
	transaction.Status = models.TransactionStatusDeleted
	transaction.StatusHistory = append(transaction.StatusHistory, transaction.Status)
	transaction.UpdateDate = now
	transaction.ExpirationDate = now.AddDate(0, 0, -1).Format(models.TimeDateOnly)
	transaction.PaymentNote = note

	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	log.Println("saving transaction to firestore...")
	err := lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
	if err != nil {
		log.Printf("error saving transaction to firestore: %s", err.Error())
		return err
	}

	log.Println("saving transaction to bigquery...")
	transaction.BigQuerySave(origin)

	if transaction.IsPay {
		// TODO: decide how to handle errors on subsequent calls
		// log and ignore, or refresh data as it was before changes
		nts := GetNetworkTransactionsByTransactionUid(transaction.Uid)
		for _, nt := range nts {
			err := DeleteNetworkTransaction(&nt)
			if err != nil {
				log.Printf("error deleting network transaction '%s': %s", nt.Uid, err.Error())
				break
			}
		}
	}

	return nil
}
