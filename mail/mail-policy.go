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
	proposalTemplateType         = "proposal"
	payTemplateType              = "pay"
	signTemplateType             = "sign"
	emittedTemplateType          = "emitted"
	reservedTemplateType         = "reserved"
	reservedApprovedTemplateType = "approved"
	reservedRejectedTemplateType = "rejected"
)

func SendMailLead(policy models.Policy, from, to, cc Address) {
	var (
		linkFormat = "https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf"
		link       = fmt.Sprintf(linkFormat, policy.Name, policy.ProductVersion)
		bodyData   = BodyData{}
		channel    = policy.Channel
	)

	setBodyData(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, proposalTemplateType))

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
		IsLink:      true,
		Link:        link,
		LinkLabel:   "Leggi documentazione",
	})
}

func SendMailPay(policy models.Policy, from, to, cc Address) {
	var (
		bodyData = BodyData{}
		channel  = policy.Channel
	)

	setBodyData(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, payTemplateType))

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

func SendMailSign(policy models.Policy, from, to, cc Address) {
	var (
		bodyData = BodyData{}
		channel  = policy.Channel
	)

	setBodyData(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, signTemplateType))

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

func SendMailContract(policy models.Policy, at *[]Attachment, from, to, cc Address) {
	var (
		bodyData = BodyData{}
		channel  = policy.Channel
	)

	setBodyData(channel, policy, &bodyData)

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

func SendMailReserved(policy models.Policy, from, to, cc Address) {
	var (
		at       []Attachment
		bodyData = BodyData{}
		channel  = policy.Channel
	)

	setBodyData(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, reservedTemplateType))

	messageBody := fillTemplate(templateFile, &bodyData)

	for _, attachment := range policy.ReservedInfo.Documents {
		if attachment.Byte == "" {
			rawDoc, err := lib.ReadFileFromGoogleStorage(attachment.Link)
			if err != nil {
				log.Printf("[sendMailReserved] error reading document %s from google storage: %s", attachment.Name, err.Error())
				return
			}
			attachment.Byte = base64.StdEncoding.EncodeToString(rawDoc)
		}

		at = append(at, Attachment{
			Name:        attachment.Name,
			Link:        attachment.Link,
			Byte:        attachment.Byte,
			FileName:    attachment.FileName,
			MimeType:    attachment.MimeType,
			Url:         attachment.Url,
			ContentType: attachment.ContentType,
		})
	}

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

func SendMailReservedResult(policy models.Policy, from, to, cc Address) {
	var (
		bodyData = BodyData{}
		channel  = policy.Channel
		template string
	)

	if policy.Status == models.PolicyStatusApproved {
		template = reservedApprovedTemplateType
	} else {
		template = reservedRejectedTemplateType
	}

	setBodyData(channel, policy, &bodyData)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", channel, template))

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
