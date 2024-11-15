package _script

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/user"
)

func getFakeAddress() *models.Address {
	return &models.Address{
		StreetName:   "GALLERIA DEL CORSO",
		StreetNumber: "1",
		City:         "MILANO",
		PostalCode:   "20122",
		Locality:     "MILANO",
		CityCode:     "MI",
		Area:         "",
	}
}

var (
	nameList      = []string{"MARIO", "PAOLO", "GIUSEPPE", "YOUSEF", "DIOGO", "LUCA", "ELIZABETH", "ALBERTO", "IVAN", "MIHAELA"}
	surnameList   = []string{"ROSSO", "VERDE", "GIALLO", "NERO", "BLU", "LILLA", "BIANCO", "MARRONE", "ARANCIONE", "AZZURRO"}
	birthDateList = []string{"1980-01-01T00:00:00Z", "1980-02-01T00:00:00Z", "1980-03-01T00:00:00Z", "1980-04-01T00:00:00Z", "1980-05-01T00:00:00Z", "1980-06-01T00:00:00Z", "1980-07-01T00:00:00Z", "1980-08-01T00:00:00Z", "1980-09-01T00:00:00Z", "1980-10-01"}
	genderList    = []string{"M", "F"}
)

func getFakePerson() models.User {
	nameIdx := (counter + time.Now().Nanosecond()) % 10
	surnameIdx := (counter + time.Now().Nanosecond()) % 10
	dateOfBirthIdx := (counter + time.Now().Nanosecond()) % 10
	genderIdx := (counter + time.Now().Nanosecond()) % 2

	u := models.User{
		Name:          nameList[nameIdx],
		Surname:       surnameList[surnameIdx],
		BirthDate:     birthDateList[dateOfBirthIdx],
		Gender:        genderList[genderIdx],
		BirthCity:     "MILANO",
		BirthProvince: "MI",
	}

	_, u, _ = user.CalculateFiscalCode(u)

	counter++

	return u
}

var counter int

func anonimizePolicy(p models.Policy) models.Policy {
	fakeAddress := getFakeAddress()
	fakeInsured := getFakePerson()
	fakeContractor := fakeInsured

	isContractorInsured := false

	for _, a := range p.Assets {
		if a.Person != nil {
			a.Person.Domicile = fakeAddress
			a.Person.Residence = fakeAddress
			a.Person.IdentityDocuments = nil
			a.Person.Name = fakeInsured.Name
			a.Person.Surname = fakeInsured.Surname
			a.Person.BirthCity = fakeInsured.BirthCity
			a.Person.BirthProvince = fakeInsured.BirthProvince
			a.Person.BirthDate = fakeInsured.BirthDate
			a.Person.FiscalCode = fakeInsured.FiscalCode
		}
		if a.Guarantees != nil {
			for i, g := range a.Guarantees {
				if g.Beneficiaries != nil {
					for j, b := range *g.Beneficiaries {
						bn := b
						if bn.BeneficiaryType != "GE" {
							fakeBeneficiary := getFakePerson()
							bn.Name = fakeBeneficiary.Name
							bn.Surname = fakeBeneficiary.Surname
							bn.FiscalCode = fakeBeneficiary.FiscalCode
							bn.Mail = ""
							bn.Phone = ""
							if bn.CompanyAddress != nil {
								bn.CompanyAddress = fakeAddress
							}
						}
						(*a.Guarantees[i].Beneficiaries)[j] = bn
					}
				}
			}
		}
		isContractorInsured = a.IsContractor
	}

	if !isContractorInsured {
		fakeContractor = getFakePerson()
	}

	if p.Contractor.Type == "legalEntity" {
		p.Contractor.CompanyAddress = fakeAddress
		p.Contractor.VatCode = "12345678910"
		p.Contractor.Name = "ACME Srl"
	} else {
		p.Contractor.Name = fakeContractor.Name
		p.Contractor.Surname = fakeContractor.Surname
		p.Contractor.BirthCity = fakeContractor.BirthCity
		p.Contractor.BirthProvince = fakeContractor.BirthProvince
		p.Contractor.BirthDate = fakeContractor.BirthDate
		p.Contractor.FiscalCode = fakeContractor.FiscalCode
		p.Contractor.Domicile = fakeAddress
		p.Contractor.Residence = fakeAddress
	}
	
	p.Contractor.IdentityDocuments = nil
	p.Contractor.Phone = "+393334455667"
	p.Contractor.Mail = fmt.Sprintf("DIOGO.CARVALHO+SEED%s@WOPTA.IT", p.CodeCompany)

	if p.Contractors != nil {
		for i, c := range *p.Contractors {
			cn := c
			fakeContractor = getFakePerson()
			cn.Name = fakeContractor.Name
			cn.Surname = fakeContractor.Surname
			cn.BirthCity = fakeContractor.BirthCity
			cn.BirthProvince = fakeContractor.BirthProvince
			cn.BirthDate = fakeContractor.BirthDate
			cn.FiscalCode = fakeContractor.FiscalCode
			cn.Domicile = fakeAddress
			cn.Residence = fakeAddress
			cn.IdentityDocuments = nil
			(*p.Contractors)[i] = cn
		}
	}

	p.ProducerCode = ""
	p.ProducerUid = ""
	p.ProducerType = ""

	return p
}

func SeedPolicies(jsonFilepath string) error {
	var (
		filePolicies []models.Policy
	)

	fileReader, err := os.Open(jsonFilepath)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(fileReader).Decode(&filePolicies); err != nil {
		return err
	}

	anonPolicies := make([]models.Policy, 0,len(filePolicies))
	
	for _, p := range filePolicies {
		p.Uid = lib.NewDoc(lib.PolicyCollection)
		p.CodeCompany = "99" + p.CodeCompany

		p = anonimizePolicy(p)

		anonPolicies = append(anonPolicies, p)
	}

	anonBytes, err := json.Marshal(anonPolicies)
	if err != nil {
		return err
	}

	return os.WriteFile("./_script/anonPolicies.json", anonBytes, os.ModePerm)
}
