package mail

import (
	"bytes"
	"fmt"

	"github.com/wopta/goworkspace/models"
)

func GetMailPolicy(
	policy *models.Policy,
	subject string,
	isLink bool,
	cc string,
	link string,
	linkLabel string,
	lines string,
	isAttachment bool,
	at *[]Attachment,
) MailRequest {
	var (
		name     string
		obj      MailRequest
		linkForm = "https://www.wopta.it/it/"
	)

	switch policy.Name {
	case "pmi":
		name = "Artigiani & Imprese"
		linkForm += "multi-rischio/"
	case "persona":
		name = "Persona"
		linkForm += "infortunio/"
	case "life":
		name = "Vita"
		linkForm += "vita/"
	case "gap":
		name = "GAP"
		// TODO: No page yet
	}

	obj.From = "anna@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Cc = cc
	obj.Message = lines
	obj.Title = "Wopta per te " + name
	obj.Subject = obj.Title + " " + subject
	obj.SubTitle = subject
	obj.IsHtml = true
	obj.IsAttachment = isAttachment
	obj.IsLink = isLink
	if isLink {
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

// func SendMailProposal(policy models.Policy) {
// 	var (
// 		name       string
// 		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
// 		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
// 		message    = make([]string, 0, 3)
// 	)

// 	message = append(
// 		message,
// 		`	<br>Grazie per aver compilato, la richiesta di preventivo per una polizza <b>`+name+`</b>.
// 	<br><br>Per poter valutare completamente la soluzione che sceglierai, ti alleghiamo tutti i documenti che ti consentiranno di prendere una decisione pienamente consapevole ed informata.<br> `,
// 	)
// 	message = append(
// 		message,
// 		`<br>Se ci fosse la necessità di richiederti alcune informazioni aggiuntive, ti contatteremo ai recapiti da te forniti.`,
// 	)
// 	message = append(
// 		message,
// 		`<br>Prima della sottoscrizione, leggi quanto trovi in questa mail, la <b>documentazione precontrattuale</b> che,
// 	 per trasparenza e tua adeguata informazione, trovi sempre disponibile cliccando sul bottone sottostante.`,
// 	)

// 	SendMail(
// 		GetMailPolicy(
// 			policy,
// 			"Documenti precontrattuali",
// 			true,
// 			link,
// 			"Leggi documentazione",
// 			message,
// 			false,
// 			nil,
// 		),
// 	)
// }

func SendMailProposal(policy *models.Policy) {
	var (
		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
		tpl        bytes.Buffer
	)

	bodyData := BodyData{}

	cc := SetBodyDataAndGetCC(policy, &bodyData)

	templateFile := GetTemplateByChannel(policy, "pay")

	FillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			policy,
			"Documenti precontrattuali",
			true,
			cc,
			link,
			"Leggi documentazione",
			tpl.String(),
			false,
			nil,
		),
	)
}

func SendMailPay(policy *models.Policy) {
	bodyData := BodyData{}
	var tpl bytes.Buffer

	cc := SetBodyDataAndGetCC(policy, &bodyData)

	templateFile := GetTemplateByChannel(policy, "pay")

	FillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			policy,
			"Paga la tua polizza"+" n° "+policy.CodeCompany,
			true,
			cc,
			policy.PayUrl,
			"Paga la tua polizza",
			tpl.String(),
			false,
			nil,
		),
	)
}

// func SendMailPay(policy models.Policy) {
// 	var message []string
// 	message = append(
// 		message,
// 		`<p>hai firmato correttamente la polizza. Sei più vicino a sentirti più protetto.</br> `,
// 	)
// 	message = append(
// 		message,
// 		`Ti invitiamo ora ad accedere a cliccare sul bottone sottostante per perfezionare il pagamento.</p> `,
// 	)
// 	message = append(
// 		message,
// 		`<p>Infatti senza pagamento la polizza non è attiva e, solo a pagamento avvenuto, ti invieremo una mail in cui trovi tutti i documenti contrattuali completi.</br> `,
// 	)
// 	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </p> `)
// 	SendMail(
// 		GetMailPolicy(
// 			policy,
// 			"Paga la tua polizza"+" n° "+policy.CodeCompany,
// 			true,
// 			policy.PayUrl,
// 			"Paga la tua polizza",
// 			message,
// 			false,
// 			nil,
// 		),
// 	)
// }

func SendMailSign(policy *models.Policy) {
	bodyData := BodyData{}
	var tpl bytes.Buffer

	cc := SetBodyDataAndGetCC(policy, &bodyData)

	templateFile := GetTemplateByChannel(policy, "sign")

	FillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			policy,
			"Firma la tua polizza"+" n° "+policy.CodeCompany,
			true,
			cc,
			policy.SignUrl,
			"Firma la tua polizza",
			tpl.String(),
			false,
			nil,
		),
	)
}

// func SendMailSign(policy models.Policy) {
// 	var message []string
// 	message = append(message, `Puoi ora completare la sottoscrizione della tua polizza.`)
// 	message = append(
// 		message,
// 		`<br>Clicca sul bottone sottostante per accedere alla procedura semplice e guidata di firma elettronica avanzata tramite
// 	utilizzo di un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato.`,
// 	)
// 	message = append(
// 		message,
// 		`<br>Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.`,
// 	)
// 	message = append(
// 		message,
// 		`<br>Ultimata la procedura di firma potrai procedere al pagamento.<br>`,
// 	)
// 	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </br> `)
// 	SendMail(
// 		GetMailPolicy(
// 			policy,
// 			"Firma la tua polizza"+" n° "+policy.CodeCompany,
// 			true,
// 			policy.SignUrl,
// 			"Firma la tua polizza",
// 			message,
// 			false,
// 			nil,
// 		),
// 	)
// }

// func SendMailContract(policy models.Policy, at *[]Attachment) {
// 	var message []string

// 	message = append(
// 		message,
// 		`<p>ti confermiamo che la protezione offerta dalla tua polizza è ora attiva.
// 	in allegato trovi la documentazione firmata tramite l’utilizzo della Firma Elettronica. Salva e conserva i documenti con cura, ti serviranno in caso di sinistro.
// 	Ti consigliamo di scaricare l’App di Wopta dagli store tramite il comodo QR code che trovi nell area sottostante per accedere
// 	alla tua area riservata nella quale troverai tutte le informazioni sulle polizze in tuo possesso e altri servizi a te riservati.</p>
// 	Puoi usare anche questi canali per effettuare una denuncia di sinistro e verificare lo stato delle tue polizze e dei pagamenti.
// Seguici su nostri canali social o sul sito e scopri le iniziative a te riservate.
// 	`,
// 	)
// 	SendMail(
// 		GetMailPolicy(
// 			policy,
// 			"Contratto"+" n° "+policy.CodeCompany,
// 			false,
// 			"",
// 			"",
// 			message,
// 			true,
// 			at,
// 		),
// 	)
// }

func SendMailContract(policy *models.Policy, at *[]Attachment) {

	bodyData := BodyData{}
	var tpl bytes.Buffer

	cc := SetBodyDataAndGetCC(policy, &bodyData)

	templateFile := GetTemplateByChannel(policy, "sign")

	FillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			policy,
			"Contratto"+" n° "+policy.CodeCompany,
			false,
			cc,
			"",
			"",
			tpl.String(),
			true,
			at,
		),
	)
}
