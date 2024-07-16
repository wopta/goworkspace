package transaction

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

func RestoreTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		isRenew     bool
		policy      models.Policy
		transaction *models.Transaction
		collection  = lib.TransactionsCollection
	)

	log.SetPrefix("[RestoreTransactionFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	transactionUid := chi.URLParam(r, "transactionUid")
	rawIsRenew := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(rawIsRenew); err != nil && rawIsRenew != "" {
		log.Printf("error: %s", err.Error())
		return "", nil, err
	}

	if !isRenew {
		if transaction = GetTransactionByUid(transactionUid, ""); transaction == nil {
			log.Printf("no transaction found with uid: %s", transactionUid)
			return "", nil, errors.New("no transaction found")
		}
		if policy, err = plc.GetPolicy(transaction.PolicyUid, ""); err != nil {
			log.Printf("error fetching policy %s from Firestore: %s", transaction.PolicyUid, err)
			return "", nil, err
		}
	} else {
		collection = lib.RenewTransactionCollection
		if transaction = trxRenew.GetRenewTransactionByUid(transactionUid); transaction == nil {
			log.Printf("no renew transaction found with uid: %s", transactionUid)
			return "", nil, errors.New("no transaction found")
		}
		if policy, err = plcRenew.GetRenewPolicyByUid(transaction.PolicyUid); err != nil {
			log.Printf("error fetching renew policy %s from Firestore: %s", transaction.PolicyUid, err)
			return "", nil, err
		}
	}

	err = ReinitializePaymentInfo(transaction, policy.Payment)
	if err != nil {
		log.Printf("error reinitializing payment info: %s", err)
		return "", nil, err
	}

	err = saveTransactionToDb(collection, *transaction)
	if err != nil {
		log.Printf("error saving transaction: %s", err)
		return "", nil, err
	}

	rawResp, err := json.Marshal(transaction)

	return string(rawResp), transaction, err
}

func saveTransactionToDb(collection string, transaction models.Transaction) error {
	err := lib.SetFirestoreErr(collection, transaction.Uid, transaction)
	if err != nil {
		return err
	}

	transaction.BigQueryParse()
	err = lib.InsertRowsBigQuery(lib.WoptaDataset, collection, transaction)
	return err
}
