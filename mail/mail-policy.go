package mail

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	proposalTemplateType = "proposal"
	payTemplateType      = "pay"
	signTemplateType     = "sign"
	emittedTemplateType  = "emitted"
)

func GetMailPolicy(policy *models.Policy, subject string, isLink bool, cc, link, linkLabel, message string, isAttachment bool, at *[]Attachment) MailRequest {
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
	obj.Message = message
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
		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
		bodyData   = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, proposalTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	SendMail(
		GetMailPolicy(
			&policy,
			"Documenti precontrattuali",
			true,
			cc,
			link,
			"Leggi documentazione",
			messageBody,
			false,
			nil,
		),
	)
}

func SendMailPay(policy models.Policy) {
	var (
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, payTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	SendMail(
		GetMailPolicy(
			&policy,
			"Paga la tua polizza"+" n° "+policy.CodeCompany,
			true,
			cc,
			policy.PayUrl,
			"Paga la tua polizza",
			messageBody,
			false,
			nil,
		),
	)
}

func SendMailSign(policy models.Policy) {
	var (
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, signTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	SendMail(
		GetMailPolicy(
			&policy,
			"Firma la tua polizza"+" n° "+policy.CodeCompany,
			true,
			cc,
			policy.SignUrl,
			"Firma la tua polizza",
			messageBody,
			false,
			nil,
		),
	)
}

func SendMailContract(policy models.Policy, at *[]Attachment) {
	var (
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, emittedTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

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
			&policy,
			"Contratto"+" n° "+policy.CodeCompany,
			false,
			cc,
			"",
			"",
			messageBody,
			true,
			at,
		),
	)
}
