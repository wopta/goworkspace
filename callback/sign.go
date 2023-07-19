package callback

import (
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
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
	log.Println("[SignFx] Handler start --------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	origin := r.URL.Query().Get("origin")
	log.Printf("[SignFx] Uid: %s, Envelope: %s, Action: %s", policyUid, envelope, action)

	switch action {
	case namirialFinished:
		namirialStepFinished(origin, policyUid)
	default:
	}

	return "", nil, nil
}

func namirialStepFinished(origin, policyUid string) {
	log.Printf("[namirialStepFinished] Policy: %s", policyUid)
	var (
		policy models.Policy
		err    error
	)

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	docSnap, err := lib.GetFirestoreErr(firePolicy, policyUid)
	if err != nil {
		log.Printf("[namirialStepFinished] ERROR getting policy from firestore: %s", err.Error())
		return
	}
	err = docSnap.DataTo(&policy)
	if err != nil {
		log.Printf("[namirialStepFinished] ERROR populating policy: %s", err.Error())
		return
	}

	if !policy.IsSign && policy.Status == models.PolicyStatusToSign {
		FillAttachments(&policy)
		Sign(&policy)
		SetToPay(&policy)

		err = lib.SetFirestoreErr(firePolicy, policyUid, policy)
		if err != nil {
			log.Printf("[namirialStepFinished] ERROR saving policy to firestore: %s", err.Error())
			return
		}
		policy.BigquerySave(origin)

		mail.SendMailPay(policy)

		return
	}

	log.Printf("[namirialStepFinished] ERROR Policy %s with status %s and flag %t cannot be signed", policyUid, policy.Status, policy.IsSign)
}

func FillAttachments(policy *models.Policy) {
	if policy.Attachments == nil {
		policy.Attachments = new([]models.Attachment)
	}
	policy.Updated = time.Now().UTC()
}

func Sign(policy *models.Policy) {
	policy.IsSign = true
	policy.Status = models.PolicyStatusSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
	policy.Updated = time.Now().UTC()
}

func SetToPay(policy *models.Policy) {
	policy.Status = models.PolicyStatusToPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
	policy.Updated = time.Now().UTC()
}
