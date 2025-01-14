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
	ContractDTO *contractDTO
}

func NewCommercialCombinedDto() *CommercialCombinedDTO {
	return &CommercialCombinedDTO{
		ContractDTO: newContractDTO(),
	}
}

func (cc *CommercialCombinedDTO) ParseFromPolicy(policy models.Policy, isProposal bool) {
	cc.ContractDTO.parseFromPolicy(policy, isProposal)
}
