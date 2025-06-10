package broker

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

type AcceptancePayload struct {
	Action  string `json:"action"` // models.PolicyStatusRejected (Rejected) | models.PolicyStatusApproved (Approved)
	Reasons string `json:"reasons"`
}

func AcceptanceFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		payload       AcceptancePayload
		policy        models.Policy
		toAddress     mail.Address
		callbackEvent base.CallbackoutAction
	)

	log.AddPrefix("AcceptanceFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.ErrorF("error getting authToken: %s", err.Error())
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin := r.Header.Get("origin")
	policyUid := chi.URLParam(r, "policyUid")
	firePolicy := lib.PolicyCollection

	log.Printf("Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = lib.CheckPayload[AcceptancePayload](body, &payload, []string{"action"})
	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	policy, err = plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.ErrorF("error retrieving policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}

	if !lib.SliceContains(models.GetWaitForApprovalStatusList(), policy.Status) {
		log.Printf("policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return "", nil, fmt.Errorf("policy uid '%s': wrong status '%s'", policy.Uid, policy.Status)
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		rejectPolicy(&policy, lib.ToUpper(payload.Reasons))
		callbackEvent = base.Rejected
	case models.PolicyStatusApproved:
		approvePolicy(&policy, lib.ToUpper(payload.Reasons))
		callbackEvent = base.Approved
	default:
		log.Printf("Unhandled action %s", payload.Action)
		return "", nil, fmt.Errorf("unhandled action %s", payload.Action)
	}

	policy.Updated = time.Now().UTC()

	log.Println("saving to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	if err != nil {
		log.ErrorF("error saving policy to firestore: %s", err.Error())
		return "", nil, err
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin)

	policyJsonLog, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling policy: %s", err.Error())
	}
	log.Printf("Policy: %s", string(policyJsonLog))

	log.Println("sending acceptance email...")

	// TODO: port acceptance into bpmn to keep code centralized and dynamic
	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}
	flowName, _ = policy.GetFlow(networkNode, warrant)
	log.Printf("flowName '%s'", flowName)

	switch policy.Channel {
	case models.MgaChannel:
		toAddress = mail.Address{
			Address: authToken.Email,
		}
	case models.NetworkChannel:
		toAddress = mail.GetNetworkNodeEmail(networkNode)
	default:
		toAddress = mail.GetContractorEmail(&policy)
	}

	log.Printf("toAddress '%s'", toAddress.String())

	mail.SendMailReservedResult(
		policy,
		mail.AddressAssunzione,
		toAddress,
		mail.Address{},
		flowName,
	)

	callback_out.Execute(networkNode, policy, callbackEvent)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func rejectPolicy(policy *models.Policy, reasons string) {
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = reasons
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s REJECTED", policy.Uid)
}

func approvePolicy(policy *models.Policy, reasons string) {
	policy.Status = models.PolicyStatusApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = reasons
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s APPROVED", policy.Uid)
}
