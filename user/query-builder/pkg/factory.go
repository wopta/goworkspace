package pkg

import "github.com/wopta/goworkspace/user/query-builder/internal/user"

func NewQueryBuilder() QueryBuilder {
	return user.NewQueryBuilder(nil)
}
