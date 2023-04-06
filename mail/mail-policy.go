package mail

import "github.com/wopta/goworkspace/models"

func GetMailPolicy(policy models.Policy, subject string, lines []string, isAttachment bool, at *[]Attachment) MailRequest {
	var name string
	var linkForm string
	if policy.Name == "pmi" {
		name = "Artigiani & Imprese"
		linkForm = "https://www.wopta.it/it/multi-rischio/"
	}
	if policy.Name == "persona" {
		name = "Persona"
		linkForm = "https://www.wopta.it/it/multi-rischio/"
	}
	if policy.Name == "life" {
		name = "Vita"
		linkForm = "https://www.wopta.it/it/multi-rischio/"
	}
	var obj MailRequest

	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `</br><p>Gentile ` + policy.Contractor.Name + ` ` + policy.Contractor.Surname + obj.Message
	for _, line := range lines {
		obj.Message = obj.Message + line
	}

	obj.Message = obj.Message + ` </br><p>Un saluto.</p><p>ll Team Wopta. Proteggiamo chi sei</p> </br>
	<p>Se hai bisogno di ulteriore supporto, non scrivere a questo indirizzo email, puoi compilare il <a class="button" href='` + linkForm + ` '>Form</a> oppure scrivere alla mail e verrai contattato da un nostro esperto.</p></br> `
	obj.Subject = "Wopta per te. " + name + " " + subject + " n° " + policy.CodeCompany
	obj.IsHtml = true
	obj.IsAttachment = isAttachment
	if isAttachment {
		obj.Attachments = at
	}

	return obj
}

func SendMailProposal(policy models.Policy, name string, link string) {
	var message []string
	message = append(message, `<p>richiedendo un preventivo per la soluzione assicurativa Wopta per Te `+name+` , dimostri interesse nel proteggere la tua Attività. </br> `)
	message = append(message, `Per poter valutare completamente la soluzione che sceglierai, ti alleghiamo tutti i documenti che ti consentiranno di prendere una decisione pienamente consapevole ed informata.</br>`)
	message = append(message, `Prima della sottoscrizione, leggi quanto trovi in questo <a class="button" href='`+link+`'>Link</a></p>`)

	SendMail(GetMailPolicy(policy, "", message, false, nil))

}
func SendMailPay(policy models.Policy) {
	var message []string
	message = append(message, `<p>Ti invitiamo ora ad accedere a questo link per perfezionare il pagamento.</br> `)
	message = append(message, `hai firmato correttamente la polizza. Sei più vicino a sentirti più protetto.</p> `)
	message = append(message, `<p><a class="button" href='`+policy.PayUrl+`'>Paga la tua polizza</a></p>`)
	message = append(message, `<p>Infatti senza pagamento la polizza non è attiva e, solo a pagamento avvenuto, ti invieremo una mail in cui trovi tutti i documenti contrattuali completi.</br> `)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </p> `)
	SendMail(GetMailPolicy(policy, "paga la tua polizza", message, false, nil))

}
func SendMailSign(policy models.Policy) {
	var message []string
	message = append(message, `<p>Puoi ora completare la sottoscrizione della tua polizza.</br> `)
	message = append(message, `<p>Qui trovi il link per accedere alla procedura semplice e guidata di firma elettronica avanzata tramite utilizzo di un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato.</p> 	`)
	message = append(message, `<p><a class="button" href='`+policy.SignUrl+`'>Firma Documento</a></p>`)
	message = append(message, `<p>Ultimata la procedura di firma potrai procedere al pagamento.</br> `)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </br> `)
	SendMail(GetMailPolicy(policy, "firma la tua polizza", message, false, nil))
}
func SendMailContract(policy models.Policy, at *[]Attachment) {
	var message []string
	message = append(message, `<p>La tua polizza in oggetto è ora attiva. 
	in allegato trovi i documenti da te firmati tramite l’utilizzo della Firma Elettronica Avanzata. Salva e conserva i documenti con cura, ti serviranno in caso di sinistro.
	In ogni caso ti consigliamo di scaricare l’App di Wopta per accedere alla tua area riservata nella quale troverai i tuoi documenti di polizza e altri servizi a te riservati.</p>
	<p><img  height="100" alt="instagram" src="https://storage.googleapis.com/documents-public-dev/mail/qr-app.png" style="-ms-interpolation-mode: bicubic;border: none;"></p>
	`)
	SendMail(GetMailPolicy(policy, "contratto", message, true, at))

}
