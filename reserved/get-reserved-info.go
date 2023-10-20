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

func GetReservedInfoByCoverage(policy *models.Policy, origin string) (bool, *models.ReservedInfo) {
	var wrapper *PolicyReservedWrapper

	switch policy.Name {
	case models.LifeProduct:
		wrapper = initWrapper(policy, &ByAssetPerson{}, lifeReservedByCoverage, origin)
	}

	if wrapper == nil {
		return false, nil
	}

	return wrapper.evaluate()
}
