package fabrick

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/payment/common"
	plc "github.com/wopta/goworkspace/policy"
	plcRenew "github.com/wopta/goworkspace/policy/renew"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	trRenew "github.com/wopta/goworkspace/transaction/renew"
)

type RefreshPayByLinkRequest struct {
	PolicyUid         string `json:"policyUid"`
	ScheduleFirstRate bool   `json:"scheduleFirstRate"`
	IsRenew           bool   `json:"isRenew"`
}

func RefreshPayByLinkFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                   error
		request               RefreshPayByLinkRequest
		policy                models.Policy
		policyCollection      = lib.PolicyCollection
		transactions          []models.Transaction
		transactionCollection = lib.TransactionsCollection
	)

	log.SetPrefix("[RefreshPayByLinkFx] ")
	defer func() {
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("error unmarshaling body")
		return "", nil, err
	}

	if request.IsRenew {
		policyCollection = lib.RenewPolicyCollection
		transactionCollection = lib.RenewTransactionCollection
		if policy, err = plcRenew.GetRenewPolicyByUid(request.PolicyUid); err != nil {
			log.Println("error getting renew policy")
			return "", nil, err
		}
		if transactions, err = trRenew.GetRenewTransactionsByPolicyUid(policy.Uid, policy.Annuity); err != nil {
			log.Println("error getting renew transactions")
			return "", nil, err
		}
	} else {
		policy = plc.GetPolicyByUid(request.PolicyUid, "")
		if transactions, err = getTransactionsList(policy.Uid); err != nil {
			log.Println("error getting transactions")
			return "", nil, err
		}
	}

	policy.SanitizePaymentData()

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
		log.Printf("error scheduling transactions on fabrick: %s", err.Error())
		return "", nil, err
	}

	err = common.SaveTransactionsToDB(updatedTransactions, transactionCollection)
	if err != nil {
		return "", nil, err
	}

	policy.PayUrl = payUrl
	policy.BigQueryParse()

	if err = lib.SetFirestoreErr(policyCollection, policy.Uid, policy); err != nil {
		log.Println("error saving policy to firestore")
		return "", nil, err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, policyCollection, policy); err != nil {
		log.Println("error saving policy to bigquery")
		return "", nil, err
	}

	if err = sendPayByLinkEmail(policy); err != nil {
		log.Println("error sending payment email")
		return "", nil, err
	}

	return "{}", nil, nil
}

func getTransactionsList(policyUid string) ([]models.Transaction, error) {
	transactions := transaction.GetPolicyTransactions("", policyUid)
	transactions = lib.SliceFilter(transactions, func(tr models.Transaction) bool {
		return (!tr.IsPay && !tr.IsDelete) || (tr.IsPay && tr.IsDelete)
	})
	if len(transactions) == 0 {
		log.Printf("no transactions to be recreated found for policy %s", policyUid)
		return nil, errors.New("no transactions to be recreated found")
	}
	log.Printf("found %02d transactions for policy %s", len(transactions), policyUid)
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
