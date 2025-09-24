package policy

import (
	"fmt"

	"gitlab.dev.wopta.it/goworkspace/models"
)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany", "producerUid"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode", "producerUid"}},
		{"contractorVatCode": []string{"contractorVatCode", "producerUid"}},
		{"contractorFiscalCode": []string{"contractorFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"rd": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
		{"contractorType": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd", "contractorType"}},
	}

	paramsWhereClause = map[string]string{
		"codeCompany": "(**tableAlias**.codeCompany = @%s)",

		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"insuredFiscalCode": "(LOWER(JSON_VALUE(**tableAlias**.data, '$.assets[0].person.fiscalCode')) = LOWER(" +
			"@%s))",
		"contractorVatCode":    "(JSON_VALUE(**tableAlias**.data, '$.contractor.vatCode') = @%s)",
		"contractorFiscalCode": "(LOWER(**tableAlias**.contractorFiscalcode) = LOWER(@%s))",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"company":       "(**tableAlias**.company = LOWER(@%s))",
		"product":       "(**tableAlias**.name = LOWER(@%s))",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",

		// rd
		"reservedYes": "(**tableAlias**.isReserved = true AND **tableAlias**.annuity = 0)",
		"reservedNo":  "(**tableAlias**.isReserved = false AND **tableAlias**.annuity = 0)",

		// status
		"notSigned": "(**tableAlias**.isSign = false AND **tableAlias**.annuity = 0)",
		"unpaid": "(**tableAlias**.isPay = false AND **tableAlias**.isSign = true AND **tableAlias**." +
			"annuity = 0)",
		"signedPaid": "(**tableAlias**.isSign = true AND **tableAlias**.isPay = true AND **tableAlias**.annuity = 0)",
		"unsolved": "(**tableAlias**.annuity > 0 AND **tableAlias**.isPay = false AND CURRENT_DATE(" +
			") >= DATE_ADD(**tableAlias**.startDate, INTERVAL **tableAlias**.annuity YEAR))",
		"renewed": "(**tableAlias**.isRenew = true AND **tableAlias**.isPay = true)",
		"deleted": "(**tableAlias**.isDeleted = true)",

		// contractorType
		models.UserLegalEntity: fmt.Sprintf("((JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'legalEntity' AND (JSON_VALUE(**tableAlias**.data,'$.contractor.fiscalCode') IS NULL OR JSON_VALUE(**tableAlias**.data,'$.contractor.fiscalCode') = '')) OR JSON_VALUE(**tableAlias**.data, '$.contractor.type') = '%v')", models.UserLegalEntity),
		models.UserIndividual:  fmt.Sprintf("((JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'legalEntity' AND JSON_VALUE(**tableAlias**.data, '$.contractor.fiscalCode')!= '') OR JSON_VALUE(**tableAlias**.data, '$.contractor.type') = '%v')", models.UserIndividual),
		models.UserPhysical:    fmt.Sprintf("((JSON_VALUE(**tableAlias**.data, '$.contractor.type') = '' OR JSON_VALUE(**tableAlias**.data, '$.contractor.type') IS NULL) OR  JSON_VALUE(**tableAlias**.data, '$.contractor.type') = '%v')", models.UserPhysical),
	}

	toBeTranslatedKeys = []string{"status", "rd", "contractorType"}
)
