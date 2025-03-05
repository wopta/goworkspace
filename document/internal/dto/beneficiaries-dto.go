package dto

import "github.com/wopta/goworkspace/models"

type BeneficiariesDTO struct {
	Contract   *contractDTO
	Contractor *contractorDTO
}

func NewBeneficiariesDto() *BeneficiariesDTO {
	return &BeneficiariesDTO{
		Contract:   newContractDTO(),
		Contractor: newContractorDTO(),
	}
}

func (b *BeneficiariesDTO) FromPolicy(policy models.Policy, product models.Product) {
	b.Contract.fromPolicy(policy, false)
	b.Contractor.fromPolicy(policy.Contractor)

	//for index, asset := range policy.Assets {
	//
	//}
}
