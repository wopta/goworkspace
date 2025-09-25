package dto

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/models"
)

type AddendumCatnatDTO struct {
	Contract   *addendumContractDTO
	Contractor *contractorDTO
	Signer     *addendumPersonDTO
}

func NewCatnatAddendumDto() *AddendumCatnatDTO {
	return &AddendumCatnatDTO{
		Contract:   &addendumContractDTO{},
		Contractor: &contractorDTO{},
		Signer:     &addendumPersonDTO{},
	}
}

func (b *AddendumCatnatDTO) FromPolicy(policy *models.Policy, now time.Time) {
	b.Contract.fromPolicy(policy, now)
	b.Contractor.fromPolicy(policy.Contractor)
	for _, signer := range *policy.Contractors {
		if signer.IsSignatory {
			b.Signer.fromPolicy(&signer)
			break
		}
	}
}
