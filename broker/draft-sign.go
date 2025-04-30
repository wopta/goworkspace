package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	draftbpnm "github.com/wopta/goworkspace/broker/draftBpnm"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

const (
	namirialFinished string = "workstepFinished" // when the workstep was finished
)

func DraftSignFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("SignFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	origin = r.URL.Query().Get("origin")
	sendEmailParam := r.URL.Query().Get("sendEmail")
	log.Printf("Uid: '%s', Envelope: '%s', Action: '%s', SendEmailParam: '%s'", policyUid, envelope, action, sendEmailParam)

	if v, err := strconv.ParseBool(sendEmailParam); err == nil {
		sendEmail = v
	} else {
		sendEmail = true
	}
	log.Printf("sendEmail: %t", sendEmail)

	switch action {
	case namirialFinished:
		if e := namirialStepFinished(origin, policyUid); e != nil {
			return "", nil, e
		}

	default:
	}

	log.Println("Handler end -------------------------------------------------")

	return "", nil, nil
}

func namirialStepFinished(origin, policyUid string) error {
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

	flow, err := getFlow(&policy, networkNode, draftbpnm.NewStorageBpnm())
	if err != nil {
		return err
	}
	err = flow.Run("sign")
	if err != nil {
		return err
	}

	policy.BigquerySave(origin)

	//callback_out.Execute(networkNode, policy, callback_out.Signed)
	return nil
}
