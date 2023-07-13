package document

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"strings"
)

func loadLifeBeneficiariesInfo(policy *models.Policy) ([]map[string]string, string, string) {
	legitimateSuccessorsChoice := "X"
	designatedSuccessorsChoice := ""
	beneficiaries := []map[string]string{
		{
			"name":       "=====",
			"fiscalCode": "=====",
			"address":    "=====",
			"mail":       "=====",
			"phone":      "=====",
			"relation":   "=====",
			"consent":    "=====",
		},
		{
			"name":           "=====",
			"fiscalCode":     "=====",
			"address":        "=====",
			"mail":           "=====",
			"phone":          "=====",
			"relation":       "=====",
			"contactConsent": "=====",
		},
	}

	deathGuarantee, err := policy.ExtractGuarantee("death")
	lib.CheckError(err)

	if deathGuarantee.Beneficiaries != nil && !(*deathGuarantee.Beneficiaries)[0].IsLegitimateSuccessors {
		legitimateSuccessorsChoice = ""
		designatedSuccessorsChoice = "X"

		for index, beneficiary := range *deathGuarantee.Beneficiaries {
			address := strings.ToUpper(beneficiary.Residence.StreetName + ", " + beneficiary.Residence.StreetNumber +
				" - " + beneficiary.Residence.PostalCode + " " + beneficiary.Residence.City +
				" (" + beneficiary.Residence.CityCode + ")")
			beneficiaries[index]["name"] = strings.ToUpper(beneficiary.Surname + " " + beneficiary.Name)
			beneficiaries[index]["fiscalCode"] = strings.ToUpper(beneficiary.FiscalCode)
			beneficiaries[index]["address"] = address
			beneficiaries[index]["mail"] = beneficiary.Mail
			beneficiaries[index]["phone"] = beneficiary.Phone
			if beneficiary.IsFamilyMember {
				beneficiaries[index]["relation"] = "Nucleo familiare (rapporto di parentela, coniuge, unione civile, " +
					"convivenza more uxorio)"
			} else {
				beneficiaries[index]["relation"] = "Altro (no rapporto parentela)"
			}
			if beneficiary.IsContactable {
				beneficiaries[index]["contactConsent"] = "SI"
			} else {
				beneficiaries[index]["contactConsent"] = "NO"
			}
		}
	}

	return beneficiaries, legitimateSuccessorsChoice, designatedSuccessorsChoice
}

func loadProducerInfo(origin string, policy *models.Policy) map[string]string {
	policyProducer := map[string]string{
		"name":            "LOMAZZI MICHELE",
		"ruiSection":      "A",
		"ruiCode":         "A000703480",
		"ruiRegistration": "02.03.2022",
	}

	if policy.AgentUid != "" {
		var agent models.Agent
		fireAgent := lib.GetDatasetByEnv(origin, models.AgentCollection)
		docsnap, err := lib.GetFirestoreErr(fireAgent, policy.AgentUid)
		lib.CheckError(err)
		err = docsnap.DataTo(&agent)
		lib.CheckError(err)
		policyProducer["name"] = strings.ToUpper(agent.Surname) + " " + strings.ToUpper(agent.Name)
		policyProducer["ruiSection"] = agent.RuiSection
		policyProducer["ruiCode"] = agent.RuiCode
		policyProducer["ruiRegistration"] = agent.RuiRegistration.Format("02.01.2006")
	} else if policy.AgencyUid != "" {
		var agency models.Agency
		fireAgency := lib.GetDatasetByEnv(origin, models.AgencyCollection)
		docsnap, err := lib.GetFirestoreErr(fireAgency, policy.AgencyUid)
		lib.CheckError(err)
		err = docsnap.DataTo(&agency)
		lib.CheckError(err)
		policyProducer["name"] = strings.ToUpper(agency.Name)
		policyProducer["ruiSection"] = agency.RuiSection
		policyProducer["ruiCode"] = agency.RuiCode
		policyProducer["ruiRegistration"] = agency.RuiRegistration.Format("02.01.2006")
	}
	return policyProducer
}
