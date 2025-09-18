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

	"gitlab.dev.wopta.it/goworkspace/callback/internal"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

/*
This handler is intended to handle all callbacks from fabrick that represent
a transaction that is the first in it's annuity. Other then pay the related
transaction, it should also update the policy state, and promote certain data
if it is the first transaction ever
*/
func (FabrickCallback) AnnuityFirstRateFx(_ http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err      error
		request  = models.FabrickPaymentsRequest{}
		response = FabrickResponse{Result: true, Locale: "it"}
	)

	log.AddPrefix("FabrickAnnuityFirstRateFx")
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
	var policy models.Policy
	if policy, err = internal.GetPolicyByUidAndCollection(policyUid, lib.PolicyCollection); err != nil {
		return "", "", ErrPolicyNotFound
	}
	err = annuityFirstRate(&policy, paymentInfo)
	if err != nil {
		log.ErrorF("error paying first annuity rate: %s", err)
		response.Result = false
	}

	stringRes, err := json.Marshal(response)
	if err != nil {
		log.ErrorF("error marshaling error response: %s", err)
	}
	policy.AddSystemNote(models.GetPayNote)
	return string(stringRes), response, nil
}

func annuityFirstRate(policy *models.Policy, paymentInfo flow.PaymentInfoBpmn) error {
	var (
		renewPolicy models.Policy
		err         error
	)

	if renewPolicy, err = internal.GetPolicyByUidAndCollection(policy.Uid, lib.RenewPolicyCollection); err == nil && renewPolicy.Uid == policy.Uid {
		policy = &renewPolicy
	}

	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("paymentInfo", &paymentInfo)
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	storage.AddGlobal("sendEmail", &flow.BoolBpmn{Bool: false})

	flow, err := bpmn.GetFlow(policy, storage)
	if err != nil {
		return err
	}
	return flow.Run("pay")
}
