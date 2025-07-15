package _script

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
	"gitlab.dev.wopta.it/goworkspace/user"
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
	nodeList      = []struct {
		Uid  string
		Code string
		Type string
	}{
		{"0aDBMMGM83xtRNE07ZYh", "W1.TestModifica", "agent"}, // mga_life_agent
		{"qKSQI7AZgHzzS2EE0dWV", "DSC.REM", "agent"},         // mga_life_agent_remittance
	}
)

func getFakePerson() models.User {
	nameIdx := rand.Int() % 10
	surnameIdx := rand.Int() % 10
	dateOfBirthIdx := rand.Int() % 10
	genderIdx := rand.Int() % 2

	u := models.User{
		Name:          nameList[nameIdx],
		Surname:       surnameList[surnameIdx],
		BirthDate:     birthDateList[dateOfBirthIdx],
		Gender:        genderList[genderIdx],
		BirthCity:     "MILANO",
		BirthProvince: "MI",
	}

	_, u, _ = user.CalculateFiscalCodeInUser(u)

	return u
}

func getFakeProducer() struct {
	Uid  string
	Code string
	Type string
} {
	producerIdx := rand.Int() % 2

	return nodeList[producerIdx]
}

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

	producer := getFakeProducer()

	p.ProducerCode = producer.Code
	p.ProducerUid = producer.Uid
	p.ProducerType = producer.Type

	return p
}

func SeedPolicies(jsonFilepath string) error {
	var (
		filePolicies []models.Policy
		batch        map[string]map[string]any = make(map[string]map[string]any)
	)

	batch[lib.PolicyCollection] = make(map[string]any)
	batch[lib.TransactionsCollection] = make(map[string]any)

	now := time.Now().UTC()

	fileReader, err := os.Open(jsonFilepath)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(fileReader).Decode(&filePolicies); err != nil {
		return err
	}

	mgaProducts := map[string]*models.Product{
		models.ProductV1: product.GetProductV2(models.LifeProduct, models.ProductV1, models.MgaChannel, nil,
			nil),
		models.ProductV2: product.GetProductV2(models.LifeProduct, models.ProductV2, models.MgaChannel, nil,
			nil),
	}

	for _, p := range filePolicies {
		p.Uid = lib.NewDoc(lib.PolicyCollection)
		p.CodeCompany = "99" + p.CodeCompany

		p = anonimizePolicy(p)
		p.BigQueryParse()

		batch[lib.PolicyCollection][p.Uid] = p

		trs := transaction.CreateTransactions(p, *mgaProducts[p.ProductVersion], func() string { return lib.NewDoc(lib.TransactionsCollection) })
		for i, tr := range trs {
			trs[i].IsPay = tr.EffectiveDate.Before(now)
			if trs[i].IsPay {
				trs[i].PayDate = now
				trs[i].PaymentMethod = "import-seed"
				trs[i].Status = models.TransactionStatusPay
				trs[i].StatusHistory = append(trs[i].StatusHistory, trs[i].Status)
			}
			trs[i].BigQueryParse()
			batch[lib.TransactionsCollection][tr.Uid] = trs[i]
		}
	}

	if err := lib.SetBatchFirestoreErr(batch); err != nil {
		return err
	}

	for col := range batch {
		data := lib.GetMapValues(batch[col])

		if err := lib.InsertRowsBigQuery(lib.WoptaDataset, col, data); err != nil {
			return err
		}
	}

	return nil
}
