package renew

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany", "producerUid"}},

		{"proposalNumber": []string{"proposalNumber", "producerUid"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode", "producerUid"}},
		{"contractorVatCode": []string{"contractorVatCode", "producerUid"}},
		{"contractorFiscalCode": []string{"contractorFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"payment": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"renewMonth": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth", "contractorType"}},
		{"contractorType": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status",
			"payment", "renewMonth", "contractorType"}},
	}

	paramsWhereClause = map[string]string{
		"codeCompany": "(**tableAlias**.codeCompany = @%s)",

		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"insuredFiscalCode":    "(LOWER(JSON_VALUE(**tableAlias**.data, '$.assets[0].person.fiscalCode')) = LOWER(@%s))",
		"contractorVatCode":    "(JSON_VALUE(**tableAlias**.data, '$.contractor.vatCode') = @%s)",
		"contractorFiscalCode": "(LOWER(**tableAlias**.contractorFiscalcode) = LOWER(@%s))",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"company":       "(**tableAlias**.company = LOWER(@%s))",
		"product":       "(**tableAlias**.name = LOWER(@%s))",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",
		"renewMonth":    "(EXTRACT(MONTH FROM **tableAlias**.startDate) = CAST(@%s AS INTEGER))",

		// status
		"paid":   "(**tableAlias**.isPay = true)",
		"unpaid": "(**tableAlias**.isPay = false)",

		// payment
		"recurrent":    "(**tableAlias**.hasMandate = true)",
		"notRecurrent": "(**tableAlias**.hasMandate = false OR **tableAlias**.hasMandate IS NULL)",

		// contractorType
		"enterprise": "(JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'legalEntity' AND (**tableAlias**." +
			"contractorFiscalcode IS NULL OR **tableAlias**.contractorFiscalcode = ''))",
		"individualCompany": "(JSON_VALUE(**tableAlias**.data, " +
			"'$.contractor.type') = 'legalEntity' AND **tableAlias**.contractorFiscalcode != '')",
		"physical": "(JSON_VALUE(p.data, '$.contractor.type') = 'individual' OR (JSON_VALUE(p.data, " +
			"'$.contractor.type') = '') OR (JSON_VALUE(p.data, '$.contractor.type') IS NULL))",
	}

	toBeTranslatedKeys = []string{"status", "payment", "contractorType"}
)
