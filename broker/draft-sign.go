package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	bpmn "gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	namirialFinished string = "workstepFinished" // when the workstep was finished
)

func DraftSignFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	log.AddPrefix("DraftSignFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	origin := r.URL.Query().Get("origin")
	sendEmailParam := r.URL.Query().Get("sendEmail")
	log.Printf("Uid: '%s', Envelope: '%s', Action: '%s', SendEmailParam: '%s'", policyUid, envelope, action, sendEmailParam)
	var sendEmail bool
	if v, err := strconv.ParseBool(sendEmailParam); err == nil {
		sendEmail = v
	} else {
		sendEmail = true
	}
	log.Printf("sendEmail: %t", sendEmail)

	switch action {
	case namirialFinished:
		if e := namirialStepFinished(origin, policyUid, sendEmail); e != nil {
			return "", nil, e
		}

	default:
	}

	log.Println("Handler end -------------------------------------------------")

	return "", nil, nil
}

func namirialStepFinished(origin, policyUid string, sendEmail bool) error {
	log.AddPrefix("namirialStepFinished")
	defer log.PopPrefix()
	log.Printf("Policy: %s", policyUid)
	var (
		policy models.Policy
		err    error
	)

	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)

	docSnap, err := lib.GetFirestoreErr(firePolicy, policyUid)
	if err != nil {
		return err
	}
	err = docSnap.DataTo(&policy)
	if err != nil {
		return err
	}

	if policy.IsSign || !lib.SliceContains(policy.StatusHistory, models.PolicyStatusToSign) {
		return fmt.Errorf(
			"ERROR cannot sign policy %s with isSign %t and statusHistory %s",
			policy.Uid, policy.IsSign, strings.Join(policy.StatusHistory, ","),
		)
	}

	log.Println("starting bpmn flow...")
	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("addresses", &flow.Addresses{
		FromAddress: mail.AddressAnna,
	})
	storage.AddGlobal("sendEmail", &flow.BoolBpmn{
		Bool: sendEmail,
	})
	flow, err := bpmn.GetFlow(&policy, origin, storage)
	if err != nil {
		return err
	}
	err = flow.Run("sign")
	if err != nil {
		return err
	}

	policy.BigquerySave(origin)

	return nil
}
