package reserved

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func lifeReserved(policy *models.Policy) {
	log.Println("[LifeReserved]")

	setContactsDetails(policy)
	setMedicalDocuments(policy)
	setReservedDocument(policy)
}

func setContactsDetails(policy *models.Policy) {
	policy.ReservedInfo.Contacts = []models.Contact{
		{
			Title:   "Tramite e-mail:",
			Type:    "e-mail",
			Address: "assunzione@wopta.it",
			Subject: fmt.Sprintf("Oggetto: %s proposta %d - UNDERWRITING MEDICO - %s", policy.NameDesc, policy.ProposalNumber,
				strings.ToUpper(policy.Contractor.Surname+" "+policy.Contractor.Name)),
		},
	}
}

func getInputData(policy *models.Policy) []byte {
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

func setMedicalDocuments(policy *models.Policy) {
	const (
		rulesFileName = "life_reserved.json"
	)

	fx := new(models.Fx)
	reservedInfo := &models.ReservedInfo{
		RequiredExams: make([]string, 0),
	}

	rulesFile := lib.GetRulesFile(rulesFileName)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, reservedInfo, getInputData(policy), nil)

	policy.ReservedInfo.RequiredExams = ruleOutput.(*models.ReservedInfo).RequiredExams
}

func setReservedDocument(policy *models.Policy) {
	attachments := make([]models.Attachment, 0)

	gsLink, b := document.LifeReserved(*policy)

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

	policy.ReservedInfo.Documents = attachments
}
