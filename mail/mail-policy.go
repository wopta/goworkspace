package mail

import (
	"github.com/wopta/goworkspace/models"
)

func GetMailPolicy(policy models.Policy, subject string, islink bool, link string, linkLabel string, lines []string, isAttachment bool, at *[]Attachment) MailRequest {
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

	obj.From = "anna@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">Ciao <strong>` + policy.Contractor.Name + " " + policy.Contractor.Surname + `</strong>,</p> `
	for _, line := range lines {
		obj.Message = obj.Message + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}

	obj.Message = obj.Message + ` 
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">Se hai bisogno di ulteriore supporto, non scrivere a questo indirizzo email, puoi compilare il <a class="button" href='` + linkForm + ` '>Form </a> oppure scrivere alla mail e verrai contattato da un nostro esperto.</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
	obj.Subject = "Wopta per te. " + name + " "
	obj.SubTitle = subject
	obj.IsHtml = true
	obj.IsAttachment = isAttachment
	obj.IsLink = islink
	if islink {
		obj.Link = link
		obj.LinkLabel = linkLabel
	} else {
		obj.IsApp = true
	}
	if isAttachment {
		obj.Attachments = at
	}

	return obj
}

func SendMailProposal(policy models.Policy) {
	var link, name string
	if policy.Name == "pmi" {

		link = "https://storage.googleapis.com/documents-public-dev/information-sets/pmi/v1/Precontrattuale.pdf"

	}
	if policy.Name == "persona" {
		link = "https://storage.googleapis.com/documents-public-dev/information-sets/persona/v1/Precontrattuale.pdf"
	}
	if policy.Name == "life" {
		link = "https://storage.googleapis.com/documents-public-dev/information-sets/life/v1/Precontrattuale.pdf"

	}
	var message []string
	message = append(message, `	<br>Grazie per aver compilato, la richiesta di preventivo per una polizza <b>`+name+`</b>.
	<br><br>Per poter valutare completamente la soluzione che sceglierai, ti alleghiamo tutti i documenti che ti consentiranno di prendere una decisione pienamente consapevole ed informata.<br> `)
	message = append(message, `<br>Se ci fosse la necessità di richiederti alcune informazioni aggiuntive, ti contatteremo ai recapiti da te forniti.`)
	message = append(message, `<br>Prima della sottoscrizione, leggi quanto trovi in questa mail, la <b>documentazione precontrattuale</b> che, per trasparenza e tua adeguata informazione, trovi sempre disponibile al link sotto.`)

	SendMail(GetMailPolicy(policy, "", true, link, "Leggi documentazione", message, false, nil))

}

func SendMailPay(policy models.Policy) {
	var message []string
	message = append(message, `<p>hai firmato correttamente la polizza. Sei più vicino a sentirti più protetto.</br> `)
	message = append(message, `Ti invitiamo ora ad accedere a questo link per perfezionare il pagamento.</p> `)
	message = append(message, `<p>Infatti senza pagamento la polizza non è attiva e, solo a pagamento avvenuto, ti invieremo una mail in cui trovi tutti i documenti contrattuali completi.</br> `)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </p> `)
	SendMail(GetMailPolicy(policy, "paga la tua polizza"+" n° "+policy.CodeCompany, true, policy.PayUrl, "Paga la tua polizza", message, false, nil))

}
func SendMailSign(policy models.Policy) {
	var message []string
	message = append(message, `Puoi ora completare la sottoscrizione della tua polizza.`)
	message = append(message, `<br>Clicca sul bottone sotto per accedere alla procedura semplice e guidata di firma elettronica avanzata tramite 
	utilizzo di un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato.`)
	message = append(message, `<br>Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.`)
	message = append(message, `<br>Ultimata la procedura di firma potrai procedere al pagamento.<br>`)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </br> `)
	SendMail(GetMailPolicy(policy, "firma la tua polizza"+" n° "+policy.CodeCompany, true, policy.SignUrl, "Firma la tua polizza", message, false, nil))
}

func SendMailContract(policy models.Policy, at *[]Attachment) {
	var message []string
	message = append(message, `<p>ti confermiamo che la protezione offerta dalla tua polizza è ora attiva. 
	in allegato trovi la documentazione firmata tramite l’utilizzo della Firma Elettronica. Salva e conserva i documenti con cura, ti serviranno in caso di sinistro.
	Ti consigliamo di scaricare l’App di Wopta dagli store tramite il comodo QR code che trovi nell area sottostante per accedere 
	alla tua area riservata nella quale troverai tutte le informazioni sulle polizze in tuo possesso e altri servizi a te riservati.</p>
	Puoi usare anche questi canali per effettuare una denuncia di sinistro e verificare lo stato delle tue polizze e dei pagamenti.
Seguici su nostri canali social o sul sito e scopri le iniziative a te riservate.
	`)
	SendMail(GetMailPolicy(policy, "contratto"+" n° "+policy.CodeCompany, false, "", "", message, true, at))

}
