package broker

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
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
	firePolicy := lib.GetDatasetByEnv(origin, "policy")

	policy, err = GetPolicy(policyUid, origin)
	lib.CheckError(err)

	if policy.Status != models.PolicyStatusWaitForApproval {
		log.Printf("[PutPolicyReservedFx] Policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return `{"success":false}`, `{"success":false}`, nil
	}

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = lib.CheckPayload[PutPolicyReservedPayload](body, &payload, []string{"action"})
	if err != nil {
		log.Printf("[PutPolicyReservedFx] ERROR: %s", err.Error())
		return `{"success":false}`, `{"success":false}`, nil
	}

	if payload.Action == models.PolicyStatusRejected {
		log.Printf("[PutPolicyReservedFx] Policy Uid %s REJECTED", policy.Uid)
		policy.IsDeleted = true
		policy.Status = models.PolicyStatusDeleted
		policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusRejected, models.PolicyStatusDeleted)
		policy.RejectReasons = payload.Reasons
		policy.Updated = time.Now().UTC()
		err := lib.SetFirestoreErr(firePolicy, policy.Uid, policy)
		lib.CheckError(err)
		policy.BigquerySave(origin)

		policyJsonLog, _ := policy.Marshal()
		log.Printf("[PutPolicyReservedFx] Policy: %s", string(policyJsonLog))

		return `{"success":true}`, `{"success":true}`, nil
	}

	if payload.Action == models.PolicyStatusApproved {
		log.Printf("[PutPolicyReservedFx] Policy Uid %s APPROVED", policy.Uid)
		policy.Status = models.PolicyStatusApproved
		policy.StatusHistory = append(policy.StatusHistory, policy.Status)

		log.Println("[PutPolicyReservedFx] Invoking Emit")
		Emit(&policy, origin)

		return `{"success":true}`, `{"success":true}`, nil
	}

	log.Printf("[PutPolicyReservedFx] Unhandled action %s", payload.Action)
	return `{"success":false}`, `{"success":false}`, nil
}
