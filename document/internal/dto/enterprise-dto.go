package dto

import (
	"github.com/wopta/goworkspace/models"
)

type enterpriseDTO struct {
	Revenue                   float64
	NorthAmericanMarket       float64
	Employer                  int64
	WorkEmployersRemuneration float64
	TotalBilled               float64
	OwnerTotalBilled          float64
	Guarantees                map[string]*guaranteeDTO
}

func newEnterpriseDTO() *enterpriseDTO {
	return &enterpriseDTO{
		Revenue:                   0,
		NorthAmericanMarket:       0,
		Employer:                  0,
		WorkEmployersRemuneration: 0,
		Guarantees:                make(map[string]*guaranteeDTO),
	}
}

func (e *enterpriseDTO) fromPolicy(enterprise models.Enterprise, guarantees []models.Guarante) {
	if enterprise.Revenue != 0.0 {
		e.Revenue = enterprise.Revenue
	}
	if enterprise.NorthAmericanMarket != 0.0 {
		e.NorthAmericanMarket = enterprise.NorthAmericanMarket
	}
	if enterprise.Employer != 0.0 {
		e.Employer = enterprise.Employer
	}
	if enterprise.WorkEmployersRemuneration != 0.0 {
		e.WorkEmployersRemuneration = enterprise.WorkEmployersRemuneration
	}

	for _, guarantee := range guarantees {
		if val, ok := e.Guarantees[guarantee.Slug]; ok {
			val.fromPolicy(guarantee)
		}

	}
}
