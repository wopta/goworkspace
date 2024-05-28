package fabrick

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

/*
This handler is intended to handle all callbacks from fabrick that represent
a transaction that is in an already valid policy annuity. It should only pay the
transaction and have no other side effects
*/
func AnnuitySingleRateFx(_ http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err            error
		requestPayload FabrickRequestPayload
		request        = new(FabrickRequest)
		response       = FabrickResponse{Result: true, Locale: "it"}
	)

	log.SetPrefix("[AnnuitySingleRateFx] ")
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

	err = annuitySingleRate(policyUid, providerId, trSchedule, paymentMethod)
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

func annuitySingleRate(policyUid, providerId, trSchedule, paymentMethod string) error {
	var (
		policy      models.Policy
		transaction models.Transaction
		networkNode *models.NetworkNode
		err         error
	)

	policy = plc.GetPolicyByUid(policyUid, "")
	if policy.Uid == "" {
		return ErrPolicyNotFound
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)

	if transaction, err = payTransaction(policyUid, providerId, trSchedule, paymentMethod, networkNode); err != nil {
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
	transaction.BigQuerySave("")

	return nil
}
