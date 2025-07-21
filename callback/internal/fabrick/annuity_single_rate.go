package fabrick

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

/*
This handler is intended to handle all callbacks from fabrick that represent
a transaction that is in an already valid policy annuity. It should only pay the
transaction and have no other side effects
*/
func (FabrickCallback) AnnuitySingleRateFx(_ http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err            error
		requestPayload FabrickRequestPayload
		request        = new(FabrickRequest)
		response       = FabrickResponse{Result: true, Locale: "it"}
	)

	log.AddPrefix("FabrickAnnuitySingleRateFx")
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

	err = annuitySingleRate(policyUid, providerId, trSchedule, paymentMethod)
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

func annuitySingleRate(policyUid, providerId, trSchedule, paymentMethod string) error {
	var (
		policy      models.Policy
		transaction models.Transaction
		networkNode *models.NetworkNode
		mgaProduct  *models.Product
		err         error
	)

	policy, err = plc.GetPolicy(policyUid)
	if err != nil {
		return err
	}
	if policy.Uid == "" {
		return ErrPolicyNotFound
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)

	if transaction, err = payTransaction(policy, providerId, trSchedule, paymentMethod, lib.TransactionsCollection, networkNode); err != nil {
		return err
	}

	firestoreBatch := map[string]map[string]interface{}{
		lib.TransactionsCollection: {
			transaction.Uid: transaction,
		},
	}
	if err = lib.SetBatchFirestoreErr(firestoreBatch); err != nil {
		return err
	}
	transaction.BigQuerySave()

	mgaProduct = prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)

	return tr.CreateNetworkTransactions(&policy, &transaction, networkNode, mgaProduct)
}
