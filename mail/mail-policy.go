package mail

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetMailPolicy(
	policy models.Policy,
	subject string,
	isLink bool,
	link string,
	linkLabel string,
	lines []string,
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
	obj.Message = `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">Ciao <strong>` + policy.Contractor.Name + " " + policy.Contractor.Surname + `</strong>,</p> `
	for _, line := range lines {
		obj.Message = obj.Message + `<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">` + line + `</p>`
	}

	obj.Message = obj.Message + ` 
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">Se hai bisogno di ulteriore supporto, non scrivere a questo indirizzo email, puoi compilare il <a class="button" href='` + linkForm + ` '>Form </a> oppure scrivere alla mail e verrai contattato da un nostro esperto.</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px"><br></p><p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#000000;font-size:14px">A presto,</p>
	<p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:17px;color:#e50075;font-size:14px"><strong>Anna</strong> di Wopta Assicurazioni</p> `
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

func SendMailProposal(policy models.Policy) {
	var (
		name       string
		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
		message    = make([]string, 0, 3)
	)

	message = append(
		message,
		`	<br>Grazie per aver compilato, la richiesta di preventivo per una polizza <b>`+name+`</b>.
	<br><br>Per poter valutare completamente la soluzione che sceglierai, ti alleghiamo tutti i documenti che ti consentiranno di prendere una decisione pienamente consapevole ed informata.<br> `,
	)
	message = append(
		message,
		`<br>Se ci fosse la necessità di richiederti alcune informazioni aggiuntive, ti contatteremo ai recapiti da te forniti.`,
	)
	message = append(
		message,
		`<br>Prima della sottoscrizione, leggi quanto trovi in questa mail, la <b>documentazione precontrattuale</b> che,
	 per trasparenza e tua adeguata informazione, trovi sempre disponibile cliccando sul bottone sottostante.`,
	)

	SendMail(
		GetMailPolicy(
			policy,
			"Documenti precontrattuali",
			true,
			link,
			"Leggi documentazione",
			message,
			false,
			nil,
		),
	)
}

func SendMailPay(policy models.Policy) {
	var message []string
	message = append(
		message,
		`<p>hai firmato correttamente la polizza. Sei più vicino a sentirti più protetto.</br> `,
	)
	message = append(
		message,
		`Ti invitiamo ora ad accedere a cliccare sul bottone sottostante per perfezionare il pagamento.</p> `,
	)
	message = append(
		message,
		`<p>Infatti senza pagamento la polizza non è attiva e, solo a pagamento avvenuto, ti invieremo una mail in cui trovi tutti i documenti contrattuali completi.</br> `,
	)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </p> `)
	SendMail(
		GetMailPolicy(
			policy,
			"Paga la tua polizza"+" n° "+policy.CodeCompany,
			true,
			policy.PayUrl,
			"Paga la tua polizza",
			message,
			false,
			nil,
		),
	)
}

func SendMailSign(policy models.Policy) {
	var message []string
	message = append(message, `Puoi ora completare la sottoscrizione della tua polizza.`)
	message = append(
		message,
		`<br>Clicca sul bottone sottostante per accedere alla procedura semplice e guidata di firma elettronica avanzata tramite 
	utilizzo di un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato.`,
	)
	message = append(
		message,
		`<br>Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.`,
	)
	message = append(
		message,
		`<br>Ultimata la procedura di firma potrai procedere al pagamento.<br>`,
	)
	message = append(message, `Qualora tu abbia già provveduto, ignora questa comunicazione </br> `)
	SendMail(
		GetMailPolicy(
			policy,
			"Firma la tua polizza"+" n° "+policy.CodeCompany,
			true,
			policy.SignUrl,
			"Firma la tua polizza",
			message,
			false,
			nil,
		),
	)
}

func SendMailContract(policy models.Policy, at *[]Attachment) {
	var message []string

	message = append(
		message,
		`<p>ti confermiamo che la protezione offerta dalla tua polizza è ora attiva. 
	in allegato trovi la documentazione firmata tramite l’utilizzo della Firma Elettronica. Salva e conserva i documenti con cura, ti serviranno in caso di sinistro.
	Ti consigliamo di scaricare l’App di Wopta dagli store tramite il comodo QR code che trovi nell area sottostante per accedere 
	alla tua area riservata nella quale troverai tutte le informazioni sulle polizze in tuo possesso e altri servizi a te riservati.</p>
	Puoi usare anche questi canali per effettuare una denuncia di sinistro e verificare lo stato delle tue polizze e dei pagamenti.
Seguici su nostri canali social o sul sito e scopri le iniziative a te riservate.
	`,
	)

	// retrocompatibility - the new use extracts the contract from the policy
	if at == nil {
		var contractbyte []byte

		filepath := fmt.Sprintf("assets/users/%s/contract_%s.pdf", policy.Contractor.Uid, policy.Uid)
		contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath)
		lib.CheckError(err)

		filenameParts := []string{policy.Contractor.Name, policy.Contractor.Surname, policy.NameDesc, "contratto.pdf"}
		filename := strings.Join(filenameParts, "_")
		filename = strings.ReplaceAll(filename, " ", "_")
		at = &[]Attachment{{
			Byte:        base64.StdEncoding.EncodeToString(contractbyte),
			ContentType: "application/pdf",
			Name:        filename,
		}}
	}

	SendMail(
		GetMailPolicy(
			policy,
			"Contratto"+" n° "+policy.CodeCompany,
			false,
			"",
			"",
			message,
			true,
			at,
		),
	)
}
