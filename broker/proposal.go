package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

type ProposalReq struct {
	BrokerBaseRequest
	SendEmail *bool `json:"sendEmail"`
}

func proposalFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error proposal body: %s", err.Error())
		return "", nil, err
	}

	policy, err = plc.GetPolicy(req.PolicyUid)
	if err != nil {
		log.ErrorF("error fetching policy %s from Firestore...: %s", req.PolicyUid, err.Error())
		return "", nil, err
	}

	if policy.Status != models.PolicyStatusInitLead {
		log.Printf("cannot save proposal for policy with status %s", policy.Status)
		return "", nil, fmt.Errorf("cannot save proposal for policy with status %s", policy.Status)
	}

	brokerUpdatePolicy(&policy, req.BrokerBaseRequest)

	err = proposal(&policy, *req.SendEmail)
	if err != nil {
		log.ErrorF("error creating proposal: %s", err.Error())
		return "", nil, err
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	policy.AddSystemNote(models.GetProposalNote)
	return string(resp), &policy, err
}
func proposal(policy *models.Policy, sendEmail bool) error {
	log.AddPrefix("proposal")
	defer log.PopPrefix()
	log.Println("starting bpmn flow...")
	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("is_PROPOSAL_V2", &flow.BoolBpmn{Bool: lib.GetBoolEnv("PROPOSAL_V2")})
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	storage.AddGlobal("sendEmail", &flow.BoolBpmn{
		Bool: sendEmail,
	})
	flow, err := bpmn.GetFlow(policy, storage)
	if err != nil {
		return err
	}
	err = flow.Run("proposal")
	if err != nil {
		return err
	}

	policy.BigquerySave()

	return nil
}
