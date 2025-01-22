package pkg

import (
	"github.com/wopta/goworkspace/policy/query-builder/internal/policy"
	"github.com/wopta/goworkspace/policy/query-builder/internal/proposal"
	"github.com/wopta/goworkspace/policy/query-builder/internal/renew"
)

func NewQueryBuilder(tmp string) QueryBuilder {
	switch tmp {
	case "policy":
		return policy.NewQueryBuilder(nil)
	case "renew":
		return renew.NewQueryBuilder(nil)
	case "proposal":
		return proposal.NewQueryBuilder(nil)
	}
	return nil
}
