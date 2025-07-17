package callback

import (
	"net/http"
	"strconv"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	namirialFinished       string = "workstepFinished"                      // when the workstep was finished
	namirialRejected       string = "workstepRejected"                      // when the workstep was rejected
	namirialDelegated      string = "workstepDelegated"                     // whe the workstep was delegated
	namirialOpened         string = "workstepOpened"                        // when the workstep was opened
	namirialNotification   string = "sendSignNotification"                  // when the sign notification was sent
	namirialExpired        string = "envelopeExpired"                       // when the envelope was expired
	namirialActionRequired string = "workstepDelegatedSenderActionRequired" // when an action from the sender is required because of the delegation
)

func SignFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("SignFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
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
		namirialStepFinished(policyUid)
	default:
	}

	log.Println("Handler end -------------------------------------------------")

	return "", nil, nil
}

func namirialStepFinished(policyUid string) {
	log.AddPrefix("namirialStepFinished")
	defer log.PopPrefix()
	log.Printf("Policy: %s", policyUid)
	var (
		policy models.Policy
		err    error
	)

	firePolicy := lib.PolicyCollection

	docSnap, err := lib.GetFirestoreErr(firePolicy, policyUid)
	if err != nil {
		log.ErrorF("error getting policy from firestore: %s", err.Error())
		return
	}
	err = docSnap.DataTo(&policy)
	if err != nil {
		log.ErrorF("error populating policy: %s", err.Error())
		return
	}

	if policy.IsSign || !lib.SliceContains(policy.StatusHistory, models.PolicyStatusToSign) {
		log.Printf(
			"ERROR cannot sign policy %s with isSign %t and statusHistory %s",
			policy.Uid, policy.IsSign, strings.Join(policy.StatusHistory, ","),
		)
		return
	}

	log.Println("starting bpmn flow...")
	state := runCallbackBpmn(&policy, signFlowKey)
	if state == nil || state.Data == nil {
		log.ErrorF("error bpmn - state not set")
		return
	}
	if state.IsFailed {
		log.ErrorF("ERROR bpmn failed")
		return
	}
	policy = *state.Data

	policy.BigquerySave()

	callback_out.Execute(networkNode, policy, callback_out.Signed)
}
