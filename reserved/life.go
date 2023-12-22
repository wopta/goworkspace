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

type ByAssetPerson struct{}

func (*ByAssetPerson) isCovered(w *PolicyReservedWrapper) (bool, []models.Policy, error) {
	var (
		result          = false
		coveredPolicies = make([]models.Policy, 0)
		//lastPaidTransaction *models.Transaction
	)
	log.Println("[ByAssetPerson.isCovered] start -----------------------------")

	query := fmt.Sprintf(
		"SELECT %s FROM `%s.%s` WHERE name = @name AND companyEmit = true AND "+
			"(isDeleted = false OR isDeleted IS NULL) AND "+
			"((@startDate >= startDate AND @startDate <= endDate) OR (@endDate >= startDate AND @endDate <= endDate)) AND "+
			"LOWER(JSON_VALUE(data, '$.assets[0].person.fiscalCode')) = LOWER(@fiscalCode) ORDER BY JSON_VALUE(data, '$.creationDate') ASC",
		"uid, data.creationDate as creationDate, startDate, endDate, codeCompany, paymentSplit, companyEmit, isPay, name, isDeleted",
		models.WoptaDataset,
		models.PoliciesViewCollection,
	)
	params := map[string]interface{}{
		"name":       w.Policy.Name,
		"startDate":  lib.GetBigQueryNullDateTime(w.Policy.StartDate),
		"endDate":    lib.GetBigQueryNullDateTime(w.Policy.EndDate),
		"fiscalCode": w.Policy.Assets[0].Person.FiscalCode,
	}

	log.Printf("[ByAssetPerson.isCovered] executing query %s with params %s", query, params)

	policies, err := lib.QueryParametrizedRowsBigQuery[models.Policy](query, params)
	if err != nil {
		log.Printf("[ByAssetPerson.isCovered] error getting policies: %s", err.Error())
	}
	log.Printf("[ByAssetPerson.isCovered] found %d policies", len(policies))

	if len(policies) > 0 {
		result = true
		coveredPolicies = policies
	}

	/*
		lateDate := time.Now().UTC().AddDate(0, 2, 0)

		for _, policy := range policies {
			log.Printf("[ByAssetPerson.isCovered] checking policy %s", policy.Uid)
			if policy.PaymentSplit == string(models.PaySplitYearly) || policy.PaymentSplit == string(models.PaySplitYear) {
				log.Printf("[ByAssetPerson.isCovered] Yearly pay: found a match! %s - %s", policy.Uid, policy.CodeCompany)
				// As of now, we have only one transaction for annual policies, so no extra control is needed
				// TODO: check behaviour when will have policy renewal
				result = true
				coveredPolicies = append(coveredPolicies, policy)
				continue
			}

			// TODO: remove me when we have a job for deleting policies that have not been paid for X months
			transactions := trn.GetPolicyTransactions(w.Origin, policy.Uid)
			for _, tr := range transactions {
				if tr.IsPay {
					lastPaidTransaction = &tr
				}
			}
			if lastPaidTransaction != nil  && !lastPaidTransaction.IsLate(lateDate) && w.Policy.StartDate.Before(policy.EndDate) {
				log.Printf("[ByAssetPerson.isCovered] Monthly pay: found a match! %s - %s", policy.Uid, policy.CodeCompany)
				result = true
				coveredPolicies = append(coveredPolicies, policy)
			}
		}
	*/

	log.Printf("[ByAssetPerson.isCovered] result '%t'", result)
	log.Println("[ByAssetPerson.isCovered] end -------------------------------")

	return result, coveredPolicies, nil
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

	docInfo := document.Reserved(policy, product)

	attachments = append(attachments, models.Attachment{
		Name: models.RvmInstructionsAttachmentName,
		FileName: strings.ReplaceAll(fmt.Sprintf(models.RvmInstructionsDocumentFormat, policy.NameDesc,
			policy.ProposalNumber), " ", "_"),
		Link:        docInfo.LinkGcs,
		ContentType: "application/pdf",
	})

	rvmLink := "gs://documents-public-dev/medical-report/" + policy.Name + "/" + policy.ProductVersion + "/rvm-life.pdf"

	attachments = append(attachments, models.Attachment{
		Name: models.RvmSurveyAttachmentName,
		FileName: strings.ReplaceAll(fmt.Sprintf(models.RvmSurveyDocumentFormat, policy.NameDesc,
			policy.ProposalNumber), " ", "_"),
		Link:        fmt.Sprintf("%s", rvmLink),
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
		setLifeContactsDetails(policy)
		setLifeReservedDocument(policy, product)
	}
}

func lifeReservedByCoverage(wrapper *PolicyReservedWrapper) (bool, *models.ReservedInfo, error) {
	log.Println("[lifeReservedByCoverage] start ------------------------------")

	var output = ReservedRuleOutput{
		IsReserved:   wrapper.Policy.IsReserved,
		ReservedInfo: wrapper.Policy.ReservedInfo,
	}

	if output.ReservedInfo == nil {
		output.ReservedInfo = &models.ReservedInfo{
			Reasons: make([]string, 0),
		}
	}

	isCovered, coveredPolicies, err := wrapper.AlreadyCovered.isCovered(wrapper)
	if err != nil {
		log.Printf("[lifeReservedByCoverage] error calculating coverage: %s", err.Error())
		return false, nil, err
	}

	if isCovered {
		policies := lib.SliceMap[models.Policy](coveredPolicies, func(p models.Policy) string { return p.CodeCompany })
		reason := fmt.Sprintf("Cliente già assicurato con le polizze numero %v", policies)
		output.IsReserved = isCovered
		output.ReservedInfo.Reasons = append(output.ReservedInfo.Reasons, reason)
	}
	jsonLog, _ := json.Marshal(output)
	log.Printf("[lifeReservedByCoverage] result: %v", string(jsonLog))

	log.Println("[lifeReservedByCoverage] end --------------------------------")
	return output.IsReserved, output.ReservedInfo, nil
}
