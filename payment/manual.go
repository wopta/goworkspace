package payment

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/consultancy"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	trn "gitlab.dev.wopta.it/goworkspace/transaction"
)

type ManualPaymentPayload struct {
	PaymentMethod   string    `json:"paymentMethod"`
	PayDate         time.Time `json:"payDate"`
	TransactionDate time.Time `json:"transactionDate"`
	Note            string    `json:"note"`
}

func ManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		payload     ManualPaymentPayload
		transaction models.Transaction
		policy      models.Policy
		flowName    string = models.ECommerceFlow
		networkNode *models.NetworkNode
		warrant     *models.Warrant
		ccAddress   = mail.Address{}
		fromAddress = mail.AddressAnna
		toAddress   = mail.Address{}
	)

	log.AddPrefix("ManualPaymentFx")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err.Error())
		}
		r.Body.Close()
		log.Println("Handler end -------------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	err = lib.CheckPayload[ManualPaymentPayload](body, &payload, []string{"paymentMethod", "payDate", "transactionDate"})
	if err != nil {
		return "", nil, err
	}

	idToken := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(idToken)
	if err != nil {
		return "", nil, err
	}

	methods := models.GetAvailableMethods(authToken.Role)
	if len(methods) == 0 {
		err = fmt.Errorf("no methods available for manual payment")
		return "", nil, err
	}

	isMethodAllowed := lib.SliceContains[string](methods, payload.PaymentMethod)
	if !isMethodAllowed {
		err = fmt.Errorf("ERROR %s", internal.ErrPaymentMethodNotAllowed)
		return "", nil, err
	}

	transactionUid := chi.URLParam(r, "transactionUid")

	docsnap, err := lib.GetFirestoreErr(lib.TransactionsCollection, transactionUid)
	if err != nil {
		return "{}", nil, err
	}
	err = docsnap.DataTo(&transaction)
	if err != nil {
		return "", nil, err
	}

	docsnap, err = lib.GetFirestoreErr(lib.PolicyCollection, transaction.PolicyUid)
	if err != nil {
		return "", nil, err
	}
	err = docsnap.DataTo(&policy)
	if err != nil {
		return "", nil, err
	}

	if policy.Channel == lib.NetworkChannel {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode == nil {
			err = errors.New("networkNode not found")
			return "", nil, err
		}
		warrant = networkNode.GetWarrant()
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
		flowName, _ = policy.GetFlow(networkNode, warrant)
		log.Printf("flowName '%s'", flowName)
	}

	canUserAccessTransaction := authToken.Role == models.UserRoleAdmin || (authToken.IsNetworkNode &&
		authToken.UserID == policy.ProducerUid && flowName == models.RemittanceMgaFlow)
	if !canUserAccessTransaction {
		err = errors.New("user cannot access transaction")
		return "", nil, err
	}

	if transaction.IsPay {
		err = internal.ErrTransactionPaid
		return "", nil, err
	}

	if transaction.IsDelete {
		err = internal.ErrTransactionDeleted
		return "", nil, err
	}

	firePolicyTransactions := trn.GetPolicyValidTransactions(transaction.PolicyUid, nil)
	log.Printf("Found transactions %v", firePolicyTransactions)
	canPayTransaction := false

	for _, t := range firePolicyTransactions {
		if !t.IsPay && t.Uid != transactionUid {
			log.Printf("Next transaction to be paid should be %s", t.Uid)
			break
		}
		if t.Uid == transactionUid {
			canPayTransaction = true
			break
		}
	}

	if !canPayTransaction {
		err = internal.ErrTransactionOutOfOrder
		return "", nil, err
	}

	if !policy.IsSign {
		err = internal.ErrPolicyNotSigned
		return "", nil, err
	}

	manualPayment(&transaction, &payload)

	err = internal.SaveTransactionsToDB([]models.Transaction{transaction}, lib.TransactionsCollection)
	if err != nil {
		err = internal.ErrPaymentFailed
		return "", nil, err
	}

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	trn.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct)

	isFirstTransactionAnnuity := lib.IsEqual(policy.StartDate.AddDate(policy.Annuity, 0, 0), transaction.EffectiveDate)
	// Update policy if needed
	if !policy.IsPay && policy.Annuity == 0 {
		policy.SanitizePaymentData()

		// Add contract to policy
		err = plc.AddSignedDocumentsInPolicy(&policy, "")
		if err != nil {
			log.ErrorF("error add contract to policy: %s", err.Error())
			return "", nil, err
		}

		// Create/Update document on user collection based on contractor fiscalCode
		err = plc.SetUserIntoPolicyContractor(&policy, "")
		if err != nil {
			log.ErrorF("error set user into policy contractor: %s", err.Error())
			return "", nil, err
		}

		// Update Policy as paid
		if isFirstTransactionAnnuity {
			if err := consultancy.GenerateInvoice(policy, transaction); err != nil {
				log.Printf("error handling consultancy: %s", err.Error())
			}
		}

		err = plc.Pay(&policy, "")
		if err != nil {
			log.ErrorF("error policy pay: %s", err.Error())
			return "", nil, err
		}

		// Update NetworkNode Portfolio
		err = network.UpdateNetworkNodePortfolio("", &policy, networkNode)
		if err != nil {
			log.ErrorF("error updating %s portfolio %s", networkNode.Type, err.Error())
			return "", nil, err
		}

		policy.BigquerySave("")

		storage := bpmnEngine.NewStorageBpnm()
		flow, err := bpmn.GetFlow(&policy, r.Header.Get("Origin"), storage)
		if err != nil {
			return "", nil, err
		}
		err = flow.Run("pay")
		if err != nil {
			return "", nil, err
		}
		// Send mail with the contract to the user
		toAddress = mail.GetContractorEmail(&policy)

		// Send mail with the contract to the user
		log.Printf(
			"Sending email from '%s', to '%s', cc '%s'",
			fromAddress.String(),
			toAddress.String(),
			ccAddress.String(),
		)
		err = mail.SendMailContract(policy, policy.Attachments, fromAddress, toAddress, ccAddress, flowName)
	} else if !policy.IsPay && policy.Annuity > 0 && isFirstTransactionAnnuity {
		policy.SanitizePaymentData()
		// Update Policy as paid and renewed
		policy.IsPay = true
		policy.Status = models.PolicyStatusRenewed
		policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay, policy.Status)
		policy.Updated = time.Now().UTC()

		if isFirstTransactionAnnuity {
			if err := consultancy.GenerateInvoice(policy, transaction); err != nil {
				log.Printf("error handling consultancy: %s", err.Error())
			}
		}

		err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
		if err != nil {
			log.ErrorF("error saving policy %s to Firestore: %s", policy.Uid, err.Error())
			return "", nil, err
		}
		policy.BigquerySave("")
	}

	return "{}", nil, err
}

func manualPayment(transaction *models.Transaction, payload *ManualPaymentPayload) {
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
}
