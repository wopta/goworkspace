package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/dustin/go-humanize"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
)

func getBodyData(policy models.Policy) BodyData {
	var bodyData BodyData

	setProductBodyData(policy, &bodyData)

	setContractorBodyData(policy, &bodyData)

	if policy.IsReserved {
		setPolicyReservedBodyData(policy, &bodyData)
	}

	node := network.GetNetworkNodeByUid(policy.ProducerUid)

	if node != nil {
		setNetworkNodeBodyData(node, &bodyData)
	}

	return bodyData
}

func setNetworkNodeBodyData(node *models.NetworkNode, bodyData *BodyData) {
	bodyData.NetworkNodeName = node.GetName()
	bodyData.NetworkNodeEmail = node.Mail
}

func setContractorBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = lib.TrimSpace(fmt.Sprintf("%s %s", policy.Contractor.Name, policy.Contractor.Surname))
	bodyData.ContractorFiscalCode = policy.Contractor.FiscalCode
}

func setProductBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ProductForm = "https://www.wopta.it/it/"
	bodyData.ProductSlug = policy.Name
	bodyData.SignUrl = policy.SignUrl
	bodyData.PayUrl = policy.PayUrl
	bodyData.PaymentMode = policy.PaymentMode
	bodyData.ProposalNumber = policy.ProposalNumber
	bodyData.PolicyUid = policy.Uid

	switch policy.Name {
	case models.PmiProduct:
		bodyData.ProductName = "Artigiani & Imprese"
		bodyData.ProductForm += "multi-rischio#contact-us"
	case models.PersonaProduct:
		bodyData.ProductName = "Persona"
		bodyData.ProductForm += "infortunio#contact-us"
	case models.LifeProduct:
		bodyData.ProductName = "Vita"
		bodyData.ProductForm += "vita#contact-us"
	case models.GapProduct:
		bodyData.ProductName = "Auto Valore Protetto"
		bodyData.ProductForm = "gap#contact-us"
	case models.CatNatProduct:
		bodyData.ProductName = "Catastrofali azienda"
		bodyData.ProductForm = "cat-nat#contact-us"
	}
	link, _ := lib.GetLastVersionSetInformativo(policy.Name, policy.ProductVersion)
	bodyData.InformationSetsUrl = fmt.Sprint(lib.BaseStorageGoogleUrl, link)

}

func setPolicyReservedBodyData(policy models.Policy, bodyData *BodyData) {
	if policy.ReservedInfo != nil && len(policy.ReservedInfo.RequiredExams) > 0 {
		bodyData.ExtraContent = policy.ReservedInfo.RequiredExams
	}
}

func getPolicyRenewDraftBodyData(policy models.Policy) BodyData {
	priceGross := policy.PriceGross
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		priceGross = policy.PriceGrossMonthly
	}

	bodyData := getBodyData(policy)
	bodyData.HasMandate = policy.HasMandate
	bodyData.PriceGross = humanize.FormatFloat("#.###,##", priceGross)
	bodyData.RenewDate = policy.StartDate.AddDate(policy.Annuity, 0, 0).Format("02/01/2006")

	return bodyData
}

func fillTemplate(htmlTemplate []byte, bodyData *BodyData) string {
	tpl := new(bytes.Buffer)
	tmplt := template.New("htmlTemplate")
	tmplt, err := tmplt.Parse(string(htmlTemplate))
	lib.CheckError(err)
	err = tmplt.Execute(tpl, bodyData)
	lib.CheckError(err)
	return tpl.String()
}

func fillTemplateV2(htmlTemplate []byte, bodyData *BodyData) (string, error) {
	tpl := new(bytes.Buffer)
	tmplt := template.New("htmlTemplate")
	tmplt, err := tmplt.Parse(string(htmlTemplate))
	if err != nil {
		return "", err
	}
	err = tmplt.Execute(tpl, bodyData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func GetEmailByChannel(policy *models.Policy) Address {
	var address Address

	switch policy.Channel {
	case models.AgentChannel:
		return GetAgentEmail(policy)
	case models.AgencyChannel:
		return GetAgencyEmail(policy)
	case models.ECommerceChannel:
		return GetContractorEmail(policy)
	}

	return address
}

func GetContractorEmail(policy *models.Policy) Address {
	name := policy.Contractor.Name + " " + policy.Contractor.Surname
	address := policy.Contractor.Mail

	if policy.Contractor.Type == models.UserLegalEntity {
		// TODO: handle multiple target signatures
		for _, c := range *policy.Contractors {
			if c.IsSignatory {
				name = c.Name + " " + c.Surname
				address = c.Mail
				break
			}
		}
	}

	return Address{
		Name:    name,
		Address: address,
	}
}

func GetAgencyEmail(policy *models.Policy) Address {
	if policy.AgencyUid == "" {
		return Address{}
	}
	agency, err := models.GetAgencyByAuthId(policy.AgencyUid)
	lib.CheckError(err)
	return Address{
		Name:    agency.Name,
		Address: agency.Email,
	}
}

func GetAgentEmail(policy *models.Policy) Address {
	if policy.AgentUid == "" {
		return Address{}
	}
	agent, err := models.GetAgentByAuthId(policy.AgentUid)
	lib.CheckError(err)
	return Address{
		Name:    agent.Name + " " + agent.Surname,
		Address: agent.Mail,
	}
}

func GetNetworkNodeEmail(networkNode *models.NetworkNode) Address {
	var address Address = Address{
		Address: networkNode.Mail,
	}

	switch networkNode.Type {
	case models.AgentNetworkNodeType:
		address.Name = networkNode.Agent.Name + " " + networkNode.Agent.Surname
	case models.AgencyNetworkNodeType:
		address.Name = networkNode.Agency.Name
	case models.BrokerNetworkNodeType:
		address.Name = networkNode.Broker.Name
	case models.AreaManagerNetworkNodeType:
		address.Name = networkNode.AreaManager.Name
	case models.PartnershipNetworkNodeType:
		address.Name = networkNode.Partnership.Name
	}

	return address
}

func getMailAttachments(policy models.Policy, attachmentNames []string) []models.Attachment {
	var (
		at     []models.Attachment
		rawDoc []byte
		err    error
	)
	log.AddPrefix("getMailAttachments")
	defer log.PopPrefix()
	if policy.Attachments == nil || len(*policy.Attachments) == 0 {
		log.Println("policy has no attachment")
		return at
	}

	at = make([]models.Attachment, 0)
	for _, attachment := range *policy.Attachments {
		if lib.SliceContains(attachmentNames, attachment.Name) {
			rawDoc, err = lib.ReadFileFromGoogleStorageEitherGsOrNot(attachment.Link)
			if err != nil {
				log.ErrorF("error reading document %s from google storage: %s", attachment.Name, err.Error())
				return nil
			}
			attachment.Byte = base64.StdEncoding.EncodeToString(rawDoc)

			at = append(at, models.Attachment{
				Name:        strings.ReplaceAll(attachment.FileName, "_", " "),
				Link:        attachment.Link,
				Byte:        attachment.Byte,
				FileName:    attachment.FileName,
				ContentType: lib.GetContentType("pdf"),
			})
		}
	}

	return at
}

func getTemplateEmail(flowName, templateType string, policy models.Policy) (string, error) {
	bodyData := getBodyData(policy)

	templateFile := lib.GetFilesByEnv(fmt.Sprintf("mail/%s/%s.html", flowName, templateType))
	template := fillTemplate(templateFile, &bodyData)
	return template, nil
}
