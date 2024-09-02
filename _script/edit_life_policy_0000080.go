package _script

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var identityDocumentMap = map[string]string{
	"01": "Carta di Identità",
	"02": "Patente di Guida",
	"03": "Passaporto",
}

func EditLifePolicy0000080(policyUid string) {
	rawData, err := os.ReadFile("./_script/policy_80.txt")
	if err != nil {
		log.Fatal(err)
	}

	df, err := csvToDataframe(rawData, ';', false)
	if err != nil {
		log.Fatal(err)
	}

	groups := groupBy(df, 2)
	delete(groups, "X2")

	policy := fetchPolicy(policyUid)

	rawPolicy := groups["80"][0]

	contractor := extractContractorData(rawPolicy)
	insured := extractInsuredData(rawPolicy, policy)
	contractors := extractContractorsData(rawPolicy)

	policy.Contractor = contractor
	policy.Contractors = &contractors
	policy.Assets[0].Person = &insured

	err = saveUser(&insured)
	if err != nil {
		log.Fatal(err)
	}

	err = saveUser(policy.Contractor.ToUser())
	if err != nil {
		log.Fatal(err)
	}

	for _, contr := range contractors {
		user := contr
		err = saveUser(&user)
		if err != nil {
			log.Fatal(err)
		}
	}

	policy.Updated = time.Now().UTC()

	err = lib.SetFirestoreErr(lib.PolicyCollection, policy.Uid, policy)
	if err != nil {
		log.Fatal(err)
	}

	policy.BigquerySave("")
}

func csvToDataframe(data []byte, delimiter rune, hasHeader bool) (dataframe.DataFrame, error) {
	reader := bytes.NewReader(data)
	df := dataframe.ReadCSV(reader,
		dataframe.WithDelimiter(delimiter),
		dataframe.HasHeader(hasHeader),
		dataframe.NaNValues(nil),
		dataframe.DetectTypes(false))
	return df, df.Error()
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

func extractContractorsData(rawPolicy []string) []models.User {
	const offset = 26
	contractors := make([]models.User, 0)

	for i := 0; i < 3; i++ {
		if strings.TrimSpace(strings.ToUpper(rawPolicy[116+(offset*i)])) == "" || strings.TrimSpace(strings.ToUpper(rawPolicy[116+(offset*i)])) == "NO" {
			break
		}
		titolareEffettivo := parsingTitolareEffettivo(rawPolicy, offset, i)
		titolareEffettivo.Normalize()
		contractors = append(contractors, titolareEffettivo)
	}
	return contractors
}

func parsingTitolareEffettivo(row []string, offset int, i int) models.User {
	isExecutor := strings.TrimSpace(strings.ToUpper(row[224])) == strings.TrimSpace(strings.ToUpper(row[121+(offset*i)]))
	now := time.Now().UTC()
	phone := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(row[132], " ", ""), " ", ""))
	if phone != "" {
		phone = fmt.Sprintf("+39%s", phone)
	}
	rawDocumentCode, _ := strconv.Atoi(strings.TrimSpace(row[136+(offset*i)]))
	identityDocumentCode := fmt.Sprintf("%02d", rawDocumentCode)
	titolareEffettivo := models.User{
		Uid:           lib.NewDoc(models.UserCollection),
		Type:          models.UserLegalEntity,
		Name:          strings.TrimSpace(lib.Capitalize(row[118+(offset*i)])),
		Surname:       strings.TrimSpace(lib.Capitalize(row[117+(offset*i)])),
		FiscalCode:    strings.TrimSpace(strings.ToUpper(row[121+(offset*i)])),
		VatCode:       fmt.Sprintf("%011s", strings.TrimSpace(row[27])),
		Gender:        strings.TrimSpace(strings.ToUpper(row[119+(offset*i)])),
		BirthDate:     parseDate(row[120+(offset*i)]).Format(time.RFC3339),
		Mail:          strings.TrimSpace(strings.ToLower(row[131+(offset*i)])),
		Phone:         phone,
		BirthCity:     strings.TrimSpace(lib.Capitalize(row[133+(offset*i)])),
		BirthProvince: strings.TrimSpace(strings.ToUpper(row[134+(offset*i)])),
		Residence: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[122+(offset*i)])),
			City:       strings.TrimSpace(lib.Capitalize(row[124+(offset*i)])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[125+(offset*i)])),
			PostalCode: strings.TrimSpace(row[123+(offset*i)]),
			Locality:   strings.TrimSpace(lib.Capitalize(row[124+(offset*i)])),
		},
		Domicile: &models.Address{
			StreetName: strings.TrimSpace(lib.Capitalize(row[126+(offset*i)])),
			City:       strings.TrimSpace(lib.Capitalize(row[128+(offset*i)])),
			CityCode:   strings.TrimSpace(strings.ToUpper(row[129+(offset*i)])),
			PostalCode: strings.TrimSpace(row[127+(offset*i)]),
			Locality:   strings.TrimSpace(lib.Capitalize(row[128+(offset*i)])),
		},
		IdentityDocuments: []*models.IdentityDocument{{
			Number:           strings.TrimSpace(strings.ToUpper(row[137+(offset*i)])),
			Code:             identityDocumentCode,
			Type:             identityDocumentMap[identityDocumentCode],
			DateOfIssue:      parseDate(row[138+(offset*i)]),
			ExpiryDate:       parseDate(row[138+(offset*i)]).AddDate(10, 0, 0),
			IssuingAuthority: strings.TrimSpace(lib.Capitalize(row[139+(offset*i)])),
			PlaceOfIssue:     strings.TrimSpace(lib.Capitalize(row[139+(offset*i)])),
			LastUpdate:       now,
		}},
		Work:            strings.TrimSpace(lib.Capitalize(row[130+(offset*i)])),
		LegalEntityType: models.TitolareEffettivo,
		IsSignatory:     isExecutor,
		IsPayer:         isExecutor,
		CreationDate:    parseDate(row[4]),
		UpdatedDate:     time.Now().UTC(),
	}
	return titolareEffettivo
}

func parseDate(rawDate string) time.Time {
	day, _ := strconv.Atoi(strings.TrimSpace(rawDate[:2]))
	month, _ := strconv.Atoi(strings.TrimSpace(rawDate[2:4]))
	year, _ := strconv.Atoi(strings.TrimSpace(rawDate[4:]))

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return date
}

func saveUser(user *models.User) error {
	err := lib.SetFirestoreErr(lib.UserCollection, user.Uid, user)
	if err != nil {
		return err
	}

	err = user.BigquerySave("")
	return err
}
