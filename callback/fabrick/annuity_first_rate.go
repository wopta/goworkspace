package fabrick

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

/*
This handler is intended to handle all callbacks from fabrick that represent
a transaction that is the first in it's annuity. Other then pay the related
transaction, it should also update the policy state, and promote certain data
if it is the first transaction ever
*/
func AnnuityFirstRateFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err            error
		requestPayload FabrickRequestPayload
		request        *FabrickRequest
		response       = FabrickResponse{Result: true, Locale: "it"}
	)

	log.SetPrefix("[AnnuityFirstRateFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	trSchedule := r.URL.Query().Get("schedule")
	providerId := ""
	paymentMethod := ""

	err = json.NewDecoder(r.Body).Decode(&requestPayload)
	defer r.Body.Close()
	if err != nil {
		log.Printf("error decoding request body: %s", err)
		return "", nil, err
	}
	strPayload, err := request.FromPayload(requestPayload)
	if err != nil {
		log.Printf("error decoding request body: %s", err)
		return "", nil, err
	}
	response.RequestPayload = strPayload

	if request.PaymentID == "" {
		log.Println(ErrProviderIdNotSet)
		return "", nil, ErrProviderIdNotSet
	}
	providerId = request.PaymentID

	if policyUid == "" {
		ext := strings.Split(request.ExternalID, "_")
		policyUid = ext[0]
		trSchedule = ext[1]
	}

	paymentMethod = strings.ToLower(request.Bill.Transactions[0].PaymentMethod)

	err = annuityFirstRate(policyUid, providerId, trSchedule, paymentMethod)
	if err != nil {
		log.Printf("error paying first annuity rate: %s", err)
		response.Result = false
	}

	stringRes, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling error response: %s", err)
	}

	return string(stringRes), response, nil
}

func annuityFirstRate(policyUid, providerId, trSchedule, paymentMethod string) error {
	var (
		policy      models.Policy
		transaction models.Transaction
		networkNode *models.NetworkNode
		warrant     *models.Warrant
		err         error
	)

	policy = plc.GetPolicyByUid(policyUid, "")
	if policy.Uid == "" {
		return ErrPolicyNotFound
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if transaction, err = payTransaction(policyUid, providerId, trSchedule, paymentMethod, networkNode); err != nil {
		return err
	}

	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)
	policy.Updated = time.Now().UTC()

	if policy.Annuity == 0 {
		// TODO: all methods save the data, they shouldn't to avoid data corruption
		if err = plc.SetUserIntoPolicyContractor(&policy, ""); err != nil {
			return err
		}

		if err = plc.AddContract(&policy, ""); err != nil {
			return err
		}

		if err = network.UpdateNetworkNodePortfolio("", &policy, networkNode); err != nil {
			return err
		}

		flowName, _ := policy.GetFlow(networkNode, warrant)
		toAddress := mail.Address{}
		ccAddress := mail.Address{}
		fromAddress := mail.AddressAnna

		switch flowName {
		case models.ProviderMgaFlow:
			toAddress = mail.GetContractorEmail(&policy)
			ccAddress = mail.GetNetworkNodeEmail(networkNode)
		case models.RemittanceMgaFlow:
			toAddress = mail.GetNetworkNodeEmail(networkNode)
		case models.MgaFlow, models.ECommerceFlow:
			toAddress = mail.GetContractorEmail(&policy)
		}

		mail.SendMailContract(policy, nil, fromAddress, toAddress, ccAddress, flowName)
	}

	firestoreBatch := map[string]map[string]interface{}{
		lib.PolicyCollection: {
			policy.Uid: policy,
		},
		lib.TransactionsCollection: {
			transaction.Uid: transaction,
		},
	}
	if err = lib.SetBatchFirestoreErr(firestoreBatch); err != nil {
		return err
	}
	policy.BigquerySave("")
	transaction.BigQuerySave("")

	return nil
}
