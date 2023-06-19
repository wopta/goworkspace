package payment

import (
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

type ManualPaymentPayload struct {
	PaymentMethod string `json:"paymentMethod"`
	Note          string `json:"note"`
}

func ManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("ManualPaymentFx")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	var payload ManualPaymentPayload

	err := lib.CheckPayload[ManualPaymentPayload](body, &payload, []string{"paymentMethod"})
	if err != nil {
		return "", nil, err
	}

	origin := r.Header.Get("origin")
	transactionUid := r.Header.Get("transactionUid")
	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")
	firePolicies := lib.GetDatasetByEnv(origin, "policy")

	var t models.Transaction
	var p models.Policy

	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		return "", nil, err
	}
	err = docsnap.DataTo(&t)
	lib.CheckError(err)

	if t.IsPay {
		log.Printf("ManualPaymentFx ERROR: %s", ERR_TR_PAID)
		return "", nil, errors.New(ERR_TR_PAID)
	}

	firePolicyTransactions := transaction.GetPolicyTransactions(origin, t.PolicyUid)
	log.Printf("ManualPaymentFx: Found transactions %v", firePolicyTransactions)
	var canPayTransaction = false

	for _, t := range firePolicyTransactions {
		if !t.IsPay && t.Uid != transactionUid {
			log.Printf("ManualPaymentFx: Next transaction to be paid should be %s", t.Uid)
			break
		}
		if t.Uid == transactionUid {
			canPayTransaction = true
			break
		}
	}

	if !canPayTransaction {
		log.Printf("ManualPaymentFx ERROR: %s", ERR_TR_OUT_OF_ORDER)
		return "", nil, errors.New(ERR_TR_OUT_OF_ORDER)
	}

	docsnap, err = lib.GetFirestoreErr(firePolicies, t.PolicyUid)
	if err != nil {
		return "", nil, err
	}
	err = docsnap.DataTo(&p)
	lib.CheckError(err)

	if !p.IsSign {
		log.Printf("ManualPaymentFx ERROR: %s", ERR_PL_NOT_SIGNED)
		return "", nil, errors.New(ERR_PL_NOT_SIGNED)
	}

	ManualPayment(&t, origin, &payload)

	// Update policy if needed
	if !p.IsPay {
		// Create/Update document on user collection based on contractor fiscalCode
		user.SetUserIntoPolicyContractor(&p, origin)

		// Get Policy contract
		gsLink := <-document.GetFileV6(p, t.PolicyUid)
		log.Println("Payment::contractGsLink: ", gsLink)

		// Update Policy as paid
		policy.SetPolicyPaid(&p, gsLink, origin)

		// Send mail with the contract to the user
		sendContractMail(&p)
	}

	return "", nil, nil
}

func ManualPayment(transaction *models.Transaction, origin string, payload *ManualPaymentPayload) {
	fireTransactions := lib.GetDatasetByEnv(origin, "transactions")

	transaction.ProviderName = "manual"
	transaction.PaymentMethod = payload.PaymentMethod
	transaction.PaymentNote = payload.Note
	transaction.IsPay = true
	transaction.PayDate = time.Now().UTC()
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)

	lib.SetFirestore(fireTransactions, transaction.Uid, transaction)
	transaction.BigQuerySave(origin)
}

func sendContractMail(policy *models.Policy) {
	log.Printf("SendContractMail: %s", policy.Uid)
	name := policy.Uid + ".pdf"
	contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "assets/users/"+
		policy.Contractor.Uid+"/contract_"+name)
	lib.CheckError(err)

	mail.SendMailContract(*policy, &[]mail.Attachment{{
		Byte:        base64.StdEncoding.EncodeToString(contractbyte),
		ContentType: "application/pdf",
		Name: policy.Contractor.Name + "_" + policy.Contractor.Surname + "_" +
			strings.ReplaceAll(policy.NameDesc, " ", "_") + "_contratto.pdf",
	}})
}
