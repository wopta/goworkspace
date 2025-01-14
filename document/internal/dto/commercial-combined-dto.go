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

func (cc *CommercialCombinedDTO) FromPolicy(policy models.Policy, isProposal bool) {
	cc.ContractDTO.fromPolicy(policy, isProposal)
	cc.ContractorDTO.fromPolicy(policy.Contractor)

	buildings := make([]*buildingDTO, 0)
	for index, asset := range policy.Assets {
		if asset.Building != nil {
			dto := newBuildingDTO()
			dto.fromPolicy(*policy.Assets[index].Building, policy.Assets[index].Guarantees)
			buildings = append(buildings, dto)
		}
		if asset.Enterprise != nil {
			dto := newEnterpriseDTO()
			dto.fromPolicy(*policy.Assets[index].Enterprise, policy.Assets[index].Guarantees)
			cc.EnterpriseDTO = dto
		}
	}

	numBuildings := len(buildings)
	for i := 0; i < 5-numBuildings; i++ {
		buildings = append(buildings, newBuildingDTO())
	}
	cc.BuildingsDTO = buildings
}
