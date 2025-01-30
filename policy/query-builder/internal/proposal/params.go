package proposal

var (
	paramsHierarchy = []map[string][]string{
		{"proposalNumber": []string{"proposalNumber", "producerUid"}},

		{"contractorVatCode": []string{"contractorVatCode", "producerUid"}},
		{"contractorFiscalCode": []string{"contractorFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
		{"rd": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
		{"buyable": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
		{"contractorType": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable", "contractorType"}},
	}

	paramsWhereClause = map[string]string{
		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"contractorVatCode":    "(JSON_VALUE(**tableAlias**.data, '$.contractor.vatCode') = @%s)",
		"contractorFiscalCode": "(LOWER(**tableAlias**.contractorFiscalcode) = LOWER(@%s))",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",

		// rd
		"toBeStarted": "(**tableAlias**.status = 'NeedsApproval')",
		"inProgress":  "(**tableAlias**.status = 'WaitForApproval')",
		"denied":      "(**tableAlias**.status = 'Rejected')",

		// buyable
		"proposal": "(**tableAlias**.status = 'Proposal')",
		"approved": "(**tableAlias**.status = 'Approved')",

		// contractorType
		"enterprise": "(JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'legalEntity' AND (**tableAlias**." +
			"contractorFiscalcode IS NULL OR **tableAlias**.contractorFiscalcode = ''))",
		"individualCompany": "(JSON_VALUE(**tableAlias**.data, " +
			"'$.contractor.type') = 'legalEntity' AND **tableAlias**.contractorFiscalcode != '')",
		"physical": "(JSON_VALUE(**tableAlias**.data, '$.contractor.type') = 'individual' OR (JSON_VALUE(**tableAlias**.data, " +
			"'$.contractor.type') = '') OR (JSON_VALUE(**tableAlias**.data, '$.contractor.type') IS NULL))",
	}

	toBeTranslatedKeys = []string{"rd", "buyable", "contractorType"}
)
