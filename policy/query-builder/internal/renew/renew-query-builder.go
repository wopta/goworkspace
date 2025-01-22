package renew

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/policy/query-builder/internal/base"
)

var (
	paramsHierarchy = []map[string][]string{
		{"codeCompany": []string{"codeCompany", "producerUid"}},

		{"proposalNumber": []string{"proposalNumber", "producerUid"}},

		{"insuredFiscalCode": []string{"insuredFiscalCode", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"company": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"product": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"status": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"payment": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
		{"renewMonth": []string{"startDateFrom", "startDateTo", "company", "product", "producerUid", "status", "payment", "renewMonth"}},
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
		"renewMonth":    "(EXTRACT(MONTH FROM **tableAlias**.startDate) = CAST(@%s AS INTEGER))",
		"paid":          "(**tableAlias**.isPay = true)",
		"unpaid":        "(**tableAlias**.isPay = false)",
		"recurrent":     "(**tableAlias**.hasMandate = true)",
		"notRecurrent":  "(**tableAlias**.hasMandate = false OR **tableAlias**.hasMandate IS NULL)",
	}

	orClausesKeys = []string{"status", "payment"}
)

type QueryBuilder struct {
	base.QueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *QueryBuilder {
	return &QueryBuilder{
		base.NewQueryBuilder(lib.RenewPolicyViewCollection, "rp", randomGenerator,
			paramsHierarchy, paramsWhereClause, orClausesKeys),
	}
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}) {
	qb.WhereClauses = []string{"(**tableAlias**.isDeleted = false OR **tableAlias**." +
		"isDeleted IS NULL)"}

	return qb.QueryBuilder.Build(params)
}
