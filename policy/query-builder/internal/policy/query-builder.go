package policy

import (
	"strings"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/policy/query-builder/internal/base"
)

type QueryBuilder struct {
	base.QueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *QueryBuilder {
	return &QueryBuilder{
		base.NewQueryBuilder(lib.PoliciesViewCollection, "p", randomGenerator,
			paramsHierarchy, paramsWhereClause, toBeTranslatedKeys),
	}
}

func (qb *QueryBuilder) Build(params map[string]string) (string, map[string]interface{}, error) {
	const (
		deleteClause = "(**tableAlias**.isDeleted = false OR **tableAlias**." +
			"isDeleted IS NULL)"
		emitClause = "(**tableAlias**.companyEmit = true)"
	)
	qb.WhereClauses = []string{emitClause, deleteClause}
	if val, ok := params["status"]; ok {
		if strings.Contains(val, "deleted") {
			qb.WhereClauses = qb.WhereClauses[:1]
		}
	}
	return qb.QueryBuilder.Build(params)
}
