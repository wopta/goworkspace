package callback

import (
	"log"
	"net/http"
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

	if action == "workstepFinished" {
		policyF := lib.GetFirestore("policy", uid)

		var policy models.Policy
		policyF.DataTo(&policy)
		log.Println("Sign "+uid+" policy.Status:", policy.Status)
		if !policy.IsSign && policy.Status == models.PolicyStatusToSign {
			policy.IsSign = true
			policy.Updated = time.Now()
			policy.Status = models.PolicyStatusToPay
			policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusSign)
			policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay)

			lib.SetFirestore("policy", uid, policy)
			e = lib.InsertRowsBigQuery("wopta", "policy", policy)

			mail.SendMailPay(policy)
			s := <-document.GetFileV6(policy.IdSign, uid)
			log.Println(s)
		}
	}

	return "", nil, e
}

/*"Gentile Nome Cognome, Spett.le ragione sociale,
hai firmato correttamente la polizza. Sei più vicino a sentirti più protetto.
Ti invitiamo ora ad accedere a questo link per perfezionare il pagamento.
Infatti senza pagamento la polizza non è attiva e, solo a pagamento avvenuto, ti invieremo una mail in cui trovi tutti i documenti contrattuali completi.
Qualora tu abbia già provveduto, ignora questa comunicazione.
Un saluto.
ll Team Wopta. Proteggiamo chi sei"
*/
