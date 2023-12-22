package payment

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	trn "github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

type ManualPaymentPayload struct {
	PaymentMethod   string    `json:"paymentMethod"`
	PayDate         time.Time `json:"payDate"`
	TransactionDate time.Time `json:"transactionDate"`
	Note            string    `json:"note"`
}

func ManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		payload     ManualPaymentPayload
		transaction models.Transaction
		policy      models.Policy
		flowName    string
		networkNode *models.NetworkNode
		warrant     *models.Warrant
		cc          = mail.Address{}
	)

	log.Println("[ManualPaymentFx] Handler start -----------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := lib.CheckPayload[ManualPaymentPayload](body, &payload, []string{"paymentMethod", "payDate", "transactionDate"})
	if err != nil {
		return "", nil, err
	}

	payloadStr, _ := json.Marshal(payload)
	log.Printf("[ManualPaymentFx] request: %s", payloadStr)

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

	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		log.Printf("[ManualPaymentFx] ERROR get transaction from firestore: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&transaction)
	lib.CheckError(err)

	if transaction.IsPay {
		log.Printf("[ManualPaymentFx] ERROR %s", errTransactionPaid)
		errorMessage := `{"success":false, "errorMessage":"` + errTransactionPaid + `"}`
		return errorMessage, errorMessage, nil
	}

	if transaction.IsDelete {
		log.Printf("[ManualPaymentFx] ERROR %s", errTransactionDeleted)
		errorMessage := `{"success":false, "errorMessage":"` + errTransactionDeleted + `"}`
		return errorMessage, errorMessage, nil
	}

	firePolicyTransactions := trn.GetPolicyActiveTransactions(origin, transaction.PolicyUid)
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

	docsnap, err = lib.GetFirestoreErr(firePolicies, transaction.PolicyUid)
	if err != nil {
		log.Printf("[ManualPaymentFx] ERROR get policy from firestore: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	err = docsnap.DataTo(&policy)
	lib.CheckError(err)

	if !policy.IsSign {
		log.Printf("[ManualPaymentFx] ERROR %s", errPolicyNotSigned)
		errorMessage := `{"success":false, "errorMessage":"` + errPolicyNotSigned + `"}`
		return errorMessage, errorMessage, nil
	}

	err = manualPayment(&transaction, origin, &payload)
	if err != nil {
		log.Printf("[ManualPaymentFx] ERROR %s", errPaymentFailed)
		errorMessage := `{"success":false, "errorMessage":"` + errPaymentFailed + `"}`
		return errorMessage, errorMessage, nil
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		cc = mail.GetNetworkNodeEmail(networkNode)
	}
	flowName, _ = policy.GetFlow(networkNode, warrant)
	log.Printf("[runBrokerBpmn] flowName '%s'", flowName)

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	trn.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct)

	// Update policy if needed
	if !policy.IsPay {
		// Create/Update document on user collection based on contractor fiscalCode
		user.SetUserIntoPolicyContractor(&policy, origin)

		// Add contract to policy
		err = plc.AddContract(&policy, origin)
		if err != nil {
			log.Printf("[ManualPaymentFx] ERROR add contract to policy: %s", err.Error())
			return `{"success":false}`, `{"success":false}`, nil
		}

		// Update Policy as paid
		plc.SetPolicyPaid(&policy, origin)

		// Send mail with the contract to the user
		mail.SendMailContract(
			policy,
			nil,
			mail.AddressAnna,
			mail.GetContractorEmail(&policy),
			cc,
			flowName,
		)
	}

	return `{"success":true}`, `{"success":true}`, nil
}

func manualPayment(transaction *models.Transaction, origin string, payload *ManualPaymentPayload) error {
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)

	transaction.ProviderName = models.ManualPaymentProvider
	transaction.PaymentMethod = payload.PaymentMethod
	transaction.PaymentNote = payload.Note
	transaction.IsPay = true
	transaction.IsDelete = false
	transaction.PayDate = payload.PayDate
	transaction.TransactionDate = payload.TransactionDate
	transaction.UpdateDate = time.Now().UTC()
	transaction.Status = models.TransactionStatusPay
	transaction.StatusHistory = append(transaction.StatusHistory, models.TransactionStatusPay)

	err := lib.SetFirestoreErr(fireTransactions, transaction.Uid, transaction)
	if err != nil {
		log.Printf("[manualPayment] error: %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	return nil
}
