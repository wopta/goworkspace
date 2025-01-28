package policy

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
		"reservedYes": "(**tableAlias**.isReserved = true)",
		"reservedNo":  "(**tableAlias**.isReserved = false)",

		// status
		"notSigned":  "(**tableAlias**.isSign = false)",
		"unpaid":     "(**tableAlias**.isPay = false) AND (**tableAlias**.isSign = true)",
		"signedPaid": "(**tableAlias**.isSign = true AND **tableAlias**.isPay = true)",
		"unsolved": "(**tableAlias**.annuity > 0 AND **tableAlias**.isPay = false AND CURRENT_DATE(" +
			") >= DATE_ADD(**tableAlias**.startDate, INTERVAL **tableAlias**.annuity YEAR))",
		"renewed": "(**tableAlias**.isRenew = true AND **tableAlias**.isPay = true)",
		"deleted": "(**tableAlias**.isDeleted = true)",

		// contractorType
		"enterprise": "(JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'legalEntity' AND (**tableAlias**." +
			"contractorFiscalcode IS NULL OR **tableAlias**.contractorFiscalcode = ''))",
		"individualCompany": "(JSON_VALUE(**tableAlias**.data, " +
			"'$.contractor.type') = 'legalEntity' AND **tableAlias**.contractorFiscalcode != '')",
		"physical": "(JSON_VALUE(p.data, '$.contractor.type') = 'individual' OR (JSON_VALUE(p.data, " +
			"'$.contractor.type') = '') OR (JSON_VALUE(p.data, '$.contractor.type') IS NULL))",
	}

	toBeTranslatedKeys = []string{"status", "rd", "contractorType"}
)
