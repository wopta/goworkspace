package dto

import (
	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type validityDateDTO struct {
	StartDate          string
	EndDate            string
	FirstAnnuityExpiry string
}

func newValidityDateDTO() *validityDateDTO {
	return &validityDateDTO{
		StartDate:          constants.EmptyField,
		EndDate:            constants.EmptyField,
		FirstAnnuityExpiry: constants.EmptyField,
	}
}

func (v *validityDateDTO) fromPolicy(p *models.Policy) {
	v.StartDate = formatDate(p.StartDate)
	v.EndDate = formatDate(p.EndDate)
	v.FirstAnnuityExpiry = formatDate(p.StartDate.AddDate(1, 0, 0))
}
