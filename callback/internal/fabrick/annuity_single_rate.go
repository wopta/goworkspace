package fabrick

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"

	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

/*
This handler is intended to handle all callbacks from fabrick that represent
a transaction that is in an already valid policy annuity. It should only pay the
transaction and have no other side effects
*/
func (FabrickCallback) AnnuitySingleRateFx(_ http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err      error
		request  = models.FabrickPaymentsRequest{}
		response = FabrickResponse{Result: true, Locale: "it"}
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

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		log.ErrorF("error decoding request body: %s", err)
		return "", nil, err
	}
	strRequest, _ := json.Marshal(request)
	response.RequestPayload = string(strRequest)

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

	paymentInfo := flow.PaymentInfoBpmn{
		Schedule:      trSchedule,
		ProviderId:    providerId,
		PaymentMethod: paymentMethod,
	}
	err = annuitySingleRate(policyUid, paymentInfo)
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

func annuitySingleRate(policyUid string, paymentInfo flow.PaymentInfoBpmn) error {
	var (
		policy models.Policy
		err    error
	)

	policy, err = plc.GetPolicy(policyUid)
	if err != nil {
		return err
	}
	if policy.Uid == "" {
		return ErrPolicyNotFound
	}

	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("paymentInfo", &paymentInfo)
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	flow, err := bpmn.GetFlow(&policy, storage)
	if err != nil {
		return err
	}
	return flow.Run("pay")
}
