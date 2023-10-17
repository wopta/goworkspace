package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func GetReservedInfo(policy *models.Policy) (bool, *models.ReservedInfo) {
	switch policy.Name {
	case models.LifeProduct:
		return lifeReserved(policy)
	default:
		return false, nil
	}
}

func SetReservedInfo(policy *models.Policy, product *models.Product) {
	switch policy.Name {
	case models.LifeProduct:
		setLifeReservedInfo(policy, product)
	}
}
