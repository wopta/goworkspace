package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wopta/goworkspace/document"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
)

type ProposalReq struct {
	PolicyUid    string `json:"policyUid"`
	PaymentSplit string `json:"paymentSplit"`
	SendEmail    bool   `json:"sendEmail"`
}

func ProposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err    error
		policy models.Policy
		req    ProposalReq
	)

	log.Println("[ProposalFx] Handler start ----------------------------------")

	log.Println("[ProposalFx] loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[ProposalFx] error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"[ProposalFx] authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin = r.Header.Get("Origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[ProposalFx] Request: %s", string(body))

	if lib.GetBoolEnv("PROPOSAL_V2") {
		err = json.Unmarshal([]byte(body), &req)
		if err != nil {
			log.Printf("[ProposalFx] error proposal body: %s", err.Error())
			return "", nil, err
		}

		sendEmail = req.SendEmail

		paymentSplit = req.PaymentSplit

		policy, err = plc.GetPolicy(req.PolicyUid, origin)
		if err != nil {
			log.Printf("[ProposalFx] error fetching policy %s from Firestore...: %s", req.PolicyUid, err.Error())
			return "", nil, err
		}

		if policy.Status != models.PolicyStatusInitLead {
			log.Printf("[ProposalFx] cannot save proposal for policy with status %s", policy.Status)
			return "", nil, fmt.Errorf("cannot save proposal for policy with status %s", policy.Status)
		}

		err = proposal(&policy)
		if err != nil {
			log.Printf("[ProposalFx] error creating proposal: %s", err.Error())
			return "", nil, err
		}
	} else {
		err = json.Unmarshal([]byte(body), &policy)
		if err != nil {
			log.Printf("[ProposalFx] error unmarshaling policy: %s", err.Error())
			return "", nil, err
		}

		err = lead(authToken, &policy)
		if err != nil {
			log.Printf("[ProposalFx] error creating lead: %s", err.Error())
			return "", nil, err
		}
		setProposalNumber(&policy)
		policy.RenewDate = policy.CreationDate.AddDate(1, 0, 0)
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.Printf("[ProposalFx] error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Printf("[ProposalFx] response: %s", string(resp))
	log.Println("[ProposalFx] Handler end ------------------------------------")

	return string(resp), &policy, err
}

func proposal(policy *models.Policy) error {
	log.Println("[proposal] starting bpmn flow...")

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	state := runBrokerBpmn(policy, proposalFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[proposal] error bpmn - state not set")
		return errors.New("error on bpmn - no data present")
	}
	if state.IsFailed {
		log.Println("[proposal] error bpmn - state failed")
		return errors.New("error bpmn - state failed")
	}

	*policy = *state.Data

	log.Printf("[proposal] saving proposal n. %d to bigquery...", policy.ProposalNumber)
	policy.BigquerySave(origin)

	return nil
}

func setProposalData(policy *models.Policy) {
	log.Println("[setProposalData]")

	setProposalNumber(policy)
	policy.Status = models.PolicyStatusProposal
	policy.PaymentSplit = paymentSplit

	if policy.IsReserved {
		log.Println("[setProposalData] setting NeedsApproval status")
		policy.Status = models.PolicyStatusNeedsApproval
	}

	result := document.Proposal(origin, policy, networkNode, mgaProduct)
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}

	filename := strings.SplitN(result.LinkGcs, "/", 3)[2]
	*policy.Attachments = append(*policy.Attachments, models.Attachment{
		Name:     "Proposta",
		Link:     result.LinkGcs,
		FileName: filename,
	})

	log.Printf("[setProposalData] policy status %s", policy.Status)

	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.Updated = time.Now().UTC()
}

func setProposalNumber(policy *models.Policy) {
	log.Println("[setProposalNumber] set proposal number start ---------------")

	if policy.ProposalNumber != 0 {
		log.Printf("[setProposalNumber] proposal number already set %d", policy.ProposalNumber)
		return
	}

	log.Println("[setProposalNumber] setting proposal number...")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	policy.ProposalNumber = GetSequenceProposal("", firePolicy)
	log.Printf("[setProposalNumber] proposal number %d", policy.ProposalNumber)
}
