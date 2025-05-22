package pkg

import (
	"gitlab.dev.wopta.it/goworkspace/policy/query-builder/internal/policy"
	"gitlab.dev.wopta.it/goworkspace/policy/query-builder/internal/proposal"
	"gitlab.dev.wopta.it/goworkspace/policy/query-builder/internal/renew"
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
