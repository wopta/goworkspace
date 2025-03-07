package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func SetReservedInfo(policy *models.Policy) {
	switch policy.Name {
	case models.LifeProduct:
		setLifeReservedInfo(policy)
	}
}
