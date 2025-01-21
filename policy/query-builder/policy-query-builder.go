package query_builder

import "github.com/wopta/goworkspace/lib"

var (
	policyParamsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany", "producerUid"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
		{"reserved": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status"}},
	}

	policyParamsWhereClause = map[string]string{
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
		"reservedYes":   "(**table.Alias**.isReserved = true)",
		"reservedNo":    "(**table.Alias**.isReserved = false)",
		"notSigned":     "(**table.Alias**.isSign = false)",
		"unpaid":        "(**table.Alias**.isPay = false)",
		"signedPaid":    "(**table.Alias**.isSign = true AND **table.Alias**.isPay = true)",
		"unsolved": "(**table.Alias**.annuity > 0 AND **table.Alias**.isPay = false AND CURRENT_DATE(" +
			") >= DATE_ADD(**table.Alias**.startDate, INTERVAL **table.Alias**.annuity YEAR))",
		"renewed": "(**table.Alias**.isRenew = true AND **table.Alias**.isPay = true)",
		"deleted": "(**table.Alias**.isDeleted = true)",
	}

	policyOrClausesKeys = []string{"status"}
)

type policyQueryBuilder struct {
	baseQueryBuilder
}

func newPolicyQueryBuilder(randomGenerator func() string) *policyQueryBuilder {
	return &policyQueryBuilder{
		newBaseQueryBuilder(lib.PoliciesViewCollection, "p", randomGenerator,
			policyParamsHierarchy, policyParamsWhereClause, policyOrClausesKeys),
	}
}

func (pqb *policyQueryBuilder) BuildQuery(params map[string]string) (string, map[string]interface{}) {
	pqb.whereClauses = []string{"(**tableAlias**.isDeleted = false OR **tableAlias**." +
		"isDeleted IS NULL)", "(**tableAlias**.companyEmit = true)"}
	return pqb.baseQueryBuilder.BuildQuery(params)

}
