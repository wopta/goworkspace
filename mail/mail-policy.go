package mail

import (
	"bytes"
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

func GetMailPolicy(policy *models.Policy, subject string, isLink bool, cc string, link string, linkLabel string, lines string, isAttachment bool, at *[]Attachment) MailRequest {
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

func SendMailProposal(policy models.Policy) {
	var (
		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
		tpl        bytes.Buffer
		bodyData   = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := getTemplateByChannel(channel, proposalTemplateType)

	fillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			&policy,
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

func SendMailPay(policy models.Policy) {
	var (
		tpl      bytes.Buffer
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := getTemplateByChannel(channel, payTemplateType)

	fillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			&policy,
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

func SendMailSign(policy models.Policy) {
	var (
		tpl      bytes.Buffer
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := getTemplateByChannel(channel, signTemplateType)

	fillTemplate(templateFile, &bodyData, &tpl)

	SendMail(
		GetMailPolicy(
			&policy,
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

func SendMailContract(policy models.Policy, at *[]Attachment) {
	var (
		tpl      bytes.Buffer
		bodyData = BodyData{}
	)

	channel := getChannel(policy)

	cc := setBodyDataAndGetCC(channel, policy, &bodyData)

	templateFile := getTemplateByChannel(channel, emittedTemplateType)

	fillTemplate(templateFile, &bodyData, &tpl)

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
			tpl.String(),
			true,
			at,
		),
	)
}
