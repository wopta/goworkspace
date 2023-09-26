package reserved

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
)

type ReservedRuleOutput struct {
	IsReserved   bool
	ReservedInfo *models.ReservedInfo
}

func lifeReserved(policy *models.Policy) (bool, *models.ReservedInfo) {
	log.Println("[lifeReserved]")

	const (
		rulesFileName = "life_reserved.json"
	)

	var output = ReservedRuleOutput{
		IsReserved: false,
		ReservedInfo: &models.ReservedInfo{
			Reasons:       make([]string, 0),
			RequiredExams: make([]string, 0),
		},
	}

	fx := new(models.Fx)
	rulesFile := lib.GetRulesFile(rulesFileName)
	input := getInputData(policy)
	log.Printf("[lifeReserved] input %v", string(input))
	data := getReservedData(policy)
	log.Printf("[lifeReserved] data %v", string(data))

	ruleOutputString, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, &output, input, data)

	log.Printf("[lifeReserved] rules output: %s", ruleOutputString)

	return ruleOutput.(*ReservedRuleOutput).IsReserved, ruleOutput.(*ReservedRuleOutput).ReservedInfo
}

func getInputData(policy *models.Policy) []byte {
	var (
		err error
		in  = make(map[string]interface{})
	)

	in["gender"] = policy.Contractor.Gender

	age, err := policy.CalculateContractorAge()
	lib.CheckError(err)
	in["age"] = int64(age)

	maxSumInsuredLimitOfIndemnity := 0.0
	in["death"] = 0.0
	in["permanent-disability"] = 0.0
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Slug == "death" {
			in["death"] = guarantee.Value.SumInsuredLimitOfIndemnity
		}
		if guarantee.Slug == "permanent-disability" {
			in["permanent-disability"] = guarantee.Value.SumInsuredLimitOfIndemnity
		}
		if guarantee.Value.SumInsuredLimitOfIndemnity > maxSumInsuredLimitOfIndemnity {
			maxSumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
		}
	}
	in["sumInsuredLimitOfIndemnity"] = maxSumInsuredLimitOfIndemnity

	in["surveys"] = false
	if policy.Surveys != nil {
		for _, survey := range *policy.Surveys {
			if survey.HasAnswer && (*survey.Answer != *survey.ExpectedAnswer) {
				in["surveys"] = true
				break
			} else if survey.HasMultipleAnswers != nil && *survey.HasMultipleAnswers {
				for _, question := range survey.Questions {
					if question.HasAnswer && question.Answer != question.ExpectedAnswer {
						in["surveys"] = true
						break
					}
				}
			}
		}
	}

	out, err := json.Marshal(in)
	lib.CheckError(err)

	return out
}

func getReservedData(policy *models.Policy) []byte {
	data := make(map[string]interface{})

	reservedAge := prd.GetReservedAge(policy.Name, models.GetChannel(policy))
	data["reservedAge"] = int64(reservedAge)

	ret, err := json.Marshal(data)
	lib.CheckError(err)

	return ret
}

func GetLifeReservedDocument(policy *models.Policy) []models.Attachment {
	attachments := make([]models.Attachment, 0)

	gsLink, _ := document.LifeReserved(*policy)

	attachments = append(attachments, models.Attachment{
		Name:        fmt.Sprintf("%s_proposta_%d_rvm_istruzioni.pdf", policy.NameDesc, policy.ProposalNumber),
		Link:        gsLink,
		ContentType: "application/pdf",
	})

	rvmLink := "medical-report/" + policy.Name + "/" + policy.ProductVersion + "/rvm-life.pdf"

	attachments = append(attachments, models.Attachment{
		Name:        fmt.Sprintf("%s_proposta_%d_rvm.pdf", policy.NameDesc, policy.ProposalNumber),
		Link:        fmt.Sprintf("gs://documents-public-dev/%s", rvmLink),
		ContentType: "application/pdf",
	})

	return attachments
}

func GetLifeContactsDetails(policy *models.Policy) []models.Contact {
	return []models.Contact{
		{
			Title:   "Tramite e-mail:",
			Type:    "e-mail",
			Address: "assunzione@wopta.it",
			Subject: fmt.Sprintf("Oggetto: %s proposta %d - UNDERWRITING MEDICO - %s", policy.NameDesc, policy.ProposalNumber,
				strings.ToUpper(policy.Contractor.Surname+" "+policy.Contractor.Name)),
		},
	}
}
