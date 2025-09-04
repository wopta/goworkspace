package mail

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	leadTemplateType             = "lead"
	proposalTemplateType         = "proposal"
	payTemplateType              = "pay"
	signTemplateType             = "sign"
	contractTemplateType         = "contract"
	reservedTemplateType         = "reserved"
	reservedApprovedTemplateType = "approved"
	reservedRejectedTemplateType = "rejected"
	renewDraftTemplateType       = "renew-draft"
)

func SendMailLead(policy models.Policy, from, to, cc Address, flowName string, attachmentNames []string) {

	messageBody, _ := getTemplateEmail(flowName, leadTemplateType, policy)

	title := policy.NameDesc
	subtitle := "Documenti precontrattuali"
	subject := fmt.Sprintf("%s %s", title, subtitle)

	at := getMailAttachments(policy, attachmentNames)

	SendMail(MailRequest{
		FromAddress:  from,
		To:           []string{to.Address},
		Cc:           cc.Address,
		Message:      messageBody,
		Title:        title,
		SubTitle:     subtitle,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: len(at) > 0,
		Attachments:  &at,
	})
}

func SendMailPay(policy models.Policy, from, to, cc Address, flowName string) {

	messageBody, _ := getTemplateEmail(flowName, payTemplateType, policy)

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
	})
}

func SendMailSign(policy models.Policy, from, to, cc Address, flowName string) {

	messageBody, _ := getTemplateEmail(flowName, signTemplateType, policy)

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
	})
}

func SendMailContract(policy models.Policy, at *[]models.Attachment, from, to, cc Address, flowName string) error {
	var (
		bcc string
	)

	messageBody, _ := getTemplateEmail(flowName, contractTemplateType, policy)
	// retrocompatibility - the new use extracts the contract from the policy
	if at == nil {
		var contractbyte []byte

		filepath := fmt.Sprintf("assets/users/%s/"+models.ContractDocumentFormat, policy.Contractor.Uid, policy.NameDesc, policy.CodeCompany)
		contractbyte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath)
		if err != nil {
			log.Error(err)
			return err
		}

		filename := fmt.Sprintf(models.ContractDocumentFormat, policy.NameDesc,
			policy.CodeCompany)
		at = &[]models.Attachment{{
			Byte:        base64.StdEncoding.EncodeToString(contractbyte),
			ContentType: lib.GetContentType("pdf"),
			FileName:    filename,
			Name:        strings.ReplaceAll(filename, "_", " "),
		}}
	}
	var attachmentToSend []models.Attachment
	for _, attach := range *at {
		if !strings.Contains(attach.Name, models.ContractAttachmentName) && !strings.Contains(attach.Name, models.PreContractAttachmentName) {
			continue
		}
		attachmentToSend = append(attachmentToSend, attach)
	}
	title := policy.NameDesc
	subtitle := fmt.Sprintf("Contratto n° %s", policy.CodeCompany)
	subject := fmt.Sprintf("%s %s", title, subtitle)

	if policy.HasPrivacyConsens() {
		bcc = os.Getenv("BCC_CONTRACT_EMAIL")
	}

	SendMail(MailRequest{
		FromAddress:  from,
		To:           []string{to.Address},
		Cc:           cc.Address,
		Bcc:          bcc,
		Message:      messageBody,
		Title:        policy.NameDesc,
		SubTitle:     subtitle,
		Subject:      subject,
		IsHtml:       true,
		IsApp:        true,
		IsAttachment: true,
		Attachments:  &attachmentToSend,
	})
	return nil
}

func SendMailReserved(policy models.Policy, from, to, cc Address, flowName string, attachmentNames []string) {
	var (
		at       []models.Attachment
		bodyData BodyData
		rawDoc   []byte
		err      error
	)
	log.AddPrefix("sendMailReserved")
	defer log.PopPrefix()

	bodyData = getBodyData(policy)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s-%s.html", flowName, reservedTemplateType, policy.Name))

	messageBody := fillTemplate(templateFile, &bodyData)

	for _, attachment := range policy.ReservedInfo.Documents {
		if attachment.Byte == "" {
			rawDoc, err = lib.ReadFileFromGoogleStorageEitherGsOrNot(attachment.Link)
			if err != nil {
				log.ErrorF("error reading document %s from google storage: %s", attachment.Name, err.Error())
				return
			}
			attachment.Byte = base64.StdEncoding.EncodeToString(rawDoc)
		}

		at = append(at, models.Attachment{
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
		bodyData     BodyData
		templateType string
	)

	if policy.Status == models.PolicyStatusApproved {
		templateType = reservedApprovedTemplateType
	} else {
		templateType = reservedRejectedTemplateType
	}

	messageBody, _ := getTemplateEmail(flowName, templateType, policy)

	bodyData = getBodyData(policy)

	title := fmt.Sprintf(
		"Wopta per te %s - Proposta %d di %s",
		bodyData.ProductName,
		policy.ProposalNumber,
		bodyData.ContractorName,
	)
	// TODO: handle multiple products reserved subtitle
	subtitle := "Esito valutazione medica assuntiva"
	subject := fmt.Sprintf("%s - %s", title, subtitle)

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
		at []models.Attachment
	)

	messageBody, _ := getTemplateEmail(flowName, proposalTemplateType, policy)
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
	})
}

func SendMailRenewDraft(policy models.Policy, from, to, cc Address, flowName string, hasMandate bool) {

	messageBody, _ := getTemplateEmail(flowName, renewDraftTemplateType, policy)

	bodyData := getPolicyRenewDraftBodyData(policy)

	title := policy.NameDesc
	subtitle := fmt.Sprintf("La tua polizza n° %s si rinnova il %s, provvedi al pagamento.", policy.CodeCompany,
		bodyData.RenewDate)
	if hasMandate {
		subtitle = fmt.Sprintf("La tua polizza n° %s si rinnova il %s, pagamento senza pensieri.",
			policy.CodeCompany, bodyData.RenewDate)
	}
	subject := fmt.Sprintf("%s - %s", title, subtitle)

	SendMail(MailRequest{
		FromName:    from.Name,
		FromAddress: from,
		To:          []string{to.Address},
		Cc:          cc.Address,
		Message:     messageBody,
		Title:       title,
		SubTitle:    subtitle,
		Subject:     subject,
		IsHtml:      true,
		IsApp:       true,
		Policy:      policy.Uid,
	})
}
