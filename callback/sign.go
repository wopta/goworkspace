package callback

import (
	"log"
	"net/http"

	broker "github.com/wopta/goworkspace/broker"
	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())
	var e error
	uid := r.URL.Query().Get("uid")
	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")

	log.Println(action)
	log.Println(envelope)
	log.Println(uid)

	if action == "workstepFinished" {
		policyF := lib.GetFirestore("policy", uid)
		var policy models.Policy
		policyF.DataTo(policy)
		company, numb := broker.GetSequenceByProduct("global")
		policy.Number = numb
		policy.NumberCompany = company
		policy.IsSign = true
		lib.SetFirestore("policy", uid, policy)
		e = lib.InsertRowsBigQuery("wopta", "policy", policy)
		mail.SendMail(getEmitMailObj(policy, policy.PayUrl))
		s := <-doc.GetFile(policy.ContractFileId, uid)
		log.Println(s)
	}

	return "", nil, e
}
func getEmitMailObj(policy models.Policy, payUrl string) mail.MailRequest {
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>Ciao ` + policy.Contractor.Name + `` + policy.Contractor.Surname + ` </p>
	<p>Polizza n° ` + policy.NumberCompany + `</p> 
	<p>Grazie per aver scelto uno dei nostri prodotti Wopta per te</p> 
	<p>Puoi ora procedere alla firma della polizza in oggetto. Qui trovi il link per
	 accedere alla procedura semplice e guidata di firma elettronica avanzata tramite utilizzo di
	  un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato. 
	Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. 
	Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.</p> 
	<p><a class="button" href='` + payUrl + `'>Firma la tua polizza:</a></p>
	<p>Ultimata la procedura di firma potrai procedere al pagamento.</p> 

	<p>Grazie per aver scelto Wopta </p> 
	<p>Proteggiamo chi sei</p> 
	`
	obj.Subject = " Wopta Contratto e pagamento"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
