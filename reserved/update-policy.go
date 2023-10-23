package reserved

import (
	"github.com/wopta/goworkspace/models"
)

func UpdatePolicyReserved(policy *models.Policy) map[string]interface{} {
	input := make(map[string]interface{}, 0)

	isReserved, reservedInfo := GetReservedInfo(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input
}

func UpdatePolicyReservedCoverage(policy *models.Policy, origin string) (map[string]interface{}, error) {
	input := make(map[string]interface{}, 0)

	isReserved, reservedInfo, err := GetReservedInfoByCoverage(policy, origin)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input, err
}
