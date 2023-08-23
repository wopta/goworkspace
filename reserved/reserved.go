package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func GetReservedInfo(policy models.Policy) *models.ReservedInfo {
	var (
		reservedInfo models.ReservedInfo
	)

	switch policy.Name {
	case models.LifeProduct:
		reservedInfo = LifeReserved(policy)
	}

	return &reservedInfo
}
