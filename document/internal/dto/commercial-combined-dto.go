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
}

func NewCommercialCombinedDto() *CommercialCombinedDTO {
	return &CommercialCombinedDTO{
		ContractDTO:   newContractDTO(),
		ContractorDTO: NewContractorDTO(),
	}
}

func (cc *CommercialCombinedDTO) FromPolicy(policy models.Policy, isProposal bool) {
	cc.ContractDTO.fromPolicy(policy, isProposal)
	cc.ContractorDTO.fromPolicy(policy.Contractor)
}
