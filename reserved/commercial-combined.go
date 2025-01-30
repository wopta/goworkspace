package reserved

import (
	"encoding/json"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

// TODO create sub package reserved/commercial-combined
func commercialCombinedReserved(p *models.Policy) (bool, *models.ReservedInfo) {
	log.Println("[commercialCombinedReserved]")

	var output = ReservedRuleOutput{
		IsReserved:   false,
		ReservedInfo: &models.ReservedInfo{},
	}

	output.ReservedInfo.CompanyApproval.Mandatory = true

	fx := new(models.Fx)
	rulesFile := lib.GetRulesFileV2(p.Name, p.ProductVersion, "reserved")

	input := getCCInputData(p)
	log.Printf("input data: %+v", string(input))

	ruleOutputString, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, &output, input, nil)

	log.Printf("[lifeReserved] rules output: %s", ruleOutputString)

	return ruleOutput.(*ReservedRuleOutput).IsReserved, ruleOutput.(*ReservedRuleOutput).ReservedInfo
}

func getCCInputData(p *models.Policy) []byte {
	var (
		enterpriseAsset models.Asset
		buildingsAssets = make([]models.Asset, 0)
		in              = make(map[string]interface{})
	)

	in["revenue"] = float64(0)

	for _, a := range p.Assets {
		switch a.Type {
		case models.AssetTypeEnterprise:
			enterpriseAsset = a
		case models.AssetTypeBuilding:
			buildingsAssets = append(buildingsAssets, a)
		}
	}

	if enterpriseAsset.Enterprise != nil {
		in["revenue"] = enterpriseAsset.Enterprise.Revenue
	}

	out, err := json.Marshal(in)
	if err != nil {
		return nil
	}

	return out
}
