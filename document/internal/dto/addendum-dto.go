package dto

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/models"
)

type addendumContractDTO struct {
	CodeHeading string
	Code        string
	StartDate   string
	EndDate     string
	Producer    string
}
type addendumContractorDTO struct {
	Name            string
	Surname         string
	FiscalCode      string
	StreetName      string
	StreetNumber    string
	City            string
	Province        string
	DomStreetName   string
	DomStreetNumber string
	DomCity         string
	DomProvince     string
	Mail            string
	Phone           string
	BirthDate       string
}

type addendumInsuredDTO struct {
	Name            string
	Surname         string
	FiscalCode      string
	StreetName      string
	StreetNumber    string
	City            string
	Province        string
	DomStreetName   string
	DomStreetNumber string
	DomCity         string
	DomProvince     string
	Mail            string
	Phone           string
	BirthDate       string
}

type addendumBeneficiaryDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Province     string
	Mail         string
	Relation     string
	Contactable  bool
	BirthDate    string
	Phone        string
}
type addendumBeneficiaryReferenceDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Province     string
	Mail         string
	Phone        string
	BirthDate    string
}

type addendumBeneficiaries []addendumBeneficiaryDTO

type AddendumBeneficiariesDTO struct {
	Contract             *addendumContractDTO
	Contractor           *addendumContractorDTO
	Insured              *addendumInsuredDTO
	Beneficiaries        *addendumBeneficiaries
	BeneficiaryReference *addendumBeneficiaryReferenceDTO
}

func NewBeneficiariesDto() *AddendumBeneficiariesDTO {
	return &AddendumBeneficiariesDTO{
		Contract:             newLifeContractDTO(),
		Contractor:           newLifeContractorDTO(),
		Insured:              newLifeInsuredDTO(),
		Beneficiaries:        newLifeBeneficiariesDTO(),
		BeneficiaryReference: newBeneficiaryReferenceDTO(),
	}
}

func (b *AddendumBeneficiariesDTO) FromPolicy(policy models.Policy, product models.Product) {
	b.Contract.fromPolicy(policy)
	b.Contractor.fromPolicy(policy.Contractor)
	for _, a := range policy.Assets {
		if a.Person != nil {
			b.Insured.fromPolicy(a.Person)
		}
		for _, g := range a.Guarantees {
			if g.Beneficiaries != nil {
				b.Beneficiaries.fromPolicy(g.Beneficiaries, g.BeneficiaryOptions)
			}
			if g.BeneficiaryReference != nil {
				b.BeneficiaryReference.fromPolicy(g.BeneficiaryReference)
			}
		}
	}
}

func newLifeContractDTO() *addendumContractDTO {
	return &addendumContractDTO{
		CodeHeading: constants.EmptyField,
		Code:        constants.EmptyField,
		StartDate:   constants.EmptyField,
		EndDate:     constants.EmptyField,
		Producer:    constants.EmptyField,
	}

}

func newLifeContractorDTO() *addendumContractorDTO {
	return &addendumContractorDTO{
		Name:            constants.EmptyField,
		Surname:         constants.EmptyField,
		FiscalCode:      constants.EmptyField,
		StreetName:      constants.EmptyField,
		StreetNumber:    constants.EmptyField,
		City:            constants.EmptyField,
		Province:        constants.EmptyField,
		DomStreetName:   constants.EmptyField,
		DomStreetNumber: constants.EmptyField,
		DomCity:         constants.EmptyField,
		DomProvince:     constants.EmptyField,
		Mail:            constants.EmptyField,
		Phone:           constants.EmptyField,
	}
}

func newLifeInsuredDTO() *addendumInsuredDTO {
	return &addendumInsuredDTO{
		Name:            constants.EmptyField,
		Surname:         constants.EmptyField,
		FiscalCode:      constants.EmptyField,
		StreetName:      constants.EmptyField,
		StreetNumber:    constants.EmptyField,
		City:            constants.EmptyField,
		Province:        constants.EmptyField,
		DomStreetName:   constants.EmptyField,
		DomStreetNumber: constants.EmptyField,
		DomCity:         constants.EmptyField,
		DomProvince:     constants.EmptyField,
		Mail:            constants.EmptyField,
		Phone:           constants.EmptyField,
	}
}

func newLifeBeneficiariesDTO() *addendumBeneficiaries {
	lb := make(addendumBeneficiaries, 0)
	l := addendumBeneficiaryDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Province:     constants.EmptyField,
		Mail:         constants.EmptyField,
		Relation:     constants.EmptyField,
		BirthDate:    constants.EmptyField,
		Phone:        constants.EmptyField,
	}
	lb = append(lb, l)
	lb = append(lb, l)
	return &lb
}

