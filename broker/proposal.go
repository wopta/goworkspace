package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/question"
	"github.com/wopta/goworkspace/reserved"
)

type ProposalReq struct {
	BrokerBaseRequest
	SendEmail *bool `json:"sendEmail"`
}

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

		err = lead(authToken, &policy)
		if err != nil {
			log.ErrorF("error creating lead: %s", err.Error())
			return "", nil, err
		}
		setProposalNumber(&policy)
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
		callback_out.Execute(networkNode, *policy, callback_out.Proposal)
	}

	return nil
}

func setProposalData(policy *models.Policy) {
	log.AddPrefix("setProposalData")
	defer log.PopPrefix()
	setProposalNumber(policy)
	policy.Status = models.PolicyStatusProposal

	if policy.IsReserved {
		log.Println("setting NeedsApproval status")
		policy.Status = models.PolicyStatusNeedsApproval

		for _, reason := range policy.ReservedInfo.Reasons {
			// TODO: add key/id for reasons so we do not have to cjeck string equallity
			if !strings.HasPrefix(reason, "Cliente già assicurato") {
				reserved.SetReservedInfo(policy)
				break
			}
		}
	}

	if policy.Statements == nil || len(*policy.Statements) == 0 {
		var err error
		policy.Statements = new([]models.Statement)

		log.Println("setting policy statements")

		*policy.Statements, err = question.GetStatements(policy)
		if err != nil {
			log.ErrorF("error setting policy statements: %s", err.Error())
			return
		}

	}

	plc.AddProposalDoc(origin, policy, networkNode, mgaProduct)

	log.Printf("policy status %s", policy.Status)

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
}

func setProposalNumber(policy *models.Policy) {
	log.AddPrefix("SetProposalNumber")
	defer log.PopPrefix()
	log.Println("set proposal number start ---------------")

	if policy.ProposalNumber != 0 {
		log.Printf("proposal number already set %d", policy.ProposalNumber)
		return
	}

	log.Println("setting proposal number...")
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	policy.ProposalNumber = GetSequenceProposal("", firePolicy)
	log.Printf("proposal number %d", policy.ProposalNumber)
}
