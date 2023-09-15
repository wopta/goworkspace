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

type PutPolicyReservedPayload struct {
	Action  string `json:"action"` // models.PolicyStatusRejected (Rejected) | models.PolicyStatusApproved (Approved)
	Reasons string `json:"reasons"`
}

func PutPolicyReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PutPolicyReservedFx] Handler start ----------------------------------------")

	var (
		err     error
		payload PutPolicyReservedPayload
		policy  models.Policy
	)

	origin := r.Header.Get("origin")
	policyUid := r.Header.Get("policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Printf("[PutPolicyReservedFx] Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	log.Printf("[PutPolicyReservedFx] Payload: %s", string(body))

	err = lib.CheckPayload[PutPolicyReservedPayload](body, &payload, []string{"action"})
	if err != nil {
		log.Printf("[PutPolicyReservedFx] ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy, err = GetPolicy(policyUid, origin)
	lib.CheckError(err)

	if policy.Status != models.PolicyStatusWaitForApproval {
		log.Printf("[PutPolicyReservedFx] Policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return `{"success":false}`, `{"success":false}`, nil
	}

	switch payload.Action {
	case models.PolicyStatusRejected:
		rejectPolicy(&policy, payload.Reasons)
	case models.PolicyStatusApproved:
		approvePolicy(&policy, payload.Reasons)
	default:
		log.Printf("[PutPolicyReservedFx] Unhandled action %s", payload.Action)
		return `{"success":false}`, `{"success":false}`, nil
	}

	policy.Updated = time.Now().UTC()

	log.Println("[PutPolicyReservedFx] saving to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	lib.CheckError(err)
	log.Println("[PutPolicyReservedFx] firestore saved!")

	policy.BigquerySave(origin)

	policyJsonLog, err := policy.Marshal()
	if err != nil {
		log.Printf("[PutPolicyReservedFx] error marshaling policy: %s", err.Error())
	}
	log.Printf("[PutPolicyReservedFx] Policy: %s", string(policyJsonLog))

	log.Println("[PutPolicyReservedFx] sending acceptance email...")
	mail.SendMailReservedResult(
		policy,
		mail.AddressAssunzione,
		mail.GetAgentEmail(&policy),
		mail.Address{},
	)

	log.Println("[PutPolicyReservedFx] saving audit trail...")
	audit, err := models.ParseHttpRequest(r, string(body))
	if err != nil {
		log.Printf("[PutPolicyReservedFx] error creating audit log: %s", err.Error())
	}
	log.Printf("[PutPolicyReservedFx] audit log: %v", audit)
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
