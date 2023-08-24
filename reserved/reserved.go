package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func GetReservedInfo(policy *models.Policy) {
	switch policy.Name {
	case models.LifeProduct:
		lifeReserved(policy)
	}
}
