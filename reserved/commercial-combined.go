package reserved

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

// TODO create sub package reserved/commercial-combined
func commercialCombinedReserved(p *models.Policy) (bool, *models.ReservedInfo) {
	log.Println("[commercialCombinedReserved]")

	var output = ReservedRuleOutput{
		IsReserved: false,
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
	return nil
}
