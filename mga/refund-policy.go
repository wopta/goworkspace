package mga

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type requestRefund struct {
	Note string `json:"note"`
}

func refundPolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var req requestRefund
	policyUid := chi.URLParam(r, "policyUid")
	transactionUid := chi.URLParam(r, "transactionUid")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	//POLICY UPDATE
	policy, err := policy.GetPolicy(policyUid, "")
	if err != nil {
		return "", nil, err
	}
	if !policy.IsPay {
		return "", nil, errors.New("policy isn't paid")
	}
	if policy.Status == models.PolicyStatusRefund {
		return "", nil, errors.New("policy is already been refund")
	}
	policy.Status = models.PolicyStatusRefund
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
	policy.IsPay = false
	policy.PayUrl = ""

	//TRANSACTION UPDATE
	tr := transaction.GetTransactionByUid(transactionUid, "")
	if tr == nil {
		return "", nil, errors.New("No transaction found")
	}
	if !tr.IsPay {
		return "", nil, errors.New("transaction isn't paid")
	}
	tr.IsPay = false
	tr.PaymentNote = req.Note
	tr.UpdateDate = time.Now().UTC()
	tr.PayUrl = ""
	tr.Status = models.TransactionStatusRefund
	tr.StatusHistory = append(tr.StatusHistory, models.TransactionStatusRefund)

	err = transaction.SaveTransaction(tr, models.TransactionsCollection)
	if err != nil {
		return "", nil, err
	}

	err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
	log.Printf("policy %s saved into Firestore", policy.Uid)
	log.Printf("saving policy %s to BigQuery...", policy.Uid)
	policy.BigquerySave("")

	return "{}", nil, nil
}
