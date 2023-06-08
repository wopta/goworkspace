package payment

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type ManualPayPayload struct {
	PaymentMethod string
	Note          string
}

func ManualPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	var payload ManualPayPayload

	err := lib.CheckPayload[ManualPayPayload](body, &payload, []string{"paymentMethod"})

	origin := r.Header.Get("origin")
	transactionUid := r.Header.Get("transactionUid")
	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")
	firePolicies := lib.GetDatasetByEnv(origin, "policy")

	var transaction models.Transaction
	var policy models.Policy

	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		return "", nil, err
	}
	err = docsnap.DataTo(&transaction)
	lib.CheckError(err)

	if transaction.IsPay {
		return "", nil, fmt.Errorf("Denied: Transaction %s already paid!", transactionUid)
	}

	docsnap, err = lib.GetFirestoreErr(firePolicies, transaction.PolicyUid)
	if err != nil {
		return "", nil, err
	}
	err = docsnap.DataTo(&policy)
	lib.CheckError(err)

	if !policy.IsSign {
		return "", nil, fmt.Errorf("Denied: Policy %s not signed!", transaction.PolicyUid)
	}

	// Update transaction
	transaction.ProviderName = "manual"
	transaction.PaymentMethod = payload.PaymentMethod
	transaction.PaymentNote = payload.Note
	transaction.IsPay = true
	transaction.PayDate = time.Now().UTC()
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)

	lib.SetFirestore(fireTransactions, transactionUid, transaction)
	transaction.BigQuerySave(origin)

	// Update policy if needed
	if !policy.IsPay {
		policy.IsPay = true
		// update policy.NextPay ????

		lib.SetFirestore(firePolicies, transaction.PolicyUid, &policy)
		policy.BigquerySave(origin)
	}

	return "", nil, nil
}
