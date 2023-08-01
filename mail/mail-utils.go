package mail

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func getChannel(policy models.Policy) string {
	if policy.AgentUid != "" {
		return models.UserRoleAgent
	}
	if policy.AgencyUid != "" {
		return models.UserRoleAgency
	}
	return "e-commerce"
}

func setBodyDataAndGetCC(channel string, policy models.Policy, bodyData *BodyData) string {
	var cc string

	switch channel {
	case models.UserRoleAgent:
		cc = getAgentBodyData(policy.AgentUid, bodyData)
	case models.UserRoleAgency:
		cc = getAgencyBodyData(policy.AgencyUid, bodyData)
	}

	getProductBodyData(policy, bodyData)

	getContractorBodyData(policy, bodyData)

	return cc
}

func getContractorBodyData(policy models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = policy.Contractor.Name
	bodyData.ContractorSurname = policy.Contractor.Surname
}

func getAgentBodyData(agentUid string, bodyData *BodyData) string {
	agent, err := models.GetAgentByAuthId(agentUid)
	lib.CheckError(err)
	bodyData.AgentName = agent.Name
	bodyData.AgentSurname = agent.Surname
	return agent.Mail
}

func getAgencyBodyData(agencyUid string, bodyData *BodyData) string {
	agency, err := models.GetAgencyByAuthId(agencyUid)
	lib.CheckError(err)
	bodyData.AgencyName = agency.Name
	return agency.Email
}

func getProductBodyData(policy models.Policy, bodyData *BodyData) {
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
		bodyData.ProductForm = "gap/"
	}
}

func getTemplateByChannel(channel, templateType string) []byte {
	var file []byte

	switch channel {
	case models.UserRoleAgency:
		file = lib.GetFilesByEnv(fmt.Sprintf("mail/agent/%s.html", templateType))
	case models.UserRoleAgent:
		file = lib.GetFilesByEnv(fmt.Sprintf("mail/agency/%s.html", templateType))
	}

	return file
}

func fillTemplate(htmlTemplate []byte, bodyData *BodyData, tpl *bytes.Buffer) {
	tmplt := template.New("htmlTemplate")
	tmplt, err := tmplt.Parse(string(htmlTemplate))
	lib.CheckError(err)
	tmplt.Execute(tpl, bodyData)
}
