package renew

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/policy/query-builder/internal/base"
)

type QueryBuilder struct {
	base.QueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *QueryBuilder {
	return &QueryBuilder{
		base.NewQueryBuilder(lib.RenewPolicyViewCollection, "rp", randomGenerator,
			paramsHierarchy, paramsWhereClause, toBeTranslatedKeys),
	}
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}, error) {
	qb.WhereClauses = []string{"(**tableAlias**.isDeleted = false OR **tableAlias**." +
		"isDeleted IS NULL)"}

	return qb.QueryBuilder.Build(params)
}
