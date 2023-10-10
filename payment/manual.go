package payment

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

type ManualPaymentPayload struct {
	PaymentMethod   string    `json:"paymentMethod"`
	PayDate         time.Time `json:"payDate"`
	TransactionDate time.Time `json:"transactionDate"`
	Note            string    `json:"note"`
}

func ManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[ManualPaymentFx] Handler start -----------------------------------------")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	var payload ManualPaymentPayload

	err := lib.CheckPayload[ManualPaymentPayload](body, &payload, []string{"paymentMethod", "payDate", "transactionDate"})
	if err != nil {
		return "", nil, err
	}

	methods := models.GetAllPaymentMethods()
	isMethodAllowed := lib.SliceContains[string](methods, payload.PaymentMethod)

	if !isMethodAllowed {
		log.Printf("[ManualPaymentFx] ERROR %s", errPaymentMethodNotAllowed)
		errorMessage := `{"success":false, "errorMessage":"` + errPaymentMethodNotAllowed + `"}`
		return errorMessage, errorMessage, nil
	}

	origin := r.Header.Get("Origin")
	transactionUid := r.Header.Get("transactionUid")
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	firePolicies := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	var t models.Transaction
	var p models.Policy

	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		log.Printf("[ManualPaymentFx] ERROR get transaction from firestore: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&t)
	lib.CheckError(err)

	if t.IsPay {
		log.Printf("[ManualPaymentFx] ERROR %s", errTransactionPaid)
		errorMessage := `{"success":false, "errorMessage":"` + errTransactionPaid + `"}`
		return errorMessage, errorMessage, nil
	}

	firePolicyTransactions := transaction.GetPolicyTransactions(origin, t.PolicyUid)
	log.Printf("[ManualPaymentFx] Found transactions %v", firePolicyTransactions)
	canPayTransaction := false

	for _, t := range firePolicyTransactions {
		if !t.IsPay && t.Uid != transactionUid {
			log.Printf("[ManualPaymentFx] Next transaction to be paid should be %s", t.Uid)
			break
		}
		if t.Uid == transactionUid {
			canPayTransaction = true
			break
		}
	}

	if !canPayTransaction {
		log.Printf("[ManualPaymentFx] ERROR %s", errTransactionOutOfOrder)
		errorMessage := `{"success":false, "errorMessage":"` + errTransactionOutOfOrder + `"}`
		return errorMessage, errorMessage, nil
	}

	docsnap, err = lib.GetFirestoreErr(firePolicies, t.PolicyUid)
	if err != nil {
		log.Printf("[ManualPaymentFx] ERROR get policy from firestore: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&p)
	lib.CheckError(err)

	if !p.IsSign {
		log.Printf("[ManualPaymentFx] ERROR %s", errPolicyNotSigned)
		errorMessage := `{"success":false, "errorMessage":"` + errPolicyNotSigned + `"}`
		return errorMessage, errorMessage, nil
	}

	ManualPayment(&t, origin, &payload)

	producerNode := network.GetNetworkNodeByUid(p.ProducerUid)
	transaction.CreateNetworkTransactions(&p, &t, producerNode)

	// Update policy if needed
	if !p.IsPay {
		// Create/Update document on user collection based on contractor fiscalCode
		user.SetUserIntoPolicyContractor(&p, origin)

		// Get Policy contract
		gsLink := <-document.GetFileV6(p, t.PolicyUid)
		log.Println("[ManualPaymentFx] contractGsLink: ", gsLink)

		// Update Policy as paid
		policy.SetPolicyPaid(&p, gsLink, origin)

		// Send mail with the contract to the user
		mail.SendMailContract(
			p,
			nil,
			mail.AddressAnna,
			mail.GetContractorEmail(&p),
			mail.Address{},
		)
	}

	return `{"success":true}`, `{"success":true}`, nil
}

func ManualPayment(transaction *models.Transaction, origin string, payload *ManualPaymentPayload) {
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	transaction.ProviderName = models.ManualPaymentProvider
	transaction.PaymentMethod = payload.PaymentMethod
	transaction.PaymentNote = payload.Note
	transaction.IsPay = true
	transaction.PayDate = payload.PayDate
	transaction.TransactionDate = payload.TransactionDate
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)

	lib.SetFirestore(fireTransactions, transaction.Uid, transaction)
	transaction.BigQuerySave(origin)
}
