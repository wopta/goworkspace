package dto

import (
	"strings"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type contractorDTO struct {
	Name         string
	Surname      string
	CompanyName  string
	FiscalCode   string
	VatCode      string
	StreetName   string
	StreetNumber string
	City         string
	PostalCode   string
	CityCode     string
	Mail         string
	Phone        string
	BirthDate    string
	Address      string
}

func newContractorDTO() *contractorDTO {
	return &contractorDTO{
		Name:         constants.EmptyField,
		Surname:      "",
		FiscalCode:   constants.EmptyField,
		VatCode:      constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		PostalCode:   constants.EmptyField,
		CityCode:     constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
		BirthDate:    constants.EmptyField,
	}
}

func (c *contractorDTO) fromPolicy(contractor models.Contractor) {
	if contractor.Name != "" {
		c.Name = contractor.Name
	}
	if contractor.Surname != "" {
		c.Surname = contractor.Surname
	}
	if contractor.CompanyName != "" {
		c.CompanyName = contractor.CompanyName
	}
	if contractor.FiscalCode != "" {
		c.FiscalCode = contractor.FiscalCode
	}
	if contractor.VatCode != "" {
		c.VatCode = contractor.VatCode
	}

	if contractor.Type == models.UserLegalEntity && contractor.CompanyAddress != nil {
		if len(contractor.CompanyAddress.StreetName) != 0 {
			c.StreetName = contractor.CompanyAddress.StreetName
		}
		if len(contractor.CompanyAddress.StreetNumber) != 0 {
			c.StreetNumber = contractor.CompanyAddress.StreetNumber
		}
		if len(contractor.CompanyAddress.PostalCode) != 0 {
			c.PostalCode = contractor.CompanyAddress.PostalCode
		}
		if len(contractor.CompanyAddress.City) != 0 {
			c.City = contractor.CompanyAddress.City
		}
		if len(contractor.CompanyAddress.CityCode) != 0 {
			c.CityCode = contractor.CompanyAddress.CityCode
		}
	} else if contractor.Residence != nil {
		if len(contractor.Residence.StreetName) != 0 {
			c.StreetName = contractor.Residence.StreetName
		}
		if len(contractor.Residence.StreetNumber) != 0 {
			c.StreetNumber = contractor.Residence.StreetNumber
		}
		if len(contractor.Residence.City) != 0 {
			c.City = contractor.Residence.City
		}
		if len(contractor.Residence.PostalCode) != 0 {
			c.PostalCode = contractor.Residence.PostalCode
		}
		if len(contractor.Residence.CityCode) != 0 {
			c.CityCode = contractor.Residence.CityCode
		}
		if len(contractor.Mail) != 0 {
			c.Mail = contractor.Mail
		}
		if len(contractor.Phone) != 0 {
			c.Phone = contractor.Phone
		}
	}
	c.Address = strings.ToUpper(c.StreetName + ", " + c.StreetNumber + "\n" + c.PostalCode + " " + c.City + " (" + c.CityCode + ")\n")

}
func (c *contractorDTO) GetFullNameContractor() (res string) {
	res = c.Name
	if c.Surname != "" {
		res += " " + c.Surname
	}
	if c.CompanyName != "" {
		if res != "" {
			res += ", "
		}
		res += c.CompanyName
	}
	return res
}

func (c *contractorDTO) GetFiscalCodeVatCode() string {
	var res = c.FiscalCode
	if c.FiscalCode != "" && c.VatCode != "" {
		res += "/"
	}
	if c.VatCode != "" {
		res += c.VatCode
	}
	return res
}
