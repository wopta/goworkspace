package mail

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	leadTemplateType             = "lead"
	proposalTemplateType         = "proposal"
	payTemplateType              = "pay"
	signTemplateType             = "sign"
	emittedTemplateType          = "emitted"
	reservedTemplateType         = "reserved"
	reservedApprovedTemplateType = "approved"
	reservedRejectedTemplateType = "rejected"
	linkFormat                   = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
)

func SendMailLead(policy models.Policy, from, to, cc Address, flowName string) {
	var (
		bodyData = BodyData{}
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, leadTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	title := policy.NameDesc
	subtitle := "Documenti precontrattuali"
	subject := fmt.Sprintf("%s %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress: from,
		To:          []string{to.Address},
		Cc:          cc.Address,
		Message:     messageBody,
		Title:       title,
		SubTitle:    subtitle,
		Subject:     subject,
		IsHtml:      true,
	})
}

func SendMailPay(policy models.Policy, from, to, cc Address, flowName string) {
	var (
		bodyData = BodyData{}
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, payTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	title := policy.NameDesc
	subtitle := fmt.Sprintf("Paga la tua polizza n° %s", policy.CodeCompany)
	subject := fmt.Sprintf("%s %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress: from,
		To:          []string{to.Address},
		Cc:          cc.Address,
		Message:     messageBody,
		Title:       title,
		SubTitle:    subtitle,
		Subject:     subject,
		IsHtml:      true,
		IsLink:      true,
		Link:        policy.PayUrl,
		LinkLabel:   "Paga la tua polizza",
	})
}

func SendMailSign(policy models.Policy, from, to, cc Address, flowName string) {
	var (
		bodyData = BodyData{}
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, signTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	title := policy.NameDesc
	subtitle := fmt.Sprintf("Firma la tua polizza n° %s", policy.CodeCompany)
	subject := fmt.Sprintf("%s %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress: from,
		To:          []string{to.Address},
		Cc:          cc.Address,
		Message:     messageBody,
		Title:       title,
		SubTitle:    subtitle,
		Subject:     subject,
		IsHtml:      true,
		IsLink:      true,
		Link:        policy.SignUrl,
		LinkLabel:   "Firma la tua polizza",
	})
}

func SendMailContract(policy models.Policy, at *[]Attachment, from, to, cc Address, flowName string) {
	var (
		bodyData = BodyData{}
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, emittedTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	// retrocompatibility - the new use extracts the contract from the policy
	if at == nil {
		var contractbyte []byte

		filepath := fmt.Sprintf("assets/users/%s/"+models.ContractDocumentFormat, policy.Contractor.Uid, policy.NameDesc, policy.CodeCompany)
		contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath)
		lib.CheckError(err)

		filename := fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc,
			policy.CodeCompany)
		at = &[]Attachment{{
			Byte:        base64.StdEncoding.EncodeToString(contractbyte),
			ContentType: "application/pdf",
			FileName:    filename,
			Name:        strings.ReplaceAll(filename, "_", " "),
		}}
	}

	title := policy.NameDesc
	subtitle := fmt.Sprintf("Contratto n° %s", policy.CodeCompany)
	subject := fmt.Sprintf("%s %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress:  from,
		To:           []string{to.Address},
		Cc:           cc.Address,
		Message:      messageBody,
		Title:        policy.NameDesc,
		SubTitle:     subtitle,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: true,
		Attachments:  at,
	})
}

func SendMailReserved(policy models.Policy, from, to, cc Address, flowName string, attachmentNames []string) {
	var (
		at       []Attachment
		bodyData = BodyData{}
		rawDoc   []byte
		err      error
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, reservedTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	for _, attachment := range policy.ReservedInfo.Documents {
		if attachment.Byte == "" {
			if strings.HasPrefix(attachment.Link, "gs://") {
				rawDoc, err = lib.ReadFileFromGoogleStorage(attachment.Link)
			} else {
				rawDoc, err = lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), attachment.Link)
			}
			if err != nil {
				log.Printf("[sendMailReserved] error reading document %s from google storage: %s", attachment.Name, err.Error())
				return
			}
			attachment.Byte = base64.StdEncoding.EncodeToString(rawDoc)
		}

		at = append(at, Attachment{
			Name:        strings.ReplaceAll(fmt.Sprintf("%s", attachment.FileName), "_", " "),
			Link:        attachment.Link,
			Byte:        attachment.Byte,
			FileName:    attachment.FileName,
			MimeType:    attachment.MimeType,
			Url:         attachment.Url,
			ContentType: attachment.ContentType,
		})
	}

	at = append(at, getMailAttachments(policy, attachmentNames)...)

	title := policy.NameDesc
	subtitle := fmt.Sprintf("Documenti Riservato proposta %d", policy.ProposalNumber)
	subject := fmt.Sprintf("%s - %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress:  from,
		To:           []string{to.Address},
		Cc:           cc.Address,
		Message:      messageBody,
		Title:        title,
		SubTitle:     subtitle,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: true,
		Attachments:  &at,
	})
}

func SendMailReservedResult(policy models.Policy, from, to, cc Address, flowName string) {
	var (
		bodyData = BodyData{}
		template string
	)

	if policy.Status == models.PolicyStatusApproved {
		template = reservedApprovedTemplateType
	} else {
		template = reservedRejectedTemplateType
	}

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, template))

	title := fmt.Sprintf(
		"Wopta per te %s - Proposta %d di %s %s",
		bodyData.ProductName,
		policy.ProposalNumber,
		bodyData.ContractorSurname,
		bodyData.ContractorName,
	)
	subtitle := "Esito valutazione medica assuntiva"
	subject := fmt.Sprintf("%s - %s", title, subtitle)

	messageBody := fillTemplate(templateFile, &bodyData)

	SendMail(MailRequest{
		FromAddress: from,
		To:          []string{to.Address},
		Cc:          cc.Address,
		Message:     messageBody,
		Title:       title,
		SubTitle:    subtitle,
		Subject:     subject,
		IsHtml:      true,
	})
}

func SendMailProposal(policy models.Policy, from, to, cc Address, flowName string, attachmentNames []string) {
	var (
		at       []Attachment
		bodyData = BodyData{}
		link     = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
	)

	setBodyData(policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, proposalTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	at = append(at, getMailAttachments(policy, attachmentNames)...)

	title := policy.NameDesc
	subtitle := fmt.Sprintf("Documento Proposta %d", policy.ProposalNumber)
	subject := fmt.Sprintf("%s - %s", title, subtitle)

	SendMail(MailRequest{
		FromAddress:  from,
		To:           []string{to.Address},
		Cc:           cc.Address,
		Message:      messageBody,
		Title:        title,
		SubTitle:     subtitle,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: true,
		Attachments:  &at,
		IsLink:       true,
		Link:         link,
		LinkLabel:    "Leggi documentazione",
	})
}
