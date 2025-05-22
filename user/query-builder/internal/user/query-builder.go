package user

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/user/query-builder/internal/base"
)

type UsersQueryBuilder struct {
	base.UsersQueryBuilder
}

func NewQueryBuilder() *UsersQueryBuilder {
	return &UsersQueryBuilder{
		base.NewQueryBuilder(lib.UsersViewCollection, "u",
			paramsHierarchy, paramsWhereClause),
	}
}

func (qb *UsersQueryBuilder) Build(params map[string]string) (string, map[string]interface{}, error) {
	return qb.UsersQueryBuilder.Build(params)
}
