package mail

import (
	"bytes"
	"text/template"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func getChannel(policy models.Policy) string {
	if policy.AgentUid != "" {
		return models.AgentChannel
	}
	if policy.AgencyUid != "" {
		return models.AgencyChannel
	}
	return models.ECommerceChannel
}

func setBodyDataAndGetCC(channel string, policy models.Policy, bodyData *BodyData) string {
	var cc string

	switch channel {
	case models.AgentChannel:
		agent, err := models.GetAgentByAuthId(policy.AgentUid)
		lib.CheckError(err)
		cc = agent.Mail
		setAgentBodyData(*agent, bodyData)
	case models.AgencyChannel:
		agency, err := models.GetAgencyByAuthId(policy.AgencyUid)
		lib.CheckError(err)
		cc = agency.Email
		setAgencyBodyData(*agency, bodyData)
	}

	setProductBodyData(policy, bodyData)

	setContractorBodyData(policy, bodyData)

	return cc
}

func setContractorBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = policy.Contractor.Name
	bodyData.ContractorSurname = policy.Contractor.Surname
}

func setAgentBodyData(agent models.Agent, bodyData *BodyData) {
	bodyData.AgentName = agent.Name
	bodyData.AgentSurname = agent.Surname
}

func setAgencyBodyData(agency models.Agency, bodyData *BodyData) {
	bodyData.AgencyName = agency.Name
}

func setProductBodyData(policy models.Policy, bodyData *BodyData) {
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
