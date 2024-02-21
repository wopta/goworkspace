package payment

import (
	"encoding/json"
	"fmt"
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
		ccAddress   = mail.Address{}
		fromAddress = mail.AddressAnna
		toAddress   = mail.Address{}
	)

	log.SetPrefix("[ManualPaymentFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := lib.CheckPayload[ManualPaymentPayload](body, &payload, []string{"paymentMethod", "payDate", "transactionDate"})
	if err != nil {
		return "", nil, err
	}

	payloadStr, _ := json.Marshal(payload)
	log.Printf("request: %s", payloadStr)

	methods := models.GetAllPaymentMethods()
	isMethodAllowed := lib.SliceContains[string](methods, payload.PaymentMethod)

	if !isMethodAllowed {
		log.Printf("ERROR %s", errPaymentMethodNotAllowed)
		return "", nil, fmt.Errorf(errPaymentMethodNotAllowed)
	}

	origin := r.Header.Get("Origin")
	transactionUid := r.Header.Get("transactionUid")
	fireTransactions := lib.GetDatasetByEnv(origin, models.TransactionsCollection)
	firePolicies := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	docsnap, err := lib.GetFirestoreErr(fireTransactions, transactionUid)
	if err != nil {
		log.Printf("ERROR get transaction from firestore: %s", err.Error())
		return "{}", nil, err
	}
	err = docsnap.DataTo(&transaction)
	lib.CheckError(err)

	if transaction.IsPay {
		log.Printf("ERROR %s", errTransactionPaid)
		return "", nil, fmt.Errorf(errTransactionPaid)
	}

	if transaction.IsDelete {
		log.Printf("ERROR %s", errTransactionDeleted)
		return "", nil, fmt.Errorf(errTransactionDeleted)
	}

	firePolicyTransactions := trn.GetPolicyActiveTransactions(origin, transaction.PolicyUid)
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
		log.Printf("ERROR %s", errTransactionOutOfOrder)
		return "", nil, fmt.Errorf(errTransactionOutOfOrder)
	}

	docsnap, err = lib.GetFirestoreErr(firePolicies, transaction.PolicyUid)
	if err != nil {
		log.Printf("ERROR get policy from firestore: %s", err.Error())
		return "", nil, err
	}
	err = docsnap.DataTo(&policy)
	lib.CheckError(err)

	if !policy.IsSign {
		log.Printf("ERROR %s", errPolicyNotSigned)
		return "", nil, fmt.Errorf(errPolicyNotSigned)
	}

	err = manualPayment(&transaction, origin, &payload)
	if err != nil {
		log.Printf("ERROR %s", errPaymentFailed)
		return "", nil, fmt.Errorf(errPaymentFailed)
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		ccAddress = mail.GetNetworkNodeEmail(networkNode)
	}
	flowName, _ = policy.GetFlow(networkNode, warrant)
	log.Printf("flowName '%s'", flowName)

	mgaProduct := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	trn.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct)

	// Update policy if needed
	if !policy.IsPay {
		// Create/Update document on user collection based on contractor fiscalCode
		err = plc.SetUserIntoPolicyContractor(&policy, origin)
		if err != nil {
			log.Printf("ERROR set user into policy contractor: %s", err.Error())
			return "", nil, err
		}

		// Add contract to policy
		err = plc.AddContract(&policy, origin)
		if err != nil {
			log.Printf("ERROR add contract to policy: %s", err.Error())
			return "", nil, err
		}

		// Update Policy as paid
		err = plc.Pay(&policy, origin)
		if err != nil {
			log.Printf("ERROR policy pay: %s", err.Error())
			return "", nil, err
		}

		// Update NetworkNode Portfolio
		err = network.UpdateNetworkNodePortfolio(origin, &policy, networkNode)
		if err != nil {
			log.Printf("[updatePolicy] error updating %s portfolio %s", networkNode.Type, err.Error())
			return "", nil, err
		}

		policy.BigquerySave(origin)

		// Send mail with the contract to the user
		switch flowName {
		case models.ProviderMgaFlow, models.MgaFlow, models.ECommerceFlow:
			toAddress = mail.GetContractorEmail(&policy)
		case models.RemittanceMgaFlow:
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		}

		// Send mail with the contract to the user
		log.Printf(
			"[updatePolicy] from '%s', to '%s', cc '%s'",
			fromAddress.String(),
			toAddress.String(),
			ccAddress.String(),
		)
		mail.SendMailContract(policy, nil, fromAddress, toAddress, ccAddress, flowName)
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
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
		log.Printf("error saving transaction to firestore: %s", err.Error())
		return err
	}

	transaction.BigQuerySave(origin)

	return nil
}
