package fabrick

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	plcRenew "gitlab.dev.wopta.it/goworkspace/policy/renew"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
	trRenew "gitlab.dev.wopta.it/goworkspace/transaction/renew"
)

type RefreshPayByLinkRequest struct {
	PolicyUid         string `json:"policyUid"`
	ScheduleFirstRate bool   `json:"scheduleFirstRate"`
}

func RefreshPayByLinkFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                   error
		request               RefreshPayByLinkRequest
		policy                models.Policy
		policyCollection      = lib.PolicyCollection
		transactions          []models.Transaction
		transactionCollection = lib.TransactionsCollection
		isRenew               bool
	)

	log.AddPrefix("RefreshPayByLinkFx")
	defer func() {
		if err != nil {
			log.Error(err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.ErrorF("error unmarshaling body")
		return "", nil, err
	}

	isRenewParam := r.URL.Query().Get("isRenew")
	if isRenew, err = strconv.ParseBool(isRenewParam); err != nil && isRenewParam != "" {
		log.ErrorF("error parsing isRenew param '%s'", isRenewParam)
		return "", nil, err
	}

	if isRenew {
		policyCollection = lib.RenewPolicyCollection
		transactionCollection = lib.RenewTransactionCollection
		if policy, err = plcRenew.GetRenewPolicyByUid(request.PolicyUid); err != nil {
			log.ErrorF("error getting renew policy")
			return "", nil, err
		}
	} else {
		policy = plc.GetPolicyByUid(request.PolicyUid, "")
	}

	if transactions, err = getTransactionsList(policy, isRenew); err != nil {
		log.ErrorF("error getting transactions")
		return "", nil, err
	}

	policy.SanitizePaymentData()

	if policy.Payment != models.FabrickPaymentProvider || policy.PaymentMode != models.PaymentModeRecurrent {
		err = fmt.Errorf("error updating payment method for policy %s with provider %s and mode %s",
			policy.Uid, policy.Payment, policy.PaymentMode)
		log.Println(err.Error())
		return "", nil, err
	}

	for index, _ := range transactions {
		transaction.ReinitializePaymentInfo(&transactions[index], policy.Payment)
		if !request.ScheduleFirstRate && index == 0 {
			transactions[index].ScheduleDate = time.Now().UTC().Format(time.DateOnly)
		}
	}

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	client := &Client{
		Policy:            policy,
		Product:           *product,
		Transactions:      transactions,
		ScheduleFirstRate: request.ScheduleFirstRate,
		CustomerId:        "",
	}
	payUrl, updatedTransactions, err := client.Update()
	if err != nil {
		log.ErrorF("error scheduling transactions on fabrick: %s", err.Error())
		return "", nil, err
	}

	err = internal.SaveTransactionsToDB(updatedTransactions, transactionCollection)
	if err != nil {
		return "", nil, err
	}

	policy.PayUrl = payUrl
	policy.BigQueryParse()

	if err = lib.SetFirestoreErr(policyCollection, policy.Uid, policy); err != nil {
		log.ErrorF("error saving policy to firestore")
		return "", nil, err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, policyCollection, policy); err != nil {
		log.ErrorF("error saving policy to bigquery")
		return "", nil, err
	}

	if err = sendPayByLinkEmail(policy); err != nil {
		log.ErrorF("error sending payment email")
		return "", nil, err
	}

	return "{}", nil, nil
}

func getTransactionsList(policy models.Policy, isRenew bool) ([]models.Transaction, error) {
	var (
		transactions []models.Transaction
		err          error
	)

	if isRenew {
		if transactions, err = trRenew.GetRenewActiveTransactionsByPolicyUid(policy.Uid, policy.Annuity); err != nil {
			log.ErrorF("error getting renew transactions")
			return nil, err
		}
	} else {
		transactions = transaction.GetPolicyTransactions("", policy.Uid)
	}

	transactions = lib.SliceFilter(transactions, func(tr models.Transaction) bool {
		return (!tr.IsPay && !tr.IsDelete) || (tr.IsPay && tr.IsDelete)
	})
	if len(transactions) == 0 {
		log.Printf("no transactions to be recreated found for policy %s", policy.Uid)
		return nil, errors.New("no transactions to be recreated found")
	}

	log.Printf("found %02d transactions for policy %s", len(transactions), policy.Uid)
	return transactions, nil
}

func sendPayByLinkEmail(policy models.Policy) error {
	var (
		flowName    string
		networkNode *models.NetworkNode
		warrant     *models.Warrant
		toAddress   mail.Address
	)

	flowName = models.ECommerceFlow
	if policy.Channel == models.NetworkChannel {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode == nil {
			return fmt.Errorf("networkNode not found")
		}
		toAddress = mail.GetNetworkNodeEmail(networkNode)
		warrant = networkNode.GetWarrant()
		if warrant == nil {
			return fmt.Errorf("warrant not found")
		}
		flowName = warrant.GetFlowName(policy.Name)
	} else {
		toAddress = mail.GetContractorEmail(&policy)
	}

	log.Printf("flowName '%s'", flowName)
	log.Printf("send pay mail to '%s'...", toAddress.String())

	mail.SendMailPay(
		policy,
		mail.AddressAnna,
		toAddress,
		mail.Address{},
		flowName,
	)

	return nil
}
