package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func setBodyData(policy models.Policy, bodyData *BodyData) {
	switch policy.Channel {
	case models.AgentChannel:
		agent, err := models.GetAgentByAuthId(policy.AgentUid)
		lib.CheckError(err)
		setAgentBodyData(*agent, bodyData)
	case models.AgencyChannel:
		agency, err := models.GetAgencyByAuthId(policy.AgencyUid)
		lib.CheckError(err)
		setAgencyBodyData(*agency, bodyData)
	}

	setProductBodyData(policy, bodyData)

	setContractorBodyData(policy, bodyData)

	if policy.IsReserved {
		setPolicyReservedBodyData(policy, bodyData)
	}
}

func setContractorBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = policy.Contractor.Name
	bodyData.ContractorSurname = policy.Contractor.Surname
}

func setAgentBodyData(agent models.Agent, bodyData *BodyData) {
	bodyData.AgentName = agent.Name
	bodyData.AgentSurname = agent.Surname
	bodyData.AgentMail = agent.Mail
}

func setAgencyBodyData(agency models.Agency, bodyData *BodyData) {
	bodyData.AgencyName = agency.Name
	bodyData.AgencyMail = agency.Email
}

func setProductBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ProductForm = "https://www.wopta.it/it/"

	switch policy.Name {
	case models.PmiProduct:
		bodyData.ProductName = "Artigiani & Imprese"
		bodyData.ProductForm += "multi-rischio/"
	case models.PersonaProduct:
		bodyData.ProductName = "Persona"
		bodyData.ProductForm += "infortunio/"
	case models.LifeProduct:
		bodyData.ProductName = "Vita"
		bodyData.ProductForm += "vita/"
	case models.GapProduct:
		bodyData.ProductName = "Auto Valore Protetto"
		bodyData.ProductForm = ""
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