func newBeneficiaryReferenceDTO() *addendumBeneficiaryReferenceDTO {
	return &addendumBeneficiaryReferenceDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Province:     constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
		BirthDate:    constants.EmptyField,
	}
}

func (l *addendumContractDTO) fromPolicy(policy models.Policy) {
	l.CodeHeading = "Variazione dati Anagrafici soggetti Polizza:"
	l.Code = policy.CodeCompany

	if !policy.StartDate.IsZero() {
		l.StartDate = policy.StartDate.Format(constants.DayMonthYearFormat)
	}

	if !policy.EndDate.IsZero() {
		l.EndDate = policy.EndDate.Format(constants.DayMonthYearFormat)
	}

	l.Producer = policy.Company

}

func parseBirthDate(dateString string) string {
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return ""
	}
	return date.Format("01/02/2006")
}

func (lc *addendumContractorDTO) fromPolicy(contr models.Contractor) {
	if contr.FiscalCode != "" {
		lc.Name = contr.Name
		lc.Surname = contr.Surname
		lc.FiscalCode = contr.FiscalCode
		lc.BirthDate = parseBirthDate(contr.BirthDate)
	}
	if contr.Residence != nil {
		lc.StreetName = contr.Residence.StreetName
		lc.StreetNumber = contr.Residence.StreetNumber
		lc.City = contr.Residence.City
		lc.Province = contr.Residence.CityCode
	}
	if contr.Domicile != nil {
		lc.DomStreetName = contr.Domicile.StreetName
		lc.DomStreetNumber = contr.Domicile.StreetNumber
		lc.DomCity = contr.Domicile.City
		lc.DomProvince = contr.Domicile.CityCode
	}
	if contr.Mail != "" {
		lc.Mail = contr.Mail
	}
	if contr.Phone != "" {
		lc.Phone = contr.Phone
	}
}

func (li *addendumInsuredDTO) fromPolicy(ins *models.User) {
	if ins != nil {
		if ins.FiscalCode != "" {
			li.Name = ins.Name
			li.Surname = ins.Surname
			li.FiscalCode = ins.FiscalCode
			li.BirthDate = parseBirthDate(ins.BirthDate)
		}
		if ins.Residence != nil {
			li.StreetName = ins.Residence.StreetName
			li.StreetNumber = ins.Residence.StreetNumber
			li.City = ins.Residence.City
			li.Province = ins.Residence.CityCode
		}

		if ins.Domicile != nil {
			li.DomStreetName = ins.Domicile.StreetName
			li.DomStreetNumber = ins.Domicile.StreetNumber
			li.DomCity = ins.Domicile.City
			li.DomProvince = ins.Domicile.CityCode
		}
		if ins.Mail != "" {
			li.Mail = ins.Mail
		}
		if ins.Phone != "" {
			li.Phone = ins.Phone
		}
	}
}

func (b *addendumBeneficiaries) fromPolicy(bens *[]models.Beneficiary, opt map[string]string) {
	if bens == nil {
		return
	}
	for i, v := range *bens {
		if i > 1 {
			break
		}
		buf := new(bytes.Buffer)
		for _, value := range opt {
			_, _ = fmt.Fprintf(buf, "%s ", value)
		}
		ben := addendumBeneficiaryDTO{
			Name:         v.Name,
			Surname:      v.Surname,
			FiscalCode:   v.FiscalCode,
			StreetName:   v.Residence.StreetName,
			StreetNumber: v.Residence.StreetNumber,
			City:         v.Residence.City,
			Province:     v.Residence.CityCode,
			Phone:        v.Phone,
			Mail:         v.Mail,
			BirthDate:    (*b)[i].BirthDate,
			Contactable:  v.IsContactable,
			Relation:     "\n" + buf.String(),
		}
		(*b)[i] = ben
	}
}

func (br *addendumBeneficiaryReferenceDTO) fromPolicy(benRef *models.User) {
	if benRef != nil {
		if benRef.FiscalCode != "" {
			br.Name = benRef.Name
			br.Surname = benRef.Surname
			br.FiscalCode = benRef.FiscalCode
		}
		if benRef.Residence != nil {
			br.StreetName = benRef.Residence.StreetName
			br.StreetNumber = benRef.Residence.StreetNumber
			br.City = benRef.Residence.City
			br.Province = benRef.Residence.CityCode
		}
		if benRef.Mail != "" {
			br.Mail = benRef.Mail
		}
		if benRef.Phone != "" {
			br.Phone = benRef.Phone
		}
	}
}
