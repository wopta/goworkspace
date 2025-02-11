package reserved

import (
	"encoding/json"
	"log"
	"slices"

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

	// enterprise data
	in["revenue"] = float64(0)
	in["northAmericanMarket"] = float64(0)
	// enterprise guarantee
	in["third-party-recourse"] = float64(0)
	in["electrical-phenomenon"] = float64(0)
	in["refrigeration-stock"] = float64(0)
	in["machinery-breakdown"] = float64(0)
	in["electronic-equipment"] = float64(0)
	in["theft"] = float64(0)
	in["daily-allowance"] = float64(0)
	in["increased-cost"] = float64(0)
	in["additional-compensation"] = float64(0)
	in["loss-rent"] = float64(0)
	in["third-party-liability-work-providers"] = float64(0)
	in["product-liability"] = float64(0)
	in["product-withdrawal"] = float64(0)
	in["management-organization"] = float64(0)
	in["cyber"] = float64(0)
	in["management-organization-asset"] = float64(0)
	in["management-organization-capital"] = float64(0)
	// building guarantee
	in["building"] = float64(0)
	in["rental-risk"] = float64(0)
	in["machinery"] = float64(0)
	in["stock"] = float64(0)
	in["stock-temporary-increase"] = float64(0)
	in["sum"] = float64(0)
	// policy
	in["hasDiscount"] = false
	in["hasBond"] = p.HasBond
	in["hasClause"] = p.Clause != ""
	in["naics"] = false

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
		in["northAmericanMarket"] = enterpriseAsset.Enterprise.NorthAmericanMarket
		for _, g := range enterpriseAsset.Guarantees {
			in[g.Slug] = g.Value.SumInsuredLimitOfIndemnity
			if g.Slug == "management-organization" {
				in["management-organization-asset"] = g.Value.LimitOfIndemnity
				in["management-organization-capital"] = g.Value.SumInsured
			}
			if g.Value.Discount != 0 {
				in["hasDiscount"] = true
			}
		}
	}

	var sum float64
	for _, b := range buildingsAssets {
		for _, g := range b.Guarantees {
			if slices.Contains([]string{"building", "machinery", "stock", "third-party-recourse",
				"stock-temporary-increase", "daily-allowance", "additional-compensation", "increased-cost"}, g.Slug) {
				sum += g.Value.SumInsuredLimitOfIndemnity
			}
			switch v := in[g.Slug].(type) {
			case float64:
				v += g.Value.SumInsuredLimitOfIndemnity
				in[g.Slug] = v
			default:
				log.Printf("%T", v)
			}
			if g.Value.Discount != 0 {
				in["hasDiscount"] = true
			}
		}
		if !b.Building.IsNaicsSellable {
			in["naics"] = true
		}
	}
	in["sum"] = sum

	out, err := json.Marshal(in)
	if err != nil {
		return nil
	}

	return out
}
