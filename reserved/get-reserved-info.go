package reserved

import (
	"encoding/json"
	"fmt"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
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

func GetReservedInfoByCoverage(policy *models.Policy) (bool, *models.ReservedInfo, error) {
	var wrapper *PolicyReservedWrapper

	switch policy.Name {
	case models.LifeProduct, models.PersonaProduct:
		wrapper = initWrapper(policy, &ByAssetPerson{}, personAssetExecutor)
	}

	if wrapper == nil {
		return false, nil, nil
	}

	return wrapper.evaluate()
}

func personAssetExecutor(wrapper *PolicyReservedWrapper) (bool, *models.ReservedInfo, error) {
	log.AddPrefix("personAssetExecutor")
	defer log.PopPrefix()

	log.Println("start ------------------------------")

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
		log.ErrorF("error calculating coverage: %s", err.Error())
		return false, nil, err
	}

	if isCovered {
		policies := lib.SliceMap[models.Policy](coveredPolicies, func(p models.Policy) string { return p.CodeCompany })
		reason := fmt.Sprintf("Cliente già assicurato con le polizze numero %v", policies)
		output.IsReserved = isCovered
		output.ReservedInfo.Reasons = append(output.ReservedInfo.Reasons, reason)
	}
	jsonLog, _ := json.Marshal(output)
	log.Printf("result: %v", string(jsonLog))

	log.Println("end --------------------------------")
	return output.IsReserved, output.ReservedInfo, nil
}
