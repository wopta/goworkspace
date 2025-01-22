package policy

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany", "producerUid"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
		{"rd": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "rd"}},
	}

	paramsWhereClause = map[string]string{
		"codeCompany": "(**tableAlias**.codeCompany = @%s)",

		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"insuredFiscalCode": "(JSON_VALUE(**tableAlias**.data, '$.assets[0].person.fiscalCode') = @%s)",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"company":       "(**tableAlias**.company = LOWER(@%s))",
		"product":       "(**tableAlias**.name = LOWER(@%s))",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",
		"reservedYes":   "(**tableAlias**.isReserved = true)",
		"reservedNo":    "(**tableAlias**.isReserved = false)",
		"notSigned":     "(**tableAlias**.isSign = false)",
		"unpaid":        "(**tableAlias**.isPay = false)",
		"signedPaid":    "(**tableAlias**.isSign = true AND **tableAlias**.isPay = true)",
		"unsolved": "(**tableAlias**.annuity > 0 AND **tableAlias**.isPay = false AND CURRENT_DATE(" +
			") >= DATE_ADD(**tableAlias**.startDate, INTERVAL **tableAlias**.annuity YEAR))",
		"renewed": "(**tableAlias**.isRenew = true AND **tableAlias**.isPay = true)",
		"deleted": "(**tableAlias**.isDeleted = true)",
	}

	orClausesKeys = []string{"status"}
)
