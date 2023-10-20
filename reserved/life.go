package reserved

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	prd "github.com/wopta/goworkspace/product"
	trn "github.com/wopta/goworkspace/transaction"
)

type ByAssetPerson struct{}

func (*ByAssetPerson) isCovered(w *PolicyReservedWrapper) (bool, *models.Policy, error) {
	var (
		result              = false
		coveredPolicy       *models.Policy
		lastPaidTransaction *models.Transaction
	)
	log.Println("[ByAssetPerson.isCovered] start -----------------------------")

	now := time.Now().UTC()
	lateDate := now.AddDate(0, 2, 0)

	// Check JSON query
	query := fmt.Sprintf(
		"SELECT * FROM `%s.%s` WHERE ",
		models.WoptaDataset,
		models.PolicyCollection,
	)

	policies, err := lib.QueryRowsBigQuery[models.Policy](query)
	if err != nil {
		log.Printf("[ByAssetPerson.isCovered] error getting network transactions: %s", err.Error())
	}

	for _, policy := range policies {
		// check if foundPolicy is in validity range (now ≥ policy.StartDate && now ≤ policy.StartDate)
		// TODO: improve comparison with inclusive dates
		if policy.IsInActiveRange() {
			// TODO: improve query with by date + 2 months
			transactions := trn.GetPolicyTransactions(w.Origin, policy.Uid)
			for _, tr := range transactions {
				if tr.IsPay {
					lastPaidTransaction = &tr
				}
			}
			// check if policy is paid valid (lastPaidTransaction + 2 months < now)
			// check if newPolicy.StartDate ≤ foundPolicy.EndDate;
			if !lastPaidTransaction.IsLate(lateDate) && w.Policy.StartDate.Before(policy.EndDate) {
				result = true
				coveredPolicy = &policy
				break
			}
		}
	}

	fmt.Printf("[ByAssetPerson.isCovered] result '%t'", result)
	fmt.Println("[ByAssetPerson.isCovered] end -------------------------------")

	return result, coveredPolicy, nil
}

func lifeReserved(policy *models.Policy) (bool, *models.ReservedInfo) {
	log.Println("[lifeReserved]")

	var output = ReservedRuleOutput{
		IsReserved: false,
		ReservedInfo: &models.ReservedInfo{
			Reasons:       make([]string, 0),
			RequiredExams: make([]string, 0),
		},
	}

	fx := new(models.Fx)
	rulesFile := lib.GetRulesFileV2(policy.Name, policy.ProductVersion, "reserved")
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
					if question.HasAnswer &&
						question.Answer != nil &&
						question.ExpectedAnswer != nil &&
						*question.Answer != *question.ExpectedAnswer {
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

	_, reservedAge := prd.GetAgeInfo(policy.Name, policy.ProductVersion, policy.Channel)
	data["reservedAge"] = int64(reservedAge)

	ret, err := json.Marshal(data)
	lib.CheckError(err)

	return ret
}

func setLifeReservedDocument(policy *models.Policy, product *models.Product) {
	attachments := make([]models.Attachment, 0)

	gsLink, _ := document.LifeReserved(*policy, product)

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

	policy.ReservedInfo.Documents = attachments
}

func setLifeContactsDetails(policy *models.Policy) {
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

func setLifeReservedInfo(policy *models.Policy, product *models.Product) {
	switch policy.ProductVersion {
	default:
		// TODO: how to handle the contents dinamically?
		setLifeReservedDocument(policy, product)
		setLifeContactsDetails(policy)
	}
}

func lifeReservedByCoverage(wrapper *PolicyReservedWrapper) (bool, *models.ReservedInfo) {
	log.Println("[lifeReservedByCoverage] start ------------------------------")

	var output = ReservedRuleOutput{
		IsReserved: false,
		ReservedInfo: &models.ReservedInfo{
			Reasons: make([]string, 0),
		},
	}

	isCovered, coveredPolicy, err := wrapper.AlreadyCovered.isCovered(wrapper)
	if err != nil {
		// TODO: check handling of error
		panic("help")
	}

	output.IsReserved = isCovered
	if isCovered {
		reason := fmt.Sprintf("Assicurato già coperto dalla polizza %s", coveredPolicy.CodeCompany)
		log.Printf("[lifeReservedByCoverage] %s", reason)
		output.ReservedInfo.Reasons = append(output.ReservedInfo.Reasons, reason)
	}

	log.Println("[lifeReservedByCoverage] end --------------------------------")
	return output.IsReserved, output.ReservedInfo
}
