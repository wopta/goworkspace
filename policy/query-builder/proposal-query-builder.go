package query_builder

import "github.com/wopta/goworkspace/lib"

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

type proposalQueryBuilder struct {
	baseQueryBuilder
}

func newProposalQueryBuilder(randomGenerator func() string) *proposalQueryBuilder {
	return &proposalQueryBuilder{
		newBaseQueryBuilder(lib.PoliciesViewCollection, "p", randomGenerator,
			proposalParamsHierarchy, proposalParamsWhereClause, proposalOrClausesKeys),
	}
}

func (pqb *proposalQueryBuilder) BuildQuery(params map[string]string) (string, map[string]interface{}) {
	pqb.whereClauses = []string{"(**tableAlias**.proposalNumber > 0)", "(**tableAlias**.companyEmit = false)"}
	return pqb.baseQueryBuilder.BuildQuery(params)
}
