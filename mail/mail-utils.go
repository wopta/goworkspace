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
		return "agent"
	}
	if policy.AgencyUid != "" {
		return "agency"
	}

	return "e-commerce"
}

func setBodyDataAndGetCC(policy models.Policy, bodyData *BodyData) string {
	var cc string
	channel := getChannel(policy)

	switch channel {
	case "agent":
		cc = getAgentBodyData(policy.AgentUid, bodyData)
	case "agency":
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
	case "pmi":
		bodyData.ProductName = "Artigiani & Imprese"
		bodyData.ProductForm += "multi-rischio/"
	case "persona":
		bodyData.ProductName = "Persona"
		bodyData.ProductForm += "infortunio/"
	case "life":
		bodyData.ProductName = "Vita"
		bodyData.ProductForm += "vita/"
	case "gap":
		bodyData.ProductName = "Auto Valore Protetto"
		bodyData.ProductForm = "gap/"
	}
}

func getTemplateByChannel(policy models.Policy, templateType string) []byte {

	var file []byte
	channel := getChannel(policy)

	if channel == "agent" {
		file = lib.GetFilesByEnv(fmt.Sprintf("mail/agent/%s.html", templateType))
	}

	if channel == "agency" {
		file = lib.GetFilesByEnv(fmt.Sprintf("mail/agency/%s.html", templateType))
	}

	return file
}

func FillTemplate(htmlTemplate []byte, bodyData *BodyData, tpl *bytes.Buffer) {
	tmplt := template.New("htmlTemplate")
	tmplt, err := tmplt.Parse(string(htmlTemplate))
	lib.CheckError(err)
	tmplt.Execute(tpl, bodyData)
}
