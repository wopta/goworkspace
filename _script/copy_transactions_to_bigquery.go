package _script

import (
	"errors"
	"log"

	"github.com/mohae/deepcopy"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/transaction"
	"google.golang.org/api/iterator"

	"time"
)

func CopyTransactionsToBigQuery() {
	//transactions := transaction.GetPolicyTransactions("", "wT6LRDMwSHViSbTCj5GM") DEV
	transactions := transaction.GetPolicyTransactions("6tztJcR7KwqBT7JHwnwI")

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
			t.BigQuerySave()
		}
		transactions[index].BigQuerySave()
	}
}

func CopyAllPoliciesTransactionToBigQuery() {
	now := time.Now().UTC()

	const batchSize = 100
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
		trans.UpdateDate = now
		trans.BigQueryParse()
		transactionsList = append(transactionsList, trans)
	}

	batches, err := divideSliceIntoBatches(transactionsList, batchSize)
	if err != nil {
		log.Printf("unable to divide slice in batches (of size %d) : %s", batchSize, err.Error())
		return
	}
	log.Printf("%d batches to save", len(batches))

	for i, batch := range batches {
		log.Printf("saving batch %d ....", i+1)
		err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, batch)
		if err != nil {
			log.Printf("unable to insert transactions into BigQuery: %s", err.Error())
			return
		}
	}

	log.Println("Finished copying all transactions to BigQuery")
}

func divideSliceIntoBatches[T any](slice []T, batchSize int) ([][]T, error) {
	if batchSize < 1 {
		return nil, errors.New("batchSize must be greater than zero")
	}

	batches := make([][]T, 0, ((len(slice)-1)/batchSize)+1)
	if len(slice) < batchSize {
		batches = append(batches, slice)
		return batches, nil
	}

	for batchSize < len(slice) {
		slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
	}
	batches = append(batches, slice)
	return batches, nil
}
