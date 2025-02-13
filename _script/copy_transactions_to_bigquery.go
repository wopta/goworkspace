package _script

import (
	"errors"
	"log"

	"github.com/mohae/deepcopy"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
	"google.golang.org/api/iterator"

	"time"
)

func CopyTransactionsToBigQuery() {
	//transactions := transaction.GetPolicyTransactions("", "wT6LRDMwSHViSbTCj5GM") DEV
	transactions := transaction.GetPolicyTransactions("", "6tztJcR7KwqBT7JHwnwI")

	for index, tr := range transactions {
		if tr.IsPay {
			t := deepcopy.Copy(tr).(models.Transaction)
			t.Status = models.TransactionStatusToPay
			t.StatusHistory = []string{models.TransactionStatusToPay}
			t.IsPay = false
			t.PaymentMethod = ""
			t.UpdateDate = t.CreationDate
			t.PayDate = time.Time{}
			t.TransactionDate = time.Time{}
			t.BigQuerySave("")
		}
		transactions[index].BigQuerySave("")
	}
}

func CopyAllPoliciesTransactionToBigQuery() {
	transactionsList := make([]models.Transaction, 0)

	queries := lib.Firequeries{
		Queries: []lib.Firequery{
			{Field: "isPay", Operator: "==", QueryValue: true},
		},
	}

	iter, err := queries.FirestoreWherefields(lib.TransactionsCollection)
	if err != nil {
		log.Printf("unable to query firestore transactions: %s", err.Error())
		return
	}
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("unable to iterate over transactions: %s", err.Error())
			return
		}

		var trans models.Transaction
		err = doc.DataTo(&trans)
		if err != nil {
			log.Printf("unable to populate transaction: %s", err.Error())
			return
		}

		trans.BigQueryParse()
		transactionsList = append(transactionsList, trans)
	}

	err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, transactionsList)
	if err != nil {
		log.Printf("unable to insert transaction list into BigQuery: %s", err.Error())
		return
	}

	log.Println("Finished copying all transactions to BigQuery")
}
