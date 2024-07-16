package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/fabrick"
	tr "github.com/wopta/goworkspace/transaction"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"
)

func DeleteTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		isRenew     bool
		transaction *models.Transaction
		collection  = lib.TransactionsCollection
	)

	defer func() {
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end -------------------------------------------------")
		log.SetPrefix("")
	}()

	log.SetPrefix("[DeleteTransactionFx] ")
	log.Println("Handler start -----------------------------------------------")

	uid := chi.URLParam(r, "uid")
	rawIsRenew := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(rawIsRenew); rawIsRenew != "" && err != nil {
		log.Printf("error: %s", err.Error())
		return "", nil, err
	}

	if !isRenew {
		transaction = tr.GetTransactionByUid(uid, "")
	} else {
		collection = lib.RenewTransactionCollection
		transaction = trxRenew.GetRenewTransactionByUid(uid)
	}

	if transaction == nil {
		log.Printf("transaction '%s' not found", uid)
		return "", nil, fmt.Errorf("transaction '%s' not found", uid)
	}

	bytes, _ := json.Marshal(transaction)
	log.Printf("found transaction: %s", string(bytes))

	if transaction.ProviderName == models.FabrickPaymentProvider {
		err = fabrick.FabrickExpireBill(transaction.ProviderId)
		if err != nil {
			log.Printf("error deleting transaction on fabrick: %s", err.Error())
			return "", nil, err
		}
	}

	tr.DeleteTransaction(transaction, "Cancellata manualmente")

	err = saveTransaction(transaction, collection)
	if err != nil {
		log.Printf("%s", err.Error())
		return "", nil, err
	}

	return "{}", nil, err
}

func saveTransaction(transaction *models.Transaction, collection string) error {
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
		log.Printf("error saving transaction %s in BigQuery: %v", transaction.Uid, err.Error())
		return err
	}
	return nil
}
