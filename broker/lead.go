package broker

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"

	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func leadFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
	)

	log.AddPrefix("LeadFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	origin := r.Header.Get("origin")
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
	policy.ProducerUid = authToken.UserID
	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.ErrorF("error unmarshalling policy: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	err = lead(authToken, &policy, origin)
	if err != nil {
		log.ErrorF("error creating lead: %s", err.Error())
		return "", nil, err
	}

	resp, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return string(resp), &policy, err
}

func lead(authToken models.AuthToken, policy *models.Policy, origin string) error {
	var err error
	log.AddPrefix("lead")
	defer log.PopPrefix()
	log.Println("start ------------------------------------------------")

	if policy.Uid != "" {
		if err = checkIfPolicyIsLead(policy); err != nil {
			return err
		}
	} else {
		policy.Uid = lib.NewDoc(lib.PolicyCollection)
	}

	if policy.Channel == "" {
		policy.Channel = authToken.GetChannelByRoleV2()
		log.Printf("setting policy channel to '%s'", policy.Channel)
	}

	log.Println("starting bpmn flow...")
	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	flowLead, e := bpmn.GetFlow(policy, storage)
	if e != nil {
		return e
	}
	e = flowLead.Run("lead")
	if e != nil {
		return e
	}

	log.Println("saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "lead", lib.GuaranteeCollection)

	log.Println("end --------------------------------------------------")
	return err
}

func checkIfPolicyIsLead(policy *models.Policy) error {
	log.AddPrefix("checkIfPolicyIsLead")
	defer log.PopPrefix()
	var recoveredPolicy models.Policy
	policyDoc, err := lib.GetFirestoreErr(lib.PolicyCollection, policy.Uid)
	if err != nil {
		log.ErrorF("error getting policy %s from firebase: %s", policy.Uid, err.Error())
		return nil
	} else if !policyDoc.Exists() {
		log.Printf("policy %s not found on Firebase", policy.Uid)
		return nil
	}

	if err := policyDoc.DataTo(&recoveredPolicy); err != nil {
		log.ErrorF("error converting policy %s data: %s", policy.Uid, err.Error())
		return nil
	}

	allowedStatus := []string{models.PolicyStatusInit, models.PolicyStatusPartnershipLead, models.PolicyStatusInitLead}
	if !slices.Contains(allowedStatus, recoveredPolicy.Status) {
		log.ErrorF("error policy %s is not a lead", policy.Uid)
		return errors.New("policy is not a lead")
	}

	log.Printf("found lead for existing policy %s", policy.Uid)

	policy.CreationDate = recoveredPolicy.CreationDate
	policy.Status = recoveredPolicy.Status
	policy.StatusHistory = recoveredPolicy.StatusHistory
	policy.ProducerUid = recoveredPolicy.ProducerUid
	policy.ProducerCode = recoveredPolicy.ProducerCode
	policy.ProducerType = recoveredPolicy.ProducerType
	policy.NetworkUid = recoveredPolicy.NetworkUid

	return nil
}
