package dto

import (
	"github.com/wopta/goworkspace/models"
)

const (
	emptyField = "======"
	no         = "NO"
	yes        = "SI"
)

type CommercialCombinedDTO struct {
	ContractDTO   *contractDTO
	ContractorDTO *contractorDTO
	EnterpriseDTO *enterpriseDTO
	BuildingsDTO  []*buildingDTO
}

func NewCommercialCombinedDto() *CommercialCombinedDTO {
	return &CommercialCombinedDTO{
		ContractDTO:   newContractDTO(),
		ContractorDTO: NewContractorDTO(),
		EnterpriseDTO: newEnterpriseDTO(),
		BuildingsDTO:  make([]*buildingDTO, 0),
	}
}

func (cc *CommercialCombinedDTO) FromPolicy(policy models.Policy, product models.Product, isProposal bool) {
	var numBuildings int64

	cc.ContractDTO.fromPolicy(policy, isProposal)
	cc.ContractorDTO.fromPolicy(policy.Contractor)

	productGuarantees := product.Companies[0].GuaranteesMap

	cc.BuildingsDTO = make([]*buildingDTO, 0, 5)
	for i := 0; i < 5; i++ {
		building := newBuildingDTO()
		for _, guarantee := range productGuarantees {
			if guarantee.Type == "building" {
				newGuarantee := newGuaranteeDTO()
				newGuarantee.Description = guarantee.CompanyName
				building.Guarantees[guarantee.Slug] = newGuarantee
			}
		}
		cc.BuildingsDTO = append(cc.BuildingsDTO, building)
	}

	cc.EnterpriseDTO = newEnterpriseDTO()
	for _, guarantee := range productGuarantees {
		if guarantee.Type == "enterprise" {
			newGuarantee := newGuaranteeDTO()
			newGuarantee.Description = guarantee.CompanyName
			cc.EnterpriseDTO.Guarantees[guarantee.Slug] = newGuarantee
		}
	}

	for index, asset := range policy.Assets {
		if asset.Building != nil {
			cc.BuildingsDTO[numBuildings].fromPolicy(*policy.Assets[index].Building, policy.Assets[index].Guarantees)
			numBuildings++
		}
		if asset.Enterprise != nil {
			cc.EnterpriseDTO.fromPolicy(*policy.Assets[index].Enterprise, policy.Assets[index].Guarantees)
		}
	}
}
