package reserved

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"strings"
)

func LifeReserved(policy models.Policy) models.ReservedInfo {
	log.Println("[LifeReserved]")

	requiredExams := getMedicalDocuments(policy)
	contacts := getContactsDetails(policy)
	documents := getReservedDocument(contacts, requiredExams, policy)

	reservedInfo := models.ReservedInfo{
		RequiredExams: requiredExams,
		Contacts:      contacts,
		Documents:     documents,
	}

	return reservedInfo
}

func getContactsDetails(policy models.Policy) []models.Contact {
	// TODO: check if we can put these info in product file
	return []models.Contact{
		{
			Title:   "Tramite Posta a:",
			Type:    "post",
			Address: "AXA PARTNERS",
			Object:  "Ufficio Underwriting Medico – Corso Como n. 17 – 20154 MILANO",
		},
		{
			Title:   "Tramite e-mail:",
			Type:    "e-mail",
			Address: "clp.it.sinistri@partners.axa",
			Object: fmt.Sprintf("Oggetto: %s proposta %d - UNDERWRITING MEDICO - %s", policy.NameDesc, policy.ProposalNumber,
				strings.ToUpper(policy.Contractor.Surname+" "+policy.Contractor.Name)),
		},
	}
}

func getInputData(policy models.Policy) []byte {
	var err error

	in := make(map[string]interface{})
	in["gender"] = policy.Contractor.Gender
	in["age"], err = policy.CalculateContractorAge()
	lib.CheckError(err)
	maxSumInsuredLimitOfIndemnity := 0.0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Value.SumInsuredLimitOfIndemnity > maxSumInsuredLimitOfIndemnity {
			maxSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
		}
	}
	in["sumInsuredLimitOfIndemnity"] = maxSumInsuredLimitOfIndemnity

	out, err := json.Marshal(in)
	lib.CheckError(err)

	return out
}

func getMedicalDocuments(policy models.Policy) []string {
	const (
		rulesFileName = "life-reserved.json"
	)

	fx := new(models.Fx)
	reservedInfo := &models.ReservedInfo{
		RequiredExams: make([]string, 0),
	}

	rulesFile := lib.GetRulesFile(rulesFileName)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, reservedInfo, getInputData(policy), nil)

	return ruleOutput.(*models.ReservedInfo).RequiredExams
}

func getReservedDocument(contacts []models.Contact, medicalDocuments []string, policy models.Policy) []models.Attachment {
	attachments := make([]models.Attachment, 0)

	gsLink, b := document.LifeReserved(contacts, medicalDocuments, policy)

	attachments = append(attachments, models.Attachment{
		Name:        fmt.Sprintf("%s_proposta_%d_rvm_istruzioni.pdf", policy.NameDesc, policy.ProposalNumber),
		Link:        gsLink,
		Byte:        base64.StdEncoding.EncodeToString(b),
		ContentType: "application/pdf",
	})

	rvmLink := "medical-report/" + policy.Name + "/" + policy.ProductVersion + "/rvm-life.pdf"
	b = lib.GetFromStorage("documents-public-dev", rvmLink, "")

	attachments = append(attachments, models.Attachment{
		Name:        fmt.Sprintf("%s_proposta_%d_rvm.pdf", policy.NameDesc, policy.ProposalNumber),
		Link:        fmt.Sprintf("gs://documents-public-dev/%s", rvmLink),
		Byte:        base64.StdEncoding.EncodeToString(b),
		ContentType: "application/pdf",
	})

	return attachments
}
