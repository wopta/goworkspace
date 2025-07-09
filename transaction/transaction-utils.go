package transaction

import (
	"errors"
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	transactionStatusReinitialized string = "Reinitialized"
	policyStatusReinitialized      string = "Reinitialized"
)

func ReinitializePaymentInfo(tr *models.Transaction, providerName string) error {
	if tr.IsPay && !tr.IsDelete {
		return errors.New("cannot reinitialize paid transaction")
	}

	now := time.Now().UTC()

	tr.ProviderName = providerName
	tr.IsPay = false
	tr.IsDelete = false
	tr.PaymentNote = ""
	tr.PaymentMethod = ""
	tr.PayDate = time.Time{}
	tr.PayUrl = ""
	tr.TransactionDate = time.Time{}
	tr.Status = models.TransactionStatusToPay
	tr.StatusHistory = append(tr.StatusHistory, transactionStatusReinitialized, models.TransactionStatusToPay)
	tr.UpdateDate = now

	if tr.EffectiveDate.IsZero() {
		tmp, err := time.Parse(time.DateOnly, tr.ScheduleDate)
		if err != nil {
			return err
		}
		tr.EffectiveDate = tmp
	}
	tr.ScheduleDate = tr.EffectiveDate.Format(time.DateOnly)
	tr.ExpirationDate = lib.AddMonths(now, 18).Format(time.DateOnly)

	return nil
}

func SaveTransactionsToDB(transactions []models.Transaction, collection string) error {
	batch := make(map[string]map[string]models.Transaction)
	batch[collection] = make(map[string]models.Transaction)

	for idx := range transactions {
		transactions[idx].BigQueryParse()
		batch[collection][transactions[idx].Uid] = transactions[idx]
	}

	if err := lib.SetBatchFirestoreErr(batch); err != nil {
		log.ErrorF("error saving transactions to firestore: %s", err.Error())
		return err
	}

	if err := lib.InsertRowsBigQuery(lib.WoptaDataset, collection, transactions); err != nil {
		log.ErrorF("error saving transactions to bigquery: %s", err.Error())
		return err
	}

	return nil
}
func SaveTransaction(transaction *models.Transaction, collection string) error {
	var (
		err error
	)

	transaction.BigQueryParse()
	err = lib.SetFirestoreErr(collection, transaction.Uid, transaction)
	if err != nil {
		return fmt.Errorf("error saving transaction %s in Firestore: %v", transaction.Uid, err.Error())
	}

	err = lib.InsertRowsBigQuery(lib.WoptaDataset, collection, transaction)
	if err != nil {
		log.ErrorF("error saving transaction %s in BigQuery: %v", transaction.Uid, err.Error())
		return err
	}
	return nil
}
