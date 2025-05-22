package pkg

import "gitlab.dev.wopta.it/goworkspace/user/query-builder/internal/user"

func NewQueryBuilder() QueryBuilder {
	return user.NewQueryBuilder()
}
