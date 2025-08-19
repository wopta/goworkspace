package dto

import "gitlab.dev.wopta.it/goworkspace/models"

type CatnatDTO struct {
	SedeDaAssicurare buildingCatnatDto
	Contractor       contractorDTO
	ValidityDate     *validityDateDTO
}

func (dto *CatnatDTO) FromPolicy(policy *models.Policy) {
	dto.SedeDaAssicurare = buildingCatnatDto{}
	dto.SedeDaAssicurare.fromPolicy(policy)
	dto.Contractor = contractorDTO{}
	dto.Contractor.fromPolicy(policy.Contractor)
	dto.ValidityDate = &validityDateDTO{}
	dto.ValidityDate.fromPolicy(policy)
}
