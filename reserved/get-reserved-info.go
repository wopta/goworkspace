package reserved

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetReservedInfo(policy *models.Policy) (bool, *models.ReservedInfo) {
	switch policy.Name {
	case models.LifeProduct:
		return lifeReserved(policy)
	case models.PersonaProduct:
		return personaReserved(policy)
	case models.CommercialCombinedProduct:
		return commercialCombinedReserved(policy)
	default:
		return false, nil
	}
}

func GetReservedInfoByCoverage(policy *models.Policy, origin string) (bool, *models.ReservedInfo, error) {
	var wrapper *PolicyReservedWrapper

	switch policy.Name {
	case models.LifeProduct, models.PersonaProduct:
		wrapper = initWrapper(policy, &ByAssetPerson{}, personAssetExecutor, origin)
	}

	if wrapper == nil {
		return false, nil, nil
	}

	return wrapper.evaluate()
}

func personAssetExecutor(wrapper *PolicyReservedWrapper) (bool, *models.ReservedInfo, error) {
	log.Println("[personAssetExecutor] start ------------------------------")

	var output = ReservedRuleOutput{
		IsReserved:   wrapper.Policy.IsReserved,
		ReservedInfo: wrapper.Policy.ReservedInfo,
	}

	if output.ReservedInfo == nil {
		output.ReservedInfo = &models.ReservedInfo{
			Reasons: make([]string, 0),
		}
	}

	isCovered, coveredPolicies, err := wrapper.AlreadyCovered.isCovered(wrapper)
	if err != nil {
		log.Printf("[personAssetExecutor] error calculating coverage: %s", err.Error())
		return false, nil, err
	}

	if isCovered {
		policies := lib.SliceMap[models.Policy](coveredPolicies, func(p models.Policy) string { return p.CodeCompany })
		reason := fmt.Sprintf("Cliente già assicurato con le polizze numero %v", policies)
		output.IsReserved = isCovered
		output.ReservedInfo.Reasons = append(output.ReservedInfo.Reasons, reason)
	}
	jsonLog, _ := json.Marshal(output)
	log.Printf("[personAssetExecutor] result: %v", string(jsonLog))

	log.Println("[personAssetExecutor] end --------------------------------")
	return output.IsReserved, output.ReservedInfo, nil
}
