package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func RequestApprovalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
	)

	log.Println("[RequestApprovalFx] Handler start ----------------------")

	origin = r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[RequestApprovalFx] request body: %s", string(body))
	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[RequestApprovalFx] error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	allowedStatus := []string{models.PolicyStatusInitLead, models.PolicyStatusNeedsApproval}

	if !policy.IsReserved || !lib.SliceContains(allowedStatus, policy.Status) {
		log.Printf("[ProposalFx] cannot request approval for policy with status %s and isReserved %t", policy.Status, policy.IsReserved)
		return "", nil, fmt.Errorf("cannot request approval for policy with status %s and isReserved %t", policy.Status, policy.IsReserved)
	}

	err = requestApproval(&policy)
	if err != nil {
		log.Printf("[RequestApprovalFx] error request approval: %s", err.Error())
		return "", nil, err
	}

	return "", nil, err
}

func requestApproval(policy *models.Policy) error {
	var (
		err error
	)

	log.Println("[RequestApproval] start --------------------")

	log.Println("[RequestApproval] starting bpmn flow...")
	state := runBrokerBpmn(policy, requestApprovalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[RequestApproval] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("[RequestApproval] error bpmn - state failed")
		return nil
	}

	*policy = *state.Data

	log.Printf("[RequestApproval] policy with uid %s to bigquery...", policy.Uid)
	policy.BigquerySave(origin)

	return err
}

func setRequestApprovalData(policy *models.Policy) {
	log.Printf("[setRequestApproval] policy uid %s: reserved flow", policy.Uid)

	setProposalNumber(policy)

	policy.Status = models.PolicyStatusWaitForApproval
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
}
