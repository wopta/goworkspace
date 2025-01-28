package broker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/callback_out"
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

const (
	approvalMga = "mga"
	approvalCompany = "company"
)

func AcceptanceFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err           error
		payload       AcceptancePayload
		policy        models.Policy
		toAddress     mail.Address
		callbackEvent string
	)

	log.SetPrefix("[AcceptanceFx]")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()

	log.Println("Handler start -----------------------------------------------")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
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

	policyUid := chi.URLParam(r, "policyUid")

	body := lib.ErrorByte(io.ReadAll(r.Body))

	if err = lib.CheckPayload(body, &payload, []string{"action"}); err != nil {
		log.Println("error checking payload")
		return "", nil, err
	}

	policy, err = plc.GetPolicy(policyUid, "")
	if err != nil {
		log.Printf("error retrieving policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}

	if !lib.SliceContains(models.GetWaitForApprovalStatusList(), policy.Status) {
		log.Printf("policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		err = fmt.Errorf("policy uid '%s': wrong status '%s'", policy.Uid, policy.Status)
		return "", nil, err
	}

	// check auth type to see what level of approval
	approvalType := ""
	switch authToken.Type {
	case lib.UserRoleAdmin:
		approvalType = approvalMga
	case lib.UserRoleCompany:
		approvalType = approvalCompany
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		callbackEvent = rejectPolicy(&policy, approvalType, lib.ToUpper(payload.Reasons))
	case models.PolicyStatusApproved:
		callbackEvent = approvePolicy(&policy, approvalType, lib.ToUpper(payload.Reasons))
	default:
		log.Printf("Unhandled action %s", payload.Action)
		return "", nil, fmt.Errorf("unhandled action %s", payload.Action)
	}

	log.Println("saving to firestore...")
	if err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, &policy); err != nil {
		log.Printf("error saving policy to firestore: %s", err.Error())
		return "", nil, err
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin)

	log.Println("sending acceptance email...")

	if callbackEvent != "" {
		// TODO: we must deceide on the process itself for comunicating changes of state
		// TODO: port acceptance into bpmn to keep code centralized and dynamic
		if networkNode = network.GetNetworkNodeByUid(policy.ProducerUid); networkNode != nil {
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
	}


	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func rejectPolicy(policy *models.Policy, approvalType, reasons string) (callbackOutEvent string) {
	if approvalType == approvalMga {
		policy.ReservedInfo.MgaApproval.Status = models.Rejected
		policy.ReservedInfo.MgaApproval.AcceptanceDate = time.Now().UTC()
		policy.ReservedInfo.MgaApproval.UpdateDate = time.Now().UTC()
		policy.ReservedInfo.MgaApproval.AcceptanceNotes = append(policy.ReservedInfo.MgaApproval.AcceptanceNotes, reasons)
		log.Printf("Policy Uid %s MGA REJECTED", policy.Uid)
	}

	if approvalType == approvalCompany {
		policy.ReservedInfo.CompanyApproval.Status = models.Rejected
		policy.ReservedInfo.CompanyApproval.AcceptanceDate = time.Now().UTC()
		policy.ReservedInfo.CompanyApproval.UpdateDate = time.Now().UTC()
		policy.ReservedInfo.CompanyApproval.AcceptanceNotes = append(policy.ReservedInfo.CompanyApproval.AcceptanceNotes, reasons)
		log.Printf("Policy Uid %s COMPANY REJECTED", policy.Uid)
	}

	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()

	log.Printf("Policy Uid %s REJECTED", policy.Uid)
	return callback_out.Rejected
}

func approvePolicy(policy *models.Policy, approvalType, reasons string) (callbackOutEvent string) {
	if approvalType == approvalMga {
		policy.ReservedInfo.MgaApproval.Status = models.Approved
		policy.ReservedInfo.MgaApproval.AcceptanceDate = time.Now().UTC()
		policy.ReservedInfo.MgaApproval.UpdateDate = time.Now().UTC()
		policy.ReservedInfo.MgaApproval.AcceptanceNotes = append(policy.ReservedInfo.MgaApproval.AcceptanceNotes, reasons)
		log.Printf("Policy Uid %s MGA APPROVED", policy.Uid)
	}

	if approvalType == approvalCompany {
		policy.ReservedInfo.CompanyApproval.Status = models.Approved
		policy.ReservedInfo.CompanyApproval.AcceptanceDate = time.Now().UTC()
		policy.ReservedInfo.CompanyApproval.UpdateDate = time.Now().UTC()
		policy.ReservedInfo.CompanyApproval.AcceptanceNotes = append(policy.ReservedInfo.CompanyApproval.AcceptanceNotes, reasons)
		log.Printf("Policy Uid %s COMPANY APPROVED", policy.Uid)
	}

	if approvalType == approvalMga && policy.ReservedInfo.CompanyApproval.Mandatory {
		policy.ReservedInfo.CompanyApproval.Status = models.WaitingApproval
		policy.ReservedInfo.CompanyApproval.UpdateDate = time.Now().UTC()
	}

	policy.Updated = time.Now().UTC()

	if approvalType == approvalCompany || approvalType == approvalMga && !policy.ReservedInfo.CompanyApproval.Mandatory {
		policy.Status = models.PolicyStatusApproved
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
		log.Printf("Policy Uid %s APPROVED", policy.Uid)
		return callback_out.Approved
	}

	return ""
}
