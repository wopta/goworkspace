package user

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/user/query-builder/internal/base"
)

type UsersQueryBuilder struct {
	base.UsersQueryBuilder
}

func NewQueryBuilder(randomGenerator func() string) *UsersQueryBuilder {
	return &UsersQueryBuilder{
		base.NewQueryBuilder(lib.UsersViewCollection, "u", randomGenerator,
			paramsHierarchy, paramsWhereClause),
	}
}

func (qb *UsersQueryBuilder) Build(params map[string]string) (string, map[string]interface{}, error) {
	//const (
	//	deleteClause = "(**tableAlias**.isDeleted = false OR **tableAlias**." +
	//		"isDeleted IS NULL)"
	//	emitClause = "(**tableAlias**.companyEmit = true)"
	//)
	//qb.WhereClauses = []string{emitClause, deleteClause}
	//if val, ok := params["status"]; ok {
	//	if strings.Contains(val, "deleted") {
	//		qb.WhereClauses = qb.WhereClauses[:1]
	//	}
	//}

	return qb.UsersQueryBuilder.Build(params)
}
