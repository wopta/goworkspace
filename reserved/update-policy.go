package reserved

import (
	"gitlab.dev.wopta.it/goworkspace/models"
)

func UpdatePolicyReserved(policy *models.Policy) map[string]interface{} {
	input := make(map[string]interface{}, 0)

	isReserved, reservedInfo := GetReservedInfo(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input
}

func UpdatePolicyReservedCoverage(policy *models.Policy) (map[string]interface{}, error) {
	input := make(map[string]interface{}, 0)

	isReserved, reservedInfo, err := GetReservedInfoByCoverage(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input, err
}
