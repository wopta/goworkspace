package policy

import (
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/policy/query-builder/internal/base"
)

var (
	policyParamsHierarchy = []map[string][]string{
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

type QueryBuilder struct {
	base.QueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *QueryBuilder {
	return &QueryBuilder{
		base.NewQueryBuilder(lib.PoliciesViewCollection, "p", randomGenerator,
			policyParamsHierarchy, policyParamsWhereClause, policyOrClausesKeys),
	}
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}) {
	const (
		deleteClause = "(**tableAlias**.isDeleted = false OR **tableAlias**." +
			"isDeleted IS NULL)"
		emitClause = "(**tableAlias**.companyEmit = true)"
	)
	qb.WhereClauses = []string{emitClause}
	if val, ok := params["status"]; ok {
		if !strings.Contains(val, "deleted") {
			qb.WhereClauses = append(qb.WhereClauses, deleteClause)
		}
	}
	return qb.QueryBuilder.Build(params)
}
