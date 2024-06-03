package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
	"github.com/wopta/goworkspace/transaction"
	"io"
	"log"
	"net/http"
	"time"
)

type FabrickRefreshPayByLinkRequest struct {
	PolicyUid         string `json:"policyUid"`
	ScheduleFirstRate bool   `json:"scheduleFirstRate"`
}

func FabrickRefreshPayByLinkFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err     error
		request FabrickRefreshPayByLinkRequest
	)

	log.SetPrefix("[FabrickRefreshPayByLinkFx] ")
	defer log.SetPrefix("")
	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("error unmarshaling body: %s", err.Error())
		return "", nil, err
	}

	policy := plc.GetPolicyByUid(request.PolicyUid, origin)

	policy.SanitizePaymentData()

	transactions, err := getTransactionsList(policy.Uid)
	if err != nil {
		return "", nil, err
	}

	for index, _ := range transactions {
		transaction.ReinitializePaymentInfo(&transactions[index], policy.Payment)
		if !request.ScheduleFirstRate && index == 0 {
			transactions[index].ScheduleDate = time.Now().UTC().Format(time.DateOnly)
		}
	}

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	payUrl, updatedTransactions, err := Controller(policy, *product, transactions, request.ScheduleFirstRate, "")
	if err != nil {
		log.Printf("error scheduling transactions on fabrick: %s", err.Error())
		return "", nil, err
	}

	err = saveTransactionsToDB(updatedTransactions)
	if err != nil {
		return "", nil, err
	}

	policy.PayUrl = payUrl

	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		return "{}", nil, err
	}
	policy.BigquerySave("")

	err = sendPayByLinkEmail(policy)
	if err != nil {
		log.Printf("error sending payment email: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

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

func saveTransactionsToDB(transactions []models.Transaction) error {
	for _, tr := range transactions {
		err := lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			log.Printf("error saving transactions to db: %s", err.Error())
			return err
		}
		tr.BigQuerySave("")
	}
	return nil
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
