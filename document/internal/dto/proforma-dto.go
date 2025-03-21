package dto

import (
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type ProformaDTO struct {
	Contractor *proformaContractorDTO
	Body       *proformaBodyDTO
}

type proformaContractorDTO struct {
	Name            string
	Surname         string
	FiscalOrVatCode string
	StreetName      string
	StreetNumber    string
	City            string
	Province        string
	PostalCode      string
	Mail            string
	Phone           string
}

type proformaBodyDTO struct {
	Date    string
	Net     string
	Vat     string
	Gross   string
	PayDate string
}

func NewProformaDTO() *ProformaDTO {
	return &ProformaDTO{
		Contractor: newProformaContractorDTO(),
		Body:       newProformaBodyDTO(),
	}
}

func newProformaContractorDTO() *proformaContractorDTO {
	return &proformaContractorDTO{
		Name:            constants.EmptyField,
		Surname:         constants.EmptyField,
		FiscalOrVatCode: constants.EmptyField,
		StreetName:      constants.EmptyField,
		StreetNumber:    constants.EmptyField,
		City:            constants.EmptyField,
		Province:        constants.EmptyField,
		PostalCode:      constants.EmptyField,
		Mail:            constants.EmptyField,
		Phone:           constants.EmptyField,
	}
}

func newProformaBodyDTO() *proformaBodyDTO {
	return &proformaBodyDTO{
		Date:    constants.EmptyField,
		Net:     constants.EmptyField,
		Vat:     constants.EmptyField,
		Gross:   constants.EmptyField,
		PayDate: constants.EmptyField,
	}
}

func (pf *ProformaDTO) FromPolicy(policy models.Policy, product models.Product) {
	pf.Contractor.fromPolicy(policy.Contractor)
	pf.Body.fromPolicy(policy.ConsultancyValue)
}

func (c *proformaContractorDTO) fromPolicy(contr models.Contractor) {
	c.Name = contr.Name
	c.Surname = contr.Surname
	if contr.VatCode != "" {
		c.FiscalOrVatCode = contr.VatCode
	} else {
		c.FiscalOrVatCode = contr.FiscalCode
	}
	if contr.Residence != nil {
		c.StreetName = contr.Residence.StreetName
		c.StreetNumber = contr.Residence.StreetNumber
		c.City = contr.Residence.Locality
		c.Province = contr.Residence.CityCode
		c.PostalCode = contr.Residence.PostalCode
	}
	c.Mail = contr.Mail
	c.Phone = contr.Phone
}

func (b *proformaBodyDTO) fromPolicy(value models.ConsultancyValue) {
	b.Gross = lib.HumanaizePriceEuro(value.Price)
	b.Net = b.Gross
	b.Vat = lib.HumanaizePriceEuro(0.00)
	b.Date = time.Now().Format(constants.DayMonthYearFormat)
}
