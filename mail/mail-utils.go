package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
)

func setBodyData(policy models.Policy, bodyData *BodyData) {
	setProductBodyData(policy, bodyData)

	setContractorBodyData(policy, bodyData)

	if policy.IsReserved {
		setPolicyReservedBodyData(policy, bodyData)
	}

	node := network.GetNetworkNodeByUid(policy.ProducerUid)

	if node != nil {
		setNetworkNodeBodyData(node, bodyData)
	}
}

func setNetworkNodeBodyData(node *models.NetworkNode, bodyData *BodyData) {
	bodyData.NetworkNodeName = node.GetName()
	bodyData.NetworkNodeEmail = node.Mail
}

func setContractorBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = policy.Contractor.Name
	bodyData.ContractorSurname = policy.Contractor.Surname
	bodyData.ContractorFiscalCode = policy.Contractor.FiscalCode
}

func setProductBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ProductForm = "https://www.wopta.it/it/"
	bodyData.ProductSlug = policy.Name
	bodyData.SignUrl = policy.SignUrl
	bodyData.PayUrl = policy.PayUrl
	bodyData.PaymentMode = policy.PaymentMode

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
	}

	bodyData.InformationSetsUrl = fmt.Sprintf(
		"https://storage.googleapis.com/documents-public-dev/information-sets/%s/%s/Precontrattuale.pdf",
		policy.Name, policy.ProductVersion)
}

func setPolicyReservedBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ProposalNumber = policy.ProposalNumber
	if policy.ReservedInfo != nil && len(policy.ReservedInfo.RequiredExams) > 0 {
		bodyData.ExtraContent = policy.ReservedInfo.RequiredExams
	}
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
	return Address{
		Name:    policy.Contractor.Name + " " + policy.Contractor.Surname,
		Address: policy.Contractor.Mail,
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

func getMailAttachments(policy models.Policy, attachmentNames []string) []Attachment {
	var (
		at     []Attachment
		rawDoc []byte
		err    error
	)

	if policy.Attachments == nil || len(*policy.Attachments) == 0 {
		log.Println("[getMailAttachments] policy has no attachment")
		return at
	}

	at = make([]Attachment, 0)

	for _, attachment := range *policy.Attachments {
		if lib.SliceContains(attachmentNames, attachment.Name) {
			if strings.HasPrefix(attachment.Link, "gs://") {
				rawDoc, err = lib.ReadFileFromGoogleStorage(attachment.Link)
			} else {
				rawDoc, err = lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), attachment.Link)
			}
			if err != nil {
				log.Printf("[getMailAttachments] error reading document %s from google storage: %s", attachment.Name, err.Error())
				return nil
			}
			attachment.Byte = base64.StdEncoding.EncodeToString(rawDoc)

			at = append(at, Attachment{
				Name:        strings.ReplaceAll(attachment.FileName, "_", " "),
				Link:        attachment.Link,
				Byte:        attachment.Byte,
				FileName:    attachment.FileName,
				ContentType: "application/pdf",
			})
		}
	}

	return at
}
