package reserved

import (
	"github.com/wopta/goworkspace/models"
)

type ReservedRuleOutput struct {
	IsReserved   bool
	ReservedInfo *models.ReservedInfo
}

type AlreadyCoveredAlgorithm interface {
	isCovered(*PolicyReservedWrapper) (bool, []models.Policy, error)
}

type PolicyReservedWrapper struct {
	Policy         *models.Policy
	AlreadyCovered AlreadyCoveredAlgorithm
	Evaluate       func(*PolicyReservedWrapper) (bool, *models.ReservedInfo, error)
	Origin         string
}

func (w *PolicyReservedWrapper) evaluate() (bool, *models.ReservedInfo, error) {
	return w.Evaluate(w)
}

func initWrapper(
	p *models.Policy,
	c AlreadyCoveredAlgorithm,
	e func(*PolicyReservedWrapper) (bool, *models.ReservedInfo, error),
	origin string,
) *PolicyReservedWrapper {
	return &PolicyReservedWrapper{
		Policy:         p,
		AlreadyCovered: c,
		Evaluate:       e,
		Origin:         origin,
	}
}
