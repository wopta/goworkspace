package transaction

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func GetTransactionByUid(transactionUid string) *models.Transaction {
	log.AddPrefix("GetTransactionByUid")
	defer log.PopPrefix()
	log.Printf("uid %s", transactionUid)

	var (
		transaction models.Transaction
		err         error
	)

	fireTransactions := models.TransactionsCollection
	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		log.ErrorF("error getting transaction from firestore: %s", err.Error())
		return nil
	}

	err = docsnap.DataTo(&transaction)
	if err != nil {
		log.ErrorF("error converting data from document: %s", err.Error())
		return nil
	}

	return &transaction
}
