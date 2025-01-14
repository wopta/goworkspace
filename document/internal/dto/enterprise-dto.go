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
	Guarantees                map[string]*GuaranteeDTO
}

func newEnterpriseDTO() *enterpriseDTO {
	return &enterpriseDTO{
		Revenue:                   0,
		NorthAmericanMarket:       0,
		Employer:                  0,
		WorkEmployersRemuneration: 0,
		Guarantees:                make(map[string]*GuaranteeDTO),
	}
}

func (e *enterpriseDTO) fromPolicy(assets []models.Asset) {
	var enterprise models.Enterprise
	var guarantees []models.Guarante

	for index, asset := range assets {
		if asset.Enterprise != nil {
			enterprise = *assets[index].Enterprise
			guarantees = assets[index].Guarantees
			break
		}
	}

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
		guaranteeDTO := newGuaranteeDTO()
		if len(guarantee.CompanyName) != 0 {
			guaranteeDTO.Description = guarantee.CompanyName
		}
		guaranteeDTO.SumInsuredLimitOfIndemnity = guarantee.Value.SumInsuredLimitOfIndemnity
		guaranteeDTO.LimitOfIndemnity = guarantee.Value.LimitOfIndemnity
		guaranteeDTO.SumInsured = guarantee.Value.SumInsured

		e.Guarantees[guarantee.Slug] = guaranteeDTO
	}
}
