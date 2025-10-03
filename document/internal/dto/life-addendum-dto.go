package dto

import (
	"fmt"
	"time"

	"gitlab.dev.wopta.it/goworkspace/document/internal/constants"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type addendumContractDTO struct {
	CodeHeading string
	Code        string
	StartDate   string
	EndDate     string
	Producer    string
	IssueDate   string
}

type addendumPersonDTO struct {
	Name            string
	Surname         string
	FiscalCode      string
	StreetName      string
	StreetNumber    string
	City            string
	BirthCity       string
	BirthProvice    string
	PostalCode      string
	Province        string
	DomStreetName   string
	DomStreetNumber string
	DomCity         string
	DomPostalCode   string
	DomProvince     string
	Mail            string
	Phone           string
	BirthDate       string
	Gender          string
}

func (a addendumPersonDTO) GetDomicilioAddress() string {
	if a.DomStreetName == "" {
		return ""
	}
	return a.DomStreetName + " " + a.DomStreetNumber + ", " + a.DomPostalCode + " " + a.DomCity + " (" + a.DomProvince + ")"
}

func (a addendumPersonDTO) GetResidenceAddress() string {
	if a.StreetName == "" {
		return ""
	}
	return a.StreetName + " " + a.StreetNumber + ", " + a.PostalCode + " " + a.City + " (" + a.Province + ")"
}
func (a addendumPersonDTO) GetBirthAddress() string {
	if a.StreetName == "" {
		return ""
	}
	return a.BirthCity + " " + a.BirthProvice
}

type addendumBeneficiaryDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	PostalCode   string
	Province     string
	Mail         string
	Relation     string
	Contactable  bool
	BirthDate    string
	Phone        string
}

func (a addendumBeneficiaryDTO) GetResidenceAddress() string {
	if a.StreetName == "" {
		return ""
	}
	return a.StreetName + " " + a.StreetNumber + ", " + a.PostalCode + " " + a.City + " (" + a.Province + ")"
}

type addendumBeneficiaryReferenceDTO struct {
	Name         string
	Surname      string
	FiscalCode   string
	StreetName   string
	StreetNumber string
	City         string
	PostalCode   string
	Province     string
	Mail         string
	Phone        string
	BirthDate    string
}

type addendumBeneficiaries []addendumBeneficiaryDTO

type AddendumLifeDTO struct {
	Contract             *addendumContractDTO
	Contractor           *ContractorDTO
	Insured              *addendumPersonDTO
	Beneficiaries        *addendumBeneficiaries
	BeneficiaryReference *addendumBeneficiaryReferenceDTO
}

func NewLifeAddendumDto() *AddendumLifeDTO {
	l := make(addendumBeneficiaries, 2)
	return &AddendumLifeDTO{
		Contract:             &addendumContractDTO{},
		Contractor:           &ContractorDTO{},
		Insured:              &addendumPersonDTO{},
		Beneficiaries:        &l,
		BeneficiaryReference: &addendumBeneficiaryReferenceDTO{},
	}
}

func (b *AddendumLifeDTO) FromPolicy(policy *models.Policy, now time.Time) {
	b.Contract.fromPolicy(policy, now)
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

func (l *addendumContractDTO) fromPolicy(policy *models.Policy, now time.Time) {
	l.CodeHeading = "Variazione dati Anagrafici soggetti Polizza:"
	l.Code = policy.CodeCompany
	if !policy.CompanyEmit {
		l.CodeHeading = "Variazione dati Anagrafici soggetti Proposta:"
		l.Code = fmt.Sprintf("%d", policy.ProposalNumber)
	}

	location, _ := time.LoadLocation("Europe/Rome")
	l.IssueDate = "Milano, il " + now.In(location).Format(constants.DayMonthYearFormat)

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

var genderToIta = map[string]string{
	"M": "Maschile",
	"F": "Femminile",
}

func (li *addendumPersonDTO) fromPolicy(ins *models.User) {
	if ins == nil {
		return
	}
	li.BirthCity = ins.BirthCity
	li.BirthProvice = ins.BirthProvince
	li.BirthDate = ins.BirthDate
	li.Gender = genderToIta[ins.Gender]
	if ins.FiscalCode != "" {
		li.Name = ins.Name
		li.Surname = ins.Surname
		li.FiscalCode = ins.FiscalCode
		birth := parseBirthDate(ins.BirthDate)
		if birth != "" {
			li.BirthDate = birth
		}
	}
	if ins.Residence != nil {
		li.StreetName = ins.Residence.StreetName
		li.StreetNumber = ins.Residence.StreetNumber
		li.City = ins.Residence.City
		li.PostalCode = ins.Residence.PostalCode
		li.Province = ins.Residence.CityCode
	}

	if ins.Domicile != nil {
		li.DomStreetName = ins.Domicile.StreetName
		li.DomStreetNumber = ins.Domicile.StreetNumber
		li.DomCity = ins.Domicile.City
		li.DomPostalCode = ins.Domicile.PostalCode
		li.DomProvince = ins.Domicile.CityCode
	}
	if ins.Mail != "" {
		li.Mail = ins.Mail
	}
	if ins.Phone != "" {
		li.Phone = ins.Phone
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

		// TODO: improve me - shouldnt be needed
		ben := addendumBeneficiaryDTO{}
		ben.Relation = " \n" + opt[v.BeneficiaryType]

		if v.BeneficiaryType == models.BeneficiaryChosenBeneficiary {
			if v.Name != "" {
				ben.Name = v.Name
			}
			if v.Surname != "" {
				ben.Surname = v.Surname
			}
			if v.FiscalCode != "" {
				ben.FiscalCode = v.FiscalCode
			}
			if v.FiscalCode != "" {
				ben.FiscalCode = v.FiscalCode
			}
			if v.Phone != "" {
				ben.Phone = v.Phone
			}
			if v.Mail != "" {
				ben.Mail = v.Mail
			}
			ben.Contactable = v.IsContactable
			if v.Residence != nil {
				ben.StreetName = v.Residence.StreetName
				ben.StreetNumber = v.Residence.StreetNumber
				ben.City = v.Residence.City
				ben.PostalCode = v.Residence.PostalCode
				ben.Province = v.Residence.CityCode
			}
		}
		(*b)[i] = ben
	}
}

func (br *addendumBeneficiaryReferenceDTO) fromPolicy(benRef *models.User) {
	*br = addendumBeneficiaryReferenceDTO{}
	if benRef.FiscalCode != "" {
		br.Name = benRef.Name
		br.Surname = benRef.Surname
		br.FiscalCode = benRef.FiscalCode
	}
	if benRef.Residence != nil {
		br.StreetName = benRef.Residence.StreetName
		br.StreetNumber = benRef.Residence.StreetNumber
		br.City = benRef.Residence.City
		br.PostalCode = benRef.Residence.PostalCode
		br.Province = benRef.Residence.CityCode
	}
	if benRef.Mail != "" {
		br.Mail = benRef.Mail
	}
	if benRef.Phone != "" {
		br.Phone = benRef.Phone
	}
}
