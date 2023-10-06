package transaction

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetTransactionByUid(transactionUid, origin string) *models.Transaction {
	log.Printf("[GetTransactionByUid] uid %s", transactionUid)

	var (
		transaction models.Transaction
		err         error
	)

	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		log.Printf("[GetTransactionByUid] error getting transaction from firestore: %s", err.Error())
		return nil
	}

	err = docsnap.DataTo(&transaction)
	if err != nil {
		log.Printf("[GetTransactionByUid] error converting data from document: %s", err.Error())
		return nil
	}

	return &transaction
}
