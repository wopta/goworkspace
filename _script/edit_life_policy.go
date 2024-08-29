package _script

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var identityDocumentMap = map[string]string{
	"01": "Carta di Identità",
	"02": "Patente di Guida",
	"03": "Passaporto",
}

func EditLifePolicy(policyUid string) {
	rawData, err := os.ReadFile("./_script/policy_80.txt")
	if err != nil {
		log.Fatal(err)
	}

	df, err := lib.CsvToDataframeV2(rawData, ';', false)
	if err != nil {
		log.Fatal(err)
	}

	groups := groupBy(df, 2)
	delete(groups, "X2")

	policy := fetchPolicy(policyUid)

	for _, rawPolicy := range groups {
		log.Printf("%v", rawPolicy)

		// TODO: implementare estrazione dati contraente persona giuridica
		contractor := extractContractorData(rawPolicy[0])
		// TODO: implementare estrazione dati assicurato
		insured := extractInsuredData(rawPolicy[0], policy)
		// TODO: implementare estrazione dati 3 titolari effettivi

		log.Printf("contractor: %+v", contractor)
		log.Printf("insured: %+v", insured)
	}

}

func fetchPolicy(policyUid string) models.Policy {
	var policy models.Policy

	docsnap, err := lib.GetFirestoreErr(lib.PolicyCollection, policyUid)
	if err != nil {
		log.Fatal(err)
	}
	err = docsnap.DataTo(&policy)
	if err != nil {
		log.Fatal(err)
	}
	return policy
}
func extractContractorData(rawPolicy []string) models.Contractor {
	contractor := models.Contractor{}

	now := time.Now().UTC()
	phone := strings.TrimSpace(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(rawPolicy[33], " ", ""), " ", "")))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}
	vatCode := fmt.Sprintf("%011s", strings.TrimSpace(rawPolicy[27]))

	contractor.Uid = vatCode
	contractor.Type = models.UserLegalEntity
	contractor.Name = strings.TrimSpace(lib.Capitalize(rawPolicy[23]))
	contractor.VatCode = vatCode
	contractor.Mail = strings.TrimSpace(lib.Capitalize(rawPolicy[32]))
	contractor.Phone = phone
	contractor.Consens = &[]models.Consens{
		{
			Title:        "Privacy",
			Consens:      "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
			Key:          2,
			Answer:       false,
			CreationDate: now,
		},
	}
	contractor.CompanyAddress = &models.Address{
		StreetName: strings.TrimSpace(lib.Capitalize(rawPolicy[28])),
		City:       strings.TrimSpace(lib.Capitalize(rawPolicy[30])),
		CityCode:   strings.TrimSpace(strings.ToUpper(rawPolicy[31])),
		PostalCode: strings.TrimSpace(rawPolicy[29]),
		Locality:   strings.TrimSpace(lib.Capitalize(rawPolicy[30])),
	}
	contractor.CreationDate = now
	contractor.UpdatedDate = now
	contractor.Normalize()
	return contractor
}

func extractInsuredData(rawPolicy []string, policy models.Policy) models.User {
	insured := models.User{}

	now := time.Now().UTC()
	phone := strings.TrimSpace(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(rawPolicy[72], " ", ""), " ", "")))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}
	rawCode, err := strconv.Atoi(strings.TrimSpace(rawPolicy[76]))
	if err != nil {
		log.Fatal(err)
	}
	identityDocumentCode := fmt.Sprintf("%02d", rawCode)

	insured.Uid = policy.Assets[0].Person.Uid
	insured.Type = models.UserIndividual
	insured.Name = strings.TrimSpace(lib.Capitalize(rawPolicy[35]))
	insured.Surname = strings.TrimSpace(lib.Capitalize(rawPolicy[34]))
	insured.FiscalCode = strings.TrimSpace(lib.Capitalize(rawPolicy[38]))
	insured.Gender = strings.TrimSpace(lib.ToUpper(rawPolicy[36]))
	insured.Mail = strings.TrimSpace(lib.ToUpper(rawPolicy[71]))
	insured.Phone = phone
	insured.BirthDate = parseDate(rawPolicy[37]).Format(time.RFC3339)
	insured.BirthCity = strings.TrimSpace(lib.Capitalize(rawPolicy[73]))
	insured.BirthProvince = strings.TrimSpace(lib.ToUpper(rawPolicy[74]))
	insured.Residence = &models.Address{
		StreetName: strings.TrimSpace(lib.Capitalize(rawPolicy[63])),
		City:       strings.TrimSpace(lib.Capitalize(rawPolicy[65])),
		CityCode:   strings.TrimSpace(strings.ToUpper(rawPolicy[66])),
		PostalCode: strings.TrimSpace(rawPolicy[64]),
		Locality:   strings.TrimSpace(lib.Capitalize(rawPolicy[65])),
	}
	insured.Domicile = &models.Address{
		StreetName: strings.TrimSpace(lib.Capitalize(rawPolicy[67])),
		City:       strings.TrimSpace(lib.Capitalize(rawPolicy[69])),
		CityCode:   strings.TrimSpace(strings.ToUpper(rawPolicy[70])),
		PostalCode: strings.TrimSpace(rawPolicy[68]),
		Locality:   strings.TrimSpace(lib.Capitalize(rawPolicy[69])),
	}
	insured.Consens = &[]models.Consens{
		{
			Title:        "Privacy",
			Consens:      "Il sottoscritto, letta e compresa l'informativa sul trattamento dei dati personali, ACCONSENTE al trattamento dei propri dati personali da parte di Wopta Assicurazioni per l'invio di comunicazioni e proposte commerciali e di marketing, incluso l'invio di newsletter e ricerche di mercato, attraverso strumenti automatizzati (sms, mms, e-mail, ecc.) e non (posta cartacea e telefono con operatore).",
			Key:          2,
			Answer:       false,
			CreationDate: now,
		},
	}
	insured.IdentityDocuments = []*models.IdentityDocument{
		{
			Number:           strings.TrimSpace(strings.ToUpper(rawPolicy[77])),
			Code:             identityDocumentCode,
			Type:             identityDocumentMap[identityDocumentCode],
			DateOfIssue:      parseDate(rawPolicy[78]),
			ExpiryDate:       parseDate(rawPolicy[78]).AddDate(10, 0, 0),
			IssuingAuthority: strings.TrimSpace(lib.Capitalize(rawPolicy[79])),
			PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(rawPolicy[79])),
			LastUpdate:       policy.EmitDate,
		},
	}
	insured.CreationDate = now
	insured.UpdatedDate = now
	insured.Normalize()
	return insured
}

func parseDate(rawDate string) time.Time {
	day, _ := strconv.Atoi(strings.TrimSpace(rawDate[:2]))
	month, _ := strconv.Atoi(strings.TrimSpace(rawDate[2:4]))
	year, _ := strconv.Atoi(strings.TrimSpace(rawDate[4:]))

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return date
}
