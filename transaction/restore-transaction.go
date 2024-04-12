package transaction

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"time"
)

func RestoreTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		transaction *models.Transaction
	)

	log.SetPrefix("[RestoreTransactionFx] ")
	defer log.SetPrefix("")
	log.Println("Handler Start -----------------------------------------------")

	transactionUid := chi.URLParam(r, "transactionUid")

	transaction = GetTransactionByUid(transactionUid, "")
	if transaction == nil {
		log.Printf("no transaction found with uid: %s", transactionUid)
		return "", nil, errors.New("no transaction found")
	}

	ReinitializePaymentInfo(transaction)
	transaction.ScheduleDate = transaction.EffectiveDate.Format(time.DateOnly)
	transaction.ExpirationDate = transaction.EffectiveDate.AddDate(10, 0, 0).Format(time.DateOnly)

	err = lib.SetFirestoreErr(models.TransactionsCollection, transaction.Uid, transaction)
	if err != nil {
		log.Printf("error saving transaction %s: %s", transaction.Uid, err)
		return "", nil, err
	}

	transaction.BigQuerySave("")
	rawResp, err := json.Marshal(transaction)

	log.Println("Handler End -----------------------------------------------")

	return string(rawResp), transaction, err
}
