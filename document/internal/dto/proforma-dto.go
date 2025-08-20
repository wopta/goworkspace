package dto

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type ProformaDTO struct {
	Contractor *proformaContractorDTO
	Body       *proformaBodyDTO
}

type proformaContractorDTO struct {
	NameAndSurname      string
	FiscalOrVatCode     string
	StreetNameAndNumber string
	PostalCodeAndCity   string
	Mail                string
	Phone               string
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
		NameAndSurname:      constants.EmptyField,
		FiscalOrVatCode:     constants.EmptyField,
		StreetNameAndNumber: constants.EmptyField,
		PostalCodeAndCity:   constants.EmptyField,
		Mail:                constants.EmptyField,
		Phone:               constants.EmptyField,
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

func (pf *ProformaDTO) FromPolicy(policy models.Policy) {
	pf.Contractor.fromPolicy(policy.Contractor)
	pf.Body.fromPolicy(policy.ConsultancyValue)
}

func (c *proformaContractorDTO) fromPolicy(contr models.Contractor) {
	if contr.Type == models.UserLegalEntity {
		c.FiscalOrVatCode = contr.VatCode
	} else {
		c.FiscalOrVatCode = contr.FiscalCode
		c.NameAndSurname = lib.TrimSpace(contr.Name + " " + contr.Surname)
	}
	if contr.Residence != nil {
		c.StreetNameAndNumber = contr.Residence.StreetName + " " + contr.Residence.StreetNumber
		c.PostalCodeAndCity = contr.Residence.PostalCode + " " + contr.Residence.Locality + " (" + contr.Residence.CityCode + ")"
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
