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
			mail.SendMail(getSignMailObj(policy, policy.PayUrl))

			s := <-document.GetFileV6(policy.IdSign, uid)
			log.Println(s)
		}
	}

	return "", nil, e
}
func getSignMailObj(policy models.Policy, payUrl string) mail.MailRequest {
	var obj mail.MailRequest
	log.Println(policy.Contractor.Mail)
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>Ciao ` + policy.Contractor.Name + `` + policy.Contractor.Surname + ` </p>
	<p>Polizza n° ` + policy.CodeCompany + `</p> 
	<p>Grazie per aver scelto uno dei nostri prodotti Wopta per te</p> 
	<p>Puoi ora procedere al pagamento della polizza in oggetto. Qui trovi il link per
	 accedere alla procedura semplice e guidata 
	Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.</p> 
	<p><a class="button" href='` + payUrl + `'>Paga la tua polizza:</a></p>
	<p>A seguito.</p>
	<p>Grazie per aver scelto Wopta </p> 
	<p>Proteggiamo chi sei</p>`
	obj.Subject = " Wopta Paga la tua polizza"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
