package dto

import (
	"strconv"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	emptyField = "======"
	no         = "NO"
	yes        = "SI"
)

type CommercialCombinedDTO struct {
	Contract   *contractDTO
	Contractor *contractorDTO
	Enterprise *enterpriseDTO
	Buildings  []*buildingDTO
	Claims     map[string]*claimDTO
}

func NewCommercialCombinedDto() *CommercialCombinedDTO {
	return &CommercialCombinedDTO{
		Contract:   newContractDTO(),
		Contractor: NewContractorDTO(),
		Enterprise: newEnterpriseDTO(),
		Buildings:  make([]*buildingDTO, 0),
	}
}

func (cc *CommercialCombinedDTO) FromPolicy(policy models.Policy, product models.Product, isProposal bool) {
	var numBuildings int64

	cc.Contract.fromPolicy(policy, isProposal)
	cc.Contractor.fromPolicy(policy.Contractor)

	productGuarantees := product.Companies[0].GuaranteesMap

	cc.Buildings = make([]*buildingDTO, 0, 5)
	for i := 0; i < 5; i++ {
		building := newBuildingDTO()
		for _, guarantee := range productGuarantees {
			if guarantee.Type == "building" {
				newGuarantee := newGuaranteeDTO()
				newGuarantee.Description = guarantee.CompanyName
				building.Guarantees[guarantee.Slug] = newGuarantee
			}
		}
		cc.Buildings = append(cc.Buildings, building)
	}

	cc.Enterprise = newEnterpriseDTO()
	for _, guarantee := range productGuarantees {
		if guarantee.Type == "enterprise" {
			newGuarantee := newGuaranteeDTO()
			newGuarantee.Description = guarantee.CompanyName
			cc.Enterprise.Guarantees[guarantee.Slug] = newGuarantee
		}
	}

	for index, asset := range policy.Assets {
		if asset.Building != nil {
			cc.Buildings[numBuildings].fromPolicy(*policy.Assets[index].Building, policy.Assets[index].Guarantees)
			numBuildings++
		}
		if asset.Enterprise != nil {
			cc.Enterprise.fromPolicy(*policy.Assets[index].Enterprise, policy.Assets[index].Guarantees)
		}
	}

	cc.Claims = make(map[string]*claimDTO)
	// TODO: improve this
	guaranteeDescriptionsMap := map[string]string{
		"property":                "Danni ai beni (escluso Furto)",
		"third-party-liability":   "Responsabilità Civile",
		"theft":                   "Furto",
		"management-organization": "D&O",
		"cyber":                   "Cyber",
	}

	for slug, description := range guaranteeDescriptionsMap {
		cc.Claims[slug] = newClaimDTO()
		cc.Claims[slug].Description = description
	}

	for _, declaredClaim := range policy.DeclaredClaims {
		slug := declaredClaim.GuaranteeSlug
		if _, ok := cc.Claims[slug]; !ok {
			continue
		}
		for _, history := range declaredClaim.History {
			cc.Claims[slug].Quantity.ValueInt += int64(history.Quantity)
			cc.Claims[slug].Quantity.Text = strconv.FormatInt(cc.Claims[slug].Quantity.ValueInt, 10)
			cc.Claims[slug].Value.ValueFloat += history.Value
			cc.Claims[slug].Value.Text = lib.HumanaizePriceEuro(cc.Claims[slug].Value.ValueFloat)
		}
	}
}
