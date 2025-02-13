package broker

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/lib"

	// "github.com/wopta/goworkspace/mail"
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
		err     error
		payload AcceptancePayload
		policy  models.Policy
		// toAddress     mail.Address
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
	case lib.UserRoleCompany:
		approvalType = companyApprovalFlow
	case lib.UserRoleCustomer:
		approvalType = customerApprovalFlow
	case lib.UserRoleAdmin:
		approvalType = mgaApprovalFlow
		if policy.ReservedInfo.CustomerApproval.Mandatory && policy.ReservedInfo.CustomerApproval.Status != models.Approved {
			approvalType = customerApprovalFlow
		}
	case lib.UserRoleAgent, lib.UserRoleAgency, lib.UserRoleAreaManager, lib.UserRoleManager:
		if authToken.UserID == policy.ProducerUid || network.IsParentOf(authToken.UserID, policy.ProducerUid) {
			approvalType = customerApprovalFlow
		}
	}

	if approvalType == "" {
		log.Println("approval flow not set")
		err = errors.New("approval flow not set")
		return "", nil, err
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		rejectPolicy(&policy, approvalType, lib.ToUpper(payload.Reasons))
		callbackEvent = callback_out.Rejected
	case models.PolicyStatusApproved:
		approvePolicy(&policy, approvalType, lib.ToUpper(payload.Reasons))
		if policy.ReservedInfo.Approved {
			callbackEvent = callback_out.Approved
		}
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

	/*
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
	*/
	callback_out.Execute(networkNode, policy, callbackEvent)

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func rejectPolicy(policy *models.Policy, approvalFlow, reasons string) (callbackOutEvent string) {
	switch approvalFlow {
	case customerApprovalFlow:
		customerRejectPolicy(policy, reasons)
	case mgaApprovalFlow:
		mgaRejectPolicy(policy, reasons)
	case companyApprovalFlow:
		companyRejectPolicy(policy, reasons)
	}
	return callback_out.Rejected
}

func approvePolicy(policy *models.Policy, approvalFlow, reasons string) {
	switch approvalFlow {
	case customerApprovalFlow:
		customerApprovePolicy(policy, reasons)
	case mgaApprovalFlow:
		mgaApprovePolicy(policy, reasons)
	case companyApprovalFlow:
		companyApprovePolicy(policy, reasons)
	}
}

func customerApprovePolicy(policy *models.Policy, reasons string) {
	log.Println("policy customer approved")
	now := time.Now().UTC()

	policy.ReservedInfo.CustomerApproval.Status = models.Approved
	policy.ReservedInfo.CustomerApproval.AcceptanceDate = now
	policy.ReservedInfo.CustomerApproval.UpdateDate = now
	policy.ReservedInfo.CustomerApproval.AcceptanceNotes = append(policy.ReservedInfo.CustomerApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusCustomerApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = now

	if policy.ReservedInfo.MgaApproval.Mandatory {
		log.Println("policy waiting for mga approval")
		policy.ReservedInfo.MgaApproval.Status = models.WaitingApproval
		policy.ReservedInfo.MgaApproval.UpdateDate = now
		policy.Status = models.PolicyStatusWaitForApproval
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}
}

func customerRejectPolicy(policy *models.Policy, reasons string) {
	log.Println("policy customer rejected")
	now := time.Now().UTC()

	policy.ReservedInfo.CustomerApproval.Status = models.Rejected
	policy.ReservedInfo.CustomerApproval.AcceptanceDate = now
	policy.ReservedInfo.CustomerApproval.UpdateDate = now
	policy.ReservedInfo.CustomerApproval.AcceptanceNotes = append(policy.ReservedInfo.CustomerApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusCustomerRejected, policy.Status)
	policy.Updated = now

	log.Println("policy reserved rejected")
}

func mgaApprovePolicy(policy *models.Policy, reasons string) {
	log.Println("policy mga approved")
	now := time.Now().UTC()

	policy.ReservedInfo.MgaApproval.Status = models.Approved
	policy.ReservedInfo.MgaApproval.AcceptanceDate = now
	policy.ReservedInfo.MgaApproval.UpdateDate = now
	policy.ReservedInfo.MgaApproval.AcceptanceNotes = append(policy.ReservedInfo.MgaApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusMgaApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)

	if policy.ReservedInfo.CompanyApproval.Mandatory {
		log.Println("policy waiting for company approval")
		policy.ReservedInfo.CompanyApproval.Status = models.WaitingApproval
		policy.ReservedInfo.CompanyApproval.UpdateDate = now
		policy.Status = models.PolicyStatusWaitForApprovalCompany
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	} else {
		log.Println("policy reserved approved")
		policy.Status = models.PolicyStatusApproved
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	}

	policy.Updated = now
}

func mgaRejectPolicy(policy *models.Policy, reasons string) {
	log.Println("policy mga rejected")
	now := time.Now().UTC()

	policy.ReservedInfo.MgaApproval.Status = models.Rejected
	policy.ReservedInfo.MgaApproval.AcceptanceDate = now
	policy.ReservedInfo.MgaApproval.UpdateDate = now
	policy.ReservedInfo.MgaApproval.AcceptanceNotes = append(policy.ReservedInfo.MgaApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusMgaRejected, policy.Status)
	policy.Updated = now

	log.Println("policy reserved rejected")
}

func companyApprovePolicy(policy *models.Policy, reasons string) {
	log.Println("policy company approved")
	now := time.Now().UTC()

	policy.ReservedInfo.CompanyApproval.Status = models.Approved
	policy.ReservedInfo.CompanyApproval.AcceptanceDate = now
	policy.ReservedInfo.CompanyApproval.UpdateDate = now
	policy.ReservedInfo.CompanyApproval.AcceptanceNotes = append(policy.ReservedInfo.CompanyApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusApproved
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusCompanyApproved, policy.Status)
	policy.Updated = now

	log.Println("policy reserved approved")
}

func companyRejectPolicy(policy *models.Policy, reasons string) {
	log.Println("policy company rejected")
	now := time.Now().UTC()

	policy.ReservedInfo.CompanyApproval.Status = models.Rejected
	policy.ReservedInfo.CompanyApproval.AcceptanceDate = now
	policy.ReservedInfo.CompanyApproval.UpdateDate = now
	policy.ReservedInfo.CompanyApproval.AcceptanceNotes = append(policy.ReservedInfo.CompanyApproval.AcceptanceNotes, reasons)
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusCompanyRejected, policy.Status)
	policy.Updated = now

	log.Println("policy reserved rejected")
}
