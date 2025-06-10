package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func ProposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
		req    ProposalReq
	)

	log.AddPrefix("ProposalFx")
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

	if lib.GetBoolEnv("PROPOSAL_V2") {
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.ErrorF("error proposal body: %s", err.Error())
			return "", nil, err
		}

		if req.SendEmail == nil {
			sendEmail = true
		} else {
			sendEmail = *req.SendEmail
		}

		policy, err = plc.GetPolicy(req.PolicyUid, origin)
		if err != nil {
			log.ErrorF("error fetching policy %s from Firestore...: %s", req.PolicyUid, err.Error())
			return "", nil, err
		}

		if policy.Status != models.PolicyStatusInitLead {
			log.Printf("cannot save proposal for policy with status %s", policy.Status)
			return "", nil, fmt.Errorf("cannot save proposal for policy with status %s", policy.Status)
		}

		brokerUpdatePolicy(&policy, req.BrokerBaseRequest)

		err = proposal(&policy)
		if err != nil {
			log.ErrorF("error creating proposal: %s", err.Error())
			return "", nil, err
		}
	} else {
		err = json.Unmarshal(body, &policy)
		if err != nil {
			log.ErrorF("error unmarshaling policy: %s", err.Error())
			return "", nil, err
		}

		err = leaddraft(authToken, &policy)
		if err != nil {
			log.ErrorF("error creating lead: %s", err.Error())
			return "", nil, err
		}
		utility.SetProposalNumber(&policy)
		policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(resp), &policy, err
}

func proposal(policy *models.Policy) error {
	log.AddPrefix("proposal")
	defer log.PopPrefix()
	log.Println("starting bpmn flow...")

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	state := runBrokerBpmn(policy, proposalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("error bpmn - state failed")
		return errors.New("error bpmn - state failed")
	}

	*policy = *state.Data

	log.Printf("saving proposal n. %d to bigquery...", policy.ProposalNumber)
	policy.BigquerySave(origin)

	if !policy.IsReserved {
		callback_out.Execute(networkNode, *policy, base.Proposal)
	}

	return nil
}
