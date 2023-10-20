package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func UpdatePolicy(policy *models.Policy) map[string]interface{} {
	input := make(map[string]interface{}, 0)

	isReserved, reservedInfo := GetReservedInfo(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input
}
