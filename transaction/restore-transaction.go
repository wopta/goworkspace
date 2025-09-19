package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	plcRenew "gitlab.dev.wopta.it/goworkspace/policy/renew"
	trxRenew "gitlab.dev.wopta.it/goworkspace/transaction/renew"
)

func restoreTransactionFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                   error
		isRenew               bool
		policy                models.Policy
		transaction           *models.Transaction
		policyCollection      = lib.PolicyCollection
		transactionCollection = lib.TransactionsCollection
	)

	log.AddPrefix("RestoreTransactionFx")
	defer func() {
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	transactionUid := chi.URLParam(r, "transactionUid")
	rawIsRenew := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(rawIsRenew); err != nil && rawIsRenew != "" {
		log.Error(err)
		return "", nil, err
	}

	if !isRenew {
		if transaction = GetTransactionByUid(transactionUid); transaction == nil {
			log.ErrorF("no transaction found with uid: %s", transactionUid)
			return "", nil, errors.New("no transaction found")
		}
		if policy, err = plc.GetPolicy(transaction.PolicyUid); err != nil {
			log.ErrorF("error fetching policy %s from Firestore: %s", transaction.PolicyUid, err)
			return "", nil, err
		}
	} else {
		policyCollection = lib.RenewPolicyCollection
		transactionCollection = lib.RenewTransactionCollection
		if transaction = trxRenew.GetRenewTransactionByUid(transactionUid); transaction == nil {
			log.ErrorF("no renew transaction found with uid: %s", transactionUid)
			return "", nil, errors.New("no transaction found")
		}
		if policy, err = plcRenew.GetRenewPolicyByUid(transaction.PolicyUid); err != nil {
			log.ErrorF("error fetching renew policy %s from Firestore: %s", transaction.PolicyUid, err)
			return "", nil, err
		}
	}

	err = ReinitializePaymentInfo(transaction, policy.Payment)
	if err != nil {
		log.ErrorF("error reinitializing payment info: %s", err)
		return "", nil, err
	}

	if lib.IsEqual(transaction.EffectiveDate, policy.StartDate.AddDate(policy.Annuity, 0, 0)) {
		policy.IsPay = false
		policy.Status = models.PolicyStatusToPay
		policy.StatusHistory = append(policy.StatusHistory, policyStatusReinitialized, policy.Status)
		policy.Updated = time.Now().UTC()
		policy.BigQueryParse()

		err = lib.SetFirestoreErr(policyCollection, policy.Uid, policy)
		if err != nil {
			return "", nil, err
		}

		err = lib.InsertRowsBigQuery(lib.WoptaDataset, policyCollection, policy)
		if err != nil {
			return "", nil, err
		}
	}

	err = SaveTransactionsToDB([]models.Transaction{*transaction}, transactionCollection)
	if err != nil {
		log.ErrorF("error saving transaction: %s", err)
		return "", nil, err
	}

	rawResp, err := json.Marshal(transaction)

	return string(rawResp), transaction, err
}
