package pkg

import "github.com/wopta/goworkspace/user/query-builder/internal/user"

func NewQueryBuilder(tmp string) QueryBuilder {
	switch tmp {
	case "user":
		return user.NewQueryBuilder(nil)
	}
	return nil
}
