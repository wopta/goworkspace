package callback

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

/*
*
workstepFinished : when the workstep was finished
workstepRejected : when the workstep was rejected
workstepDelegated : whe the workstep was delegated
workstepOpened : when the workstep was opened
sendSignNotification : when the sign notification was sent
envelopeExpired : when the envelope was expired
workstepDelegatedSenderActionRequired : when an action from the sender is required because of the delegation
*/
func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())
	var e error
	uid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	log.Println("Sign "+uid+" ", action)
	log.Println("Sign "+uid+" ", envelope)
	origin := r.URL.Query().Get("origin")
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	if action == "workstepFinished" {
		//policyFire=lib.GetDatasetByContractorName(policy.Contractor.Name,"policy")
		policyF := lib.GetFirestore(firePolicy, uid)
		var policy models.Policy
		policyF.DataTo(&policy)

		log.Println("Sign "+uid+" policy.Status:", policy.Status)
		if !policy.IsSign && policy.Status == models.PolicyStatusToSign {
			policy.IsSign = true
			policy.Updated = time.Now()
			policy.Status = models.PolicyStatusToPay
			policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
			policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)
			if policy.Attachments == nil {
				policy.Attachments = new([]models.Attachment)
			}
			productName := strings.ReplaceAll(policy.NameDesc, " ", "_")
			*policy.Attachments = append(*policy.Attachments, models.Attachment{Name: "Contratto",
				FileName: "Contratto_" + productName + "_" + policy.CodeCompany + ".pdf", Link: "gs://" +
					os.Getenv("GOOGLE_STORAGE_BUCKET") + "/contracts/" + policy.Uid + ".pdf"})
			lib.SetFirestore(firePolicy, uid, policy)
			policy.BigquerySave(r.Header.Get("origin"))
			mail.SendMailPay(policy)
			s := <-document.GetFileV6(policy.IdSign, uid)

			log.Println(s)
		}
	}

	return "", nil, e
}
