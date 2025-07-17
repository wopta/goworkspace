package broker

import (
	"encoding/json"
	"io"
	"net/http"

	bpmn "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func DraftLeadFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err    error
		policy models.Policy
	)

	log.AddPrefix("LeadFx")
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

	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.ErrorF("error unmarshalling policy: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	err = leaddraft(authToken, &policy)
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

func leaddraft(authToken models.AuthToken, policy *models.Policy) error {
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

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	log.Println("starting bpmn flow...")
	storage := bpmn.NewStorageBpnm()
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	flowLead, e := getFlow(policy, storage)
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
