package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/reserved"
)

type RequestApprovalReq = BrokerBaseRequest

func RequestApprovalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		req    RequestApprovalReq
		policy models.Policy
	)

	log.SetPrefix("[RequestApprovalFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin = r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	log.Printf("fetching policy %s from Firestore...", req.PolicyUid)
	policy, err = plc.GetPolicy(req.PolicyUid, origin)
	if err != nil {
		log.Printf("error fetching policy %s from Firestore...", req.PolicyUid)
		return "", nil, err
	}

	if policy.ProducerUid != authToken.UserID {
		log.Printf("user %s cannot request approval for policy %s because producer not equal to request user",
			authToken.UserID, policy.Uid)
		return "", nil, errors.New("operation not allowed")
	}

	allowedStatus := []string{models.PolicyStatusInitLead, models.PolicyStatusNeedsApproval}

	if !policy.IsReserved || !lib.SliceContains(allowedStatus, policy.Status) {
		log.Printf("cannot request approval for policy with status %s and isReserved %t", policy.Status, policy.IsReserved)
		return "", nil, fmt.Errorf("cannot request approval for policy with status %s and isReserved %t", policy.Status, policy.IsReserved)
	}

	brokerUpdatePolicy(&policy, req)

	err = requestApproval(&policy)
	if err != nil {
		log.Printf("error request approval: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := policy.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(jsonOut), policy, err
}

func requestApproval(policy *models.Policy) error {
	var (
		err error
	)

	log.Println("[requestApproval] start -------------------------------------")

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("[requestApproval] starting bpmn flow...")

	state := runBrokerBpmn(policy, requestApprovalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[requestApproval] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("[requestApproval] error bpmn - state failed")
		return nil
	}

	*policy = *state.Data

	log.Printf("[requestApproval] saving policy with uid %s to bigquery...", policy.Uid)
	policy.BigquerySave(origin)

	callback_out.Execute(networkNode, *policy)

	log.Println("[requestApproval] end ---------------------------------------")

	return err
}

func setRequestApprovalData(policy *models.Policy) {
	log.Printf("[setRequestApprovalData] policy uid %s: reserved flow", policy.Uid)

	setProposalNumber(policy)

	if policy.Status == models.PolicyStatusInitLead {
		plc.AddProposalDoc(origin, policy, networkNode, mgaProduct)

		for _, reason := range policy.ReservedInfo.Reasons {
			// TODO: add key/id for reasons so we do not have to check string equallity
			if !strings.HasPrefix(reason, "Cliente già assicurato") {
				reserved.SetReservedInfo(policy, mgaProduct)
				break
			}
		}
	}

	policy.Status = models.PolicyStatusWaitForApproval
	for _, reason := range policy.ReservedInfo.Reasons {
		// TODO: add key/id for reasons so we do not have to check string equallity
		if strings.HasPrefix(reason, "Cliente già assicurato") {
			policy.Status = models.PolicyStatusWaitForApprovalMga
			break
		}
	}

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
}
