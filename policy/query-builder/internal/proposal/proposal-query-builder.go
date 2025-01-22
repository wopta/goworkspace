package proposal

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/policy/query-builder/internal/base"
)

var (
	proposalParamsHierarchy = []map[string][]string{
		{"proposalNumber": []string{"proposalNumber", "producerUid"}},

		{"contractorName": []string{"contractorName", "contractorSurname", "producerUid"}},
		{"contractorSurname": []string{"contractorName", "contractorSurname", "producerUid"}},

		{"startDateFrom": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable"}},
		{"startDateTo": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable"}},
		{"producerUid": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable"}},
		{"rd": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable"}},
		{"buyable": []string{"startDateFrom", "startDateTo", "producerUid", "rd", "buyable"}},
	}

	proposalParamsWhereClause = map[string]string{
		"proposalNumber": "(**tableAlias**.proposalNumber = CAST(@%s AS INTEGER))",

		"contractorName":    "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.name')), LOWER(@%s)))",
		"contractorSurname": "(REGEXP_CONTAINS(LOWER(JSON_VALUE(**tableAlias**.data, '$.contractor.surname')), LOWER(@%s)))",

		"startDateFrom": "(**tableAlias**.startDate >= @%s)",
		"startDateTo":   "(**tableAlias**.startDate <= @%s)",
		"producerUid":   "(**tableAlias**.producerUid IN (%s))",
		"toBeStarted":   "(**tableAlias**.status = 'NeedsApproval')",
		"inProgress":    "(**tableAlias**.status = 'WaitForApproval')",
		"denied":        "(**tableAlias**.status = 'Rejected')",
		"proposal":      "(**tableAlias**.status = 'Proposal')",
		"approved":      "(**tableAlias**.status = 'Approved')",
	}

	proposalOrClausesKeys = []string{"rd", "buyable"}
)

type QueryBuilder struct {
	base.QueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *QueryBuilder {
	return &QueryBuilder{
		base.NewQueryBuilder(lib.PoliciesViewCollection, "p", randomGenerator,
			proposalParamsHierarchy, proposalParamsWhereClause, proposalOrClausesKeys),
	}
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}) {
	qb.WhereClauses = []string{"(**tableAlias**.proposalNumber > 0)", "(**tableAlias**.companyEmit = false)"}
	return qb.QueryBuilder.Build(params)
}
