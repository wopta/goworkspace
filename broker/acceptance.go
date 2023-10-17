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

	log.Println("[AcceptanceFx] Handler start --------------------------------")

	log.Println("[AcceptanceFx] loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[AcceptanceFx] error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"[AcceptanceFx] authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin := r.Header.Get("origin")
	policyUid := r.Header.Get("policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Printf("[AcceptanceFx] Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[AcceptanceFx] Payload: %s", string(body))

	err = lib.CheckPayload[AcceptancePayload](body, &payload, []string{"action"})
	if err != nil {
		log.Printf("[AcceptanceFx] ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy, err = plc.GetPolicy(policyUid, origin)
	lib.CheckError(err)

	if policy.Status != models.PolicyStatusWaitForApproval {
		log.Printf("[AcceptanceFx] Policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return `{"success":false}`, `{"success":false}`, nil
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		rejectPolicy(&policy, payload.Reasons)
	case models.PolicyStatusApproved:
		approvePolicy(&policy, payload.Reasons)
	default:
		log.Printf("[AcceptanceFx] Unhandled action %s", payload.Action)
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy.Updated = time.Now().UTC()

	log.Println("[AcceptanceFx] saving to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	lib.CheckError(err)
	log.Println("[AcceptanceFx] firestore saved!")

	policy.BigquerySave(origin)

	policyJsonLog, err := policy.Marshal()
	if err != nil {
		log.Printf("[AcceptanceFx] error marshaling policy: %s", err.Error())
	}
	log.Printf("[AcceptanceFx] Policy: %s", string(policyJsonLog))

	log.Println("[AcceptanceFx] sending acceptance email...")

	flowName = models.ECommerceFlow
	if policy.Channel == models.MgaChannel {
		flowName = models.MgaFlow
	} else {
		networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
		if networkNode != nil {
			warrant = networkNode.GetWarrant()
			if warrant != nil {
				flowName = warrant.GetFlowName(policy.Name)
			}
		}
	}
	log.Printf("[AcceptanceFx] flowName '%s'", flowName)

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

	log.Printf("[AcceptanceFx] toAddress '%s'", toAddress.String())

	mail.SendMailReservedResult(
		policy,
		mail.AddressAssunzione,
		toAddress,
		mail.Address{},
		flowName,
	)

	models.CreateAuditLog(r, string(body))

	log.Println("[AcceptanceFx] Handler end ----------------------------------")

	return `{"success":true}`, `{"success":true}`, nil
}

func rejectPolicy(policy *models.Policy, reasons string) {
	log.Printf("[rejectPolicy] Policy Uid %s REJECTED", policy.Uid)
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = reasons
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
}

func approvePolicy(policy *models.Policy, reasons string) {
	log.Printf("[approvePolicy] Policy Uid %s APPROVED", policy.Uid)
	policy.Status = models.PolicyStatusApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = reasons
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
}
