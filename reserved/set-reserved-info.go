package reserved

import (
	"gitlab.dev.wopta.it/goworkspace/models"
)

func SetReservedInfo(policy *models.Policy) {
	switch policy.Name {
	case models.LifeProduct:
		setLifeReservedInfo(policy)
	}
}
