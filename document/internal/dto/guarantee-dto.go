package dto

import (
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/models"
)

type GuaranteeDTO struct {
	Description                string
	SumInsuredLimitOfIndemnity float64
	LimitOfIndemnity           float64
	SumInsured                 float64
	StartDate                  string
	Duration                   int64
}

func newGuaranteeDTO() *GuaranteeDTO {
	return &GuaranteeDTO{
		Description:                emptyField,
		SumInsuredLimitOfIndemnity: 0,
		LimitOfIndemnity:           0,
		SumInsured:                 0,
		StartDate:                  emptyField,
		Duration:                   0,
	}
}

func (g *GuaranteeDTO) fromPolicy(guarantee models.Guarante) {
	if guarantee.CompanyName != "" {
		g.Description = guarantee.Description
	}
	g.SumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
	g.LimitOfIndemnity = guarantee.Value.LimitOfIndemnity
	g.SumInsured = guarantee.Value.SumInsured
	if guarantee.Value.StartDate != nil && !guarantee.Value.StartDate.IsZero() {
		g.StartDate = guarantee.Value.StartDate.Format(constants.DayMonthYearFormat)
	}
	if guarantee.Value.Duration != nil {
		g.Duration = int64(guarantee.Value.Duration.Year)
	}
}
