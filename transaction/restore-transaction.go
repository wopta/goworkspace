package transaction

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"log"
	"net/http"
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

	policy, err := plc.GetPolicy(transaction.PolicyUid, "")
	if err != nil {
		log.Printf("error fetching policy %s from Firestore: %s", transaction.PolicyUid, err)
		return "", nil, err
	}

	err = ReinitializePaymentInfo(transaction, policy.Payment)
	if err != nil {
		log.Printf("error reinitializing payment info: %s", err)
		return "", nil, err
	}

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
