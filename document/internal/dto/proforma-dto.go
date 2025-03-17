package dto

import (
	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/models"
)

type ProformaDTO struct {
	Contractor *ContractorDTO
	Body       *BodyDTO
}

type ContractorDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Province     string
	PostalCode   string
	Mail         string
	Phone        string
}

type BodyDTO struct {
	Date    string
	Net     string
	Vat     string
	Gross   string
	PayDate string
}

func NewProformaDTO() *ProformaDTO {
	return &ProformaDTO{
		Contractor: NewContractorDTO(),
		Body:       NewBodyDTO(),
	}
}

func NewContractorDTO() *ContractorDTO {
	return &ContractorDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Province:     constants.EmptyField,
		PostalCode:   constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
	}
}

func NewBodyDTO() *BodyDTO {
	return &BodyDTO{
		Date:    constants.EmptyField,
		Net:     constants.EmptyField,
		Vat:     constants.EmptyField,
		Gross:   constants.EmptyField,
		PayDate: constants.EmptyField,
	}
}

func (cc *ProformaDTO) FromPolicy(policy models.Policy, product models.Product) {

}
