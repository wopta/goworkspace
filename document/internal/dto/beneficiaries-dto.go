package dto

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wopta/goworkspace/document/internal/constants"
	"github.com/wopta/goworkspace/models"
)

type lifeContractDTO struct {
	CodeHeading string
	Code        string
	StartDate   string
	EndDate     string
	Producer    string
}
type lifeContractorDTO struct {
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

type lifeInsuredDTO struct {
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

type lifeBeneficiaryDTO struct {
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
type BeneficiaryReferenceDTO struct {
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

type Beneficiaries []lifeBeneficiaryDTO

type BeneficiariesDTO struct {
	Contract             *lifeContractDTO
	Contractor           *lifeContractorDTO
	Insured              *lifeInsuredDTO
	Beneficiaries        *Beneficiaries
	BeneficiaryReference *BeneficiaryReferenceDTO
}

func NewBeneficiariesDto() *BeneficiariesDTO {
	return &BeneficiariesDTO{
		Contract:             newLifeContractDTO(),
		Contractor:           newLifeContractorDTO(),
		Insured:              newLifeInsuredDTO(),
		Beneficiaries:        newLifeBeneficiariesDTO(),
		BeneficiaryReference: newBeneficiaryReferenceDTO(),
	}
}

func (b *BeneficiariesDTO) FromPolicy(policy models.Policy, product models.Product) {
	b.Contract.fromPolicy(policy)
	b.Contractor.fromPolicy(policy.Contractor)
	b.Insured.fromPolicy(policy.Assets[0].Person)
	//b.Beneficiaries.fromPolicy(policy.Assets[0].Guarantees)
	b.Beneficiaries.fromPolicy(policy.Assets[0].Guarantees[0].Beneficiaries, policy.Assets[0].Guarantees[0].BeneficiaryOptions)
	b.BeneficiaryReference.fromPolicy(policy.Assets[0].Guarantees[0].BeneficiaryReference)
}

func newLifeContractDTO() *lifeContractDTO {
	return &lifeContractDTO{
		CodeHeading: constants.EmptyField,
		Code:        constants.EmptyField,
		StartDate:   constants.EmptyField,
		EndDate:     constants.EmptyField,
		Producer:    constants.EmptyField,
	}

}

func newLifeContractorDTO() *lifeContractorDTO {
	return &lifeContractorDTO{
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

func newLifeInsuredDTO() *lifeInsuredDTO {
	return &lifeInsuredDTO{
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

func newLifeBeneficiariesDTO() *Beneficiaries {
	lb := make(Beneficiaries, 0)
	l := lifeBeneficiaryDTO{
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

func newBeneficiaryReferenceDTO() *BeneficiaryReferenceDTO {
	return &BeneficiaryReferenceDTO{
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

func (l *lifeContractDTO) fromPolicy(policy models.Policy) {
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

func (lc *lifeContractorDTO) fromPolicy(contr models.Contractor) {
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

func (li *lifeInsuredDTO) fromPolicy(ins *models.User) {
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

func (b *Beneficiaries) fromPolicy(bens *[]models.Beneficiary, opt map[string]string) {
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
		ben := lifeBeneficiaryDTO{
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

func (br *BeneficiaryReferenceDTO) fromPolicy(benRef *models.User) {
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
