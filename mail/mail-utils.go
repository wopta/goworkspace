package mail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func CheckChannel(policy *models.Policy) string {
	agentUid := policy.AgentUid
	agencyUid := policy.AgencyUid

	if agentUid != "" {
		return "agent"
	}
	if agencyUid != "" {
		return "agency"
	}

	return "e-commerce"

}

func SetBodyDataAndGetCC(policy *models.Policy, bodyData *BodyData) string {
	var cc string
	channel := CheckChannel(policy)

	switch channel {
	case "agent":
		cc = GetAgentBodyData(policy.AgentUid, bodyData)
	case "agency":
		cc = GetAgencyBodyData(policy.AgencyUid, bodyData)
	}

	GetProductBodyData(policy, bodyData)

	GetContractorBodyData(policy, bodyData)

	return cc
}

func GetContractorBodyData(policy *models.Policy, bodyData *BodyData) {
	bodyData.ContractorName = policy.Contractor.Name
	bodyData.ContractorSurname = policy.Contractor.Surname
}

func GetAgentBodyData(agentUid string, bodyData *BodyData) string {
	agent, err := models.GetAgentByAuthId(agentUid)
	lib.CheckError(err)
	bodyData.AgentName = agent.Name
	bodyData.AgentSurname = agent.Surname
	return agent.Mail
}

func GetAgencyBodyData(agencyUid string, bodyData *BodyData) string {
	agency, err := models.GetAgencyByAuthId(agencyUid)
	lib.CheckError(err)
	bodyData.AgencyName = agency.Name
	return agency.Email
}

func GetProductBodyData(policy *models.Policy, bodyData *BodyData) {
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
		bodyData.ProductName = "GAP"
		bodyData.ProductForm = "gap/"
	}
}

func GetTemplateByChannel(policy *models.Policy, templateType string) []byte {

	var file []byte
	channel := CheckChannel(policy)

	if channel == "agent" {
		switch os.Getenv("env") {
		case "local":
			file = lib.ErrorByte(ioutil.ReadFile(fmt.Sprintf("../function-data/dev/mail/agent/%s.html", templateType)))
			// case "dev":
			// 	file = lib.GetFromStorage("function-data", "mail/mail_template.html", "")
			// case "prod":
			// 	file = lib.GetFromStorage("core-350507-function-data", "mail/mail_template.html", "")
		}
	}

	if channel == "agency" {
		switch os.Getenv("env") {
		case "local":
			file = lib.ErrorByte(ioutil.ReadFile(fmt.Sprintf("../function-data/dev/mail/agency/%s.html", templateType)))
			// case "dev":
			// 	file = lib.GetFromStorage("function-data", "mail/mail_template.html", "")
			// case "prod":
			// 	file = lib.GetFromStorage("core-350507-function-data", "mail/mail_template.html", "")
		}
	}

	return file

}

func FillTemplate(htmlTemplate []byte, bodyData *BodyData, tpl *bytes.Buffer) {

	tmplt := template.New("htmlTemplate")
	tmplt, err := tmplt.Parse(string(htmlTemplate))
	lib.CheckError(err)
	tmplt.Execute(tpl, bodyData)
}
