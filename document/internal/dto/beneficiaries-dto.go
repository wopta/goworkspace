package dto

import (
	"strconv"

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
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Mail         string
	Phone        string
	BirthDate    string
}

type lifeInsuredDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Mail         string
	Phone        string
	BirthDate    string
}

type lifeBeneficiaryDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Mail         string
	Relation     string
	Contactable  bool
	BirthDate    string
	Phone        string
}

type Beneficiaries []lifeBeneficiaryDTO

type BeneficiariesDTO struct {
	Contract             *lifeContractDTO
	Contractor           *lifeContractorDTO
	Insured              *lifeInsuredDTO
	Beneficiaries        *Beneficiaries
	BeneficiaryReference *BeneficiaryReferenceDTO
}

type BeneficiaryReferenceDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	Mail         string
	Phone        string
	BirthDate    string
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
	b.Beneficiaries.fromPolicy(policy.Assets[0].Guarantees[0].Beneficiaries)
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
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
	}
}

func newLifeInsuredDTO() *lifeInsuredDTO {
	return &lifeInsuredDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
	}
}

func newLifeBeneficiariesDTO() *Beneficiaries {
	lb := make([]lifeBeneficiaryDTO, 0)
	l := lifeBeneficiaryDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Mail:         constants.EmptyField,
		Relation:     constants.EmptyField,
		BirthDate:    constants.EmptyField,
		Phone:        constants.EmptyField,
	}
	lb = append(lb, l)
	lb = append(lb, l)
	var b Beneficiaries = lb
	return &b
}

func newBeneficiaryReferenceDTO() *BeneficiaryReferenceDTO {
	return &BeneficiaryReferenceDTO{
		Name:         constants.EmptyField,
		Surname:      constants.EmptyField,
		FiscalCode:   constants.EmptyField,
		StreetName:   constants.EmptyField,
		StreetNumber: constants.EmptyField,
		City:         constants.EmptyField,
		Mail:         constants.EmptyField,
		Phone:        constants.EmptyField,
		BirthDate:    constants.EmptyField,
	}
}

func (l *lifeContractDTO) fromPolicy(policy models.Policy) {
	l.CodeHeading = "Variazione dati Anagrafici soggetti Polizza:"
	l.Code = strconv.Itoa(policy.Number)

	if !policy.StartDate.IsZero() {
		l.StartDate = policy.StartDate.Format(constants.DayMonthYearFormat)
	}

	if !policy.EndDate.IsZero() {
		l.EndDate = policy.EndDate.Format(constants.DayMonthYearFormat)
	}

	l.Producer = policy.Company

}

func (lc *lifeContractorDTO) fromPolicy(contr models.Contractor) {
	lc.Name = contr.Name
	lc.Surname = contr.Surname
	lc.FiscalCode = contr.FiscalCode
	lc.StreetName = contr.Residence.StreetName
	lc.StreetNumber = contr.Residence.StreetNumber
	lc.City = contr.Residence.City
	lc.Mail = contr.Mail
	lc.Phone = contr.Phone
	lc.BirthDate = contr.BirthDate
}

func (li *lifeInsuredDTO) fromPolicy(ins *models.User) {
	li.Name = ins.Name
	li.Surname = ins.Surname
	li.FiscalCode = ins.FiscalCode
	li.StreetName = ins.Residence.StreetName
	li.StreetNumber = ins.Residence.StreetNumber
	li.City = ins.Residence.City
	li.Mail = ins.Mail
	li.Phone = ins.Phone
	li.BirthDate = ins.BirthDate
}

func (b *Beneficiaries) fromPolicy(bens *[]models.Beneficiary) {

	for i, v := range *bens {
		if i > 1 {
			break
		}
		ben := lifeBeneficiaryDTO{
			Name:         v.Name,
			Surname:      v.Surname,
			FiscalCode:   v.FiscalCode,
			StreetName:   v.Residence.StreetName,
			StreetNumber: v.Residence.StreetNumber,
			City:         v.Residence.City,
			Phone:        v.Phone,
			Mail:         v.Mail,
			BirthDate:    (*b)[i].BirthDate,
			Contactable:  v.IsContactable,
		}
		(*b)[i] = ben
	}
}

func (br *BeneficiaryReferenceDTO) fromPolicy(benRef *models.User) {
	br.Name = benRef.Name
	br.Surname = benRef.Surname
	br.FiscalCode = benRef.FiscalCode
	br.StreetName = benRef.Residence.StreetName
	br.StreetNumber = benRef.Residence.StreetNumber
	br.City = benRef.Residence.City
	br.Mail = benRef.Mail
	br.Phone = benRef.Phone
}
