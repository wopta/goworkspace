package _script

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

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

	for _, rawPolicy := range groups {
		log.Printf("%v", rawPolicy)

		// TODO: implementare estrazione dati contraente persona giuridica
		contractor := extractContractorData(rawPolicy[0])
		// TODO: implementare estrazione dati 3 titolari effettivi
		// TODO: implementare estrazione dati assicurato

		log.Printf("contractor: %+v", contractor)
	}

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
