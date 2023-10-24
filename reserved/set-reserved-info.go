package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func SetReservedInfo(policy *models.Policy, product *models.Product) {
	switch policy.Name {
	case models.LifeProduct:
		setLifeReservedInfo(policy, product)
	}
}
