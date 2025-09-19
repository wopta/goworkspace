package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/consultancy"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
	plcRenew "gitlab.dev.wopta.it/goworkspace/policy/renew"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
	trxRenew "gitlab.dev.wopta.it/goworkspace/transaction/renew"
)

func renewManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		payload     ManualPaymentPayload
		policy      models.Policy
		transaction *models.Transaction
		mgaProduct  *models.Product
		networkNode *models.NetworkNode
		flowName    string = models.ECommerceFlow
	)

	log.AddPrefix("RenewManualPaymentFx")
	defer func() {
		if err != nil {
			log.ErrorF("error: %s", err)
		}
		log.Println("Handler end -----------------------------------------------")
		log.PopPrefix()
	}()
	log.Println("Handler start ----------------------------------------------")

	err = json.NewDecoder(r.Body).Decode(&payload)
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
		return "", nil, fmt.Errorf("no methods available for manual payment")
	}

	isMethodAllowed := lib.SliceContains(methods, payload.PaymentMethod)
	if !isMethodAllowed {
		err = internal.ErrPaymentMethodNotAllowed
		return "", nil, err
	}

	transactionUid := chi.URLParam(r, "transactionUid")
	transaction = trxRenew.GetRenewTransactionByUid(transactionUid)
	if transaction == nil {
		err = errors.New("no renew transaction found")
		return "", nil, err
	}
	if transaction.IsPay {
		err = internal.ErrTransactionPaid
		return "", nil, err
	}

	policy, err = plcRenew.GetRenewPolicyByUid(transaction.PolicyUid)
	if err != nil {
		return "", nil, err
	}

	isFirstTransactionAnnuity := lib.IsEqual(policy.StartDate.AddDate(policy.Annuity, 0, 0), transaction.EffectiveDate)

	if !isFirstTransactionAnnuity {
		err = errors.New("cannot pay transaction that is not the first")
		return "", nil, err
	}

	if policy.Channel == lib.NetworkChannel {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode == nil {
			err = errors.New("networkNode not found")
			return "", nil, err
		}
		warrant := networkNode.GetWarrant()
		flowName, _ = policy.GetFlow(networkNode, warrant)
		log.Printf("flowName '%s'", flowName)
	}

	canUserAccessTransaction := authToken.Role == models.UserRoleAdmin || (authToken.IsNetworkNode &&
		authToken.UserID == policy.ProducerUid && flowName == models.RemittanceMgaFlow)
	if !canUserAccessTransaction {
		err = errors.New("user cannot access transaction")
		return "", nil, err
	}

	manualPayment(transaction, &payload)

	policy.IsPay = true
	policy.Status = models.PolicyStatusRenewed
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay, policy.Status)
	policy.Updated = time.Now().UTC()

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)

	err = tr.CreateNetworkTransactions(&policy, transaction, networkNode, mgaProduct)
	if err != nil {
		return "", nil, err
	}

	err = internal.SaveTransactionsToDB([]models.Transaction{*transaction}, lib.RenewTransactionCollection)
	if err != nil {
		return "", nil, err
	}

	if err := consultancy.GenerateInvoice(policy, *transaction); err != nil {
		log.Printf("error handling consultancy: %s", err.Error())
	}

	err = lib.SetFirestoreErr(lib.RenewPolicyCollection, policy.Uid, policy)
	if err != nil {
		return "", nil, err
	}

	policy.BigQueryParse()
	err = lib.InsertRowsBigQuery(lib.WoptaDataset, lib.RenewPolicyCollection, policy)
	if err != nil {
		return "", nil, err
	}

	policy.AddSystemNote(models.GetManualRenewNote)
	return "{}", nil, err
}
