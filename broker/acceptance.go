package broker

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
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

	log.Println("[AcceptanceFx] Handler start ----------------------------------------")

	origin := r.Header.Get("origin")
	policyUid := r.Header.Get("policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Println("[AcceptanceFx] loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[AcceptanceFx] error getting authToken")
		return "", nil, err
	}

	log.Printf("[AcceptanceFx] Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[AcceptanceFx] Payload: %s", string(body))

	err = lib.CheckPayload[AcceptancePayload](body, &payload, []string{"action"})
	if err != nil {
		log.Printf("[AcceptanceFx] ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy, err = GetPolicy(policyUid, origin)
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

	switch authToken.Role {
	case models.UserRoleAdmin:
		toAddress = mail.Address{
			Address: authToken.Email,
		}
	case models.UserRoleAgent:
		toAddress = mail.GetAgentEmail(&policy)
	default:
		toAddress = mail.GetContractorEmail(&policy)
	}

	mail.SendMailReservedResult(
		policy,
		mail.AddressAssunzione,
		toAddress,
		mail.Address{},
	)

	log.Println("[AcceptanceFx] saving audit trail...")
	audit, err := models.ParseHttpRequest(r, string(body))
	if err != nil {
		log.Printf("[AcceptanceFx] error creating audit log: %s", err.Error())
	}
	log.Printf("[AcceptanceFx] audit log: %v", audit)
	audit.SaveToBigQuery()

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
