package fabrick

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/callback/internal"
	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/payment/consultancy"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
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
		request        = new(FabrickRequest)
		response       = FabrickResponse{Result: true, Locale: "it"}
	)

	log.AddPrefix("AnnuityFirstRateFx")
	defer func() {
		log.Println("Handler end ---------------------------------------------")
		log.PopPrefix()
	}()

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	trSchedule := r.URL.Query().Get("schedule")
	providerId := ""
	paymentMethod := ""

	err = json.NewDecoder(r.Body).Decode(&requestPayload)
	defer r.Body.Close()
	if err != nil {
		log.ErrorF("error decoding request body: %s", err)
		return "", nil, err
	}
	strPayload, err := request.FromPayload(requestPayload)
	if err != nil {
		log.ErrorF("error decoding request body: %s", err)
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
		log.ErrorF("error paying first annuity rate: %s", err)
		response.Result = false
	}

	stringRes, err := json.Marshal(response)
	if err != nil {
		log.ErrorF("error marshaling error response: %s", err)
	}

	return string(stringRes), response, nil
}

func annuityFirstRate(policyUid, providerId, trSchedule, paymentMethod string) error {
	var (
		policy                 models.Policy
		renewPolicy            models.Policy
		transaction            models.Transaction
		networkNode            *models.NetworkNode
		mgaProduct             *models.Product
		warrant                *models.Warrant
		err                    error
		policyCollection       string = lib.PolicyCollection
		transactionsCollection string = lib.TransactionsCollection
	)

	if policy, err = internal.GetPolicyByUidAndCollection(policyUid, lib.PolicyCollection); err != nil {
		return ErrPolicyNotFound
	}

	if renewPolicy, err = internal.GetPolicyByUidAndCollection(policyUid, lib.RenewPolicyCollection); err == nil && renewPolicy.Uid == policyUid {
		policy = renewPolicy
		policyCollection = lib.RenewPolicyCollection
		transactionsCollection = lib.RenewTransactionCollection
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if transaction, err = payTransaction(policy, providerId, trSchedule, paymentMethod, transactionsCollection, networkNode); err != nil {
		return err
	}

	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusPay)

	if policy.Annuity > 0 {
		policy.Status = models.PolicyStatusRenewed
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}

	policy.Updated = time.Now().UTC()

	policy.BigQueryParse()
	transaction.BigQueryParse()

	if policy.Annuity == 0 {
		// TODO: all methods save the data, they shouldn't to avoid data corruption
		if err = plc.SetUserIntoPolicyContractor(&policy, ""); err != nil {
			return err
		}

		if err = plc.AddNamirialDocumentsInPolicy(&policy, ""); err != nil {
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

		defer callback_out.Execute(networkNode, policy, callback_out.Paid)
	}

	if err := consultancy.GenerateInvoice(policy, transaction); err != nil {
		log.Printf("error handling consultancy: %s", err.Error())
	}

	policy.BigQueryParse()
	transaction.BigQueryParse()

	firestoreBatch := map[string]map[string]interface{}{
		policyCollection: {
			policy.Uid: policy,
		},
		transactionsCollection: {
			transaction.Uid: transaction,
		},
	}
	if err = lib.SetBatchFirestoreErr(firestoreBatch); err != nil {
		return err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, policyCollection, policy); err != nil {
		return err
	}
	if err = lib.InsertRowsBigQuery(lib.WoptaDataset, transactionsCollection, transaction); err != nil {
		return err
	}

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	return tr.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct)
}
