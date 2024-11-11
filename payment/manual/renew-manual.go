package manual

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/payment/common"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	prd "github.com/wopta/goworkspace/product"
	tr "github.com/wopta/goworkspace/transaction"
	trxRenew "github.com/wopta/goworkspace/transaction/renew"
)

func RenewManualPaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		payload     ManualPaymentPayload
		policy      models.Policy
		transaction *models.Transaction
		mgaProduct  *models.Product
		networkNode *models.NetworkNode
	)

	log.SetPrefix("[RenewManualPaymentFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end -----------------------------------------------")
		log.SetPrefix("")
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

	methods := models.GetAllPaymentMethods(authToken.Role)
	isMethodAllowed := lib.SliceContains(methods, payload.PaymentMethod)

	if !isMethodAllowed {
		err = errors.New(errPaymentMethodNotAllowed)
		return "", nil, err
	}

	transactionUid := chi.URLParam(r, "transactionUid")
	transaction = trxRenew.GetRenewTransactionByUid(transactionUid)
	if transaction == nil {
		err = errors.New("no renew transaction found")
		return "", nil, err
	}
	if transaction.IsPay {
		err = errors.New(errTransactionPaid)
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

	err = common.SaveTransactionsToDB([]models.Transaction{*transaction}, lib.RenewTransactionCollection)
	if err != nil {
		return "", nil, err
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

	return "{}", nil, err
}
