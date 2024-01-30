package broker

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

type AcceptancePayload struct {
	Action  string `json:"action"` // models.PolicyStatusRejected (Rejected) | models.PolicyStatusApproved (Approved)
	Reasons string `json:"reasons"`
}

func AcceptanceFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		payload   AcceptancePayload
		policy    models.Policy
		toAddress mail.Address
	)

	log.SetPrefix("[AcceptanceFx]")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken: %s", err.Error())
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
	policyUid := r.Header.Get("policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Printf("Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("Request Payload: %s", string(body))

	err = lib.CheckPayload[AcceptancePayload](body, &payload, []string{"action"})
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy, err = plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("error retrieving policy %s from Firestore: %s", policyUid, err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	if !lib.SliceContains(models.GetWaitForApprovalStatusList(), policy.Status) {
		log.Printf("policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return `{"success":false}`, `{"success":false}`, nil
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		rejectPolicy(&policy, lib.ToUpper(payload.Reasons))
	case models.PolicyStatusApproved:
		approvePolicy(&policy, lib.ToUpper(payload.Reasons))
	default:
		log.Printf("Unhandled action %s", payload.Action)
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy.Updated = time.Now().UTC()

	log.Println("saving to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	if err != nil {
		log.Printf("error saving policy to firestore: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin)

	policyJsonLog, err := policy.Marshal()
	if err != nil {
		log.Printf("error marshaling policy: %s", err.Error())
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

	models.CreateAuditLog(r, string(body))

	log.Println("Handler end ----------------------------------")

	return `{"success":true}`, `{"success":true}`, nil
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
