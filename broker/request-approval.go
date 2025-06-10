package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

type RequestApprovalReq = BrokerBaseRequest

func RequestApprovalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		req    RequestApprovalReq
		policy models.Policy
	)

	log.AddPrefix("RequestApprovalFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.ErrorF("error getting authToken")
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
		log.ErrorF("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	log.Printf("fetching policy %s from Firestore...", req.PolicyUid)
	policy, err = plc.GetPolicy(req.PolicyUid, origin)
	if err != nil {
		log.ErrorF("error fetching policy %s from Firestore...", req.PolicyUid)
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
		log.ErrorF("error request approval: %s", err.Error())
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
	log.AddPrefix("requestApproval")
	defer log.PopPrefix()

	log.Println("start -------------------------------------")

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("starting bpmn flow...")

	state := runBrokerBpmn(policy, requestApprovalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("error bpmn - state failed")
		return nil
	}

	*policy = *state.Data

	log.Printf("saving policy with uid %s to bigquery...", policy.Uid)
	policy.BigquerySave(origin)

	callback_out.Execute(networkNode, *policy, base.RequestApproval)

	log.Println("end ---------------------------------------")

	return err
}
