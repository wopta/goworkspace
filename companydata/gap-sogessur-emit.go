package companydata

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	gapDateFormat                  = "02/01/2006"
	gapCsvFilenameFormat           = "Contratti_GAP_%02d_%04d.csv"
	storagePath                    = "track/" + models.SogessurCompany + "/" + models.GapProduct + "/"
	tmpPath                        = "../tmp/"
	subscriptionOperation          = "A"
	withdrawalOperation            = "P"
	earlyWithdrawalOperation       = "S"
	variationWithPriceOperation    = "I"
	variationWithoutPriceOperation = "V"
)

var (
	// Maps for replacing the values in policy into valid values for the gap csv
	boolAnswers = map[bool]string{
		true:  "SI",
		false: "NO",
	}
	gapOptions = map[string]string{
		"base":     "OPTION 1",
		"complete": "OPTION 2",
	}
	isVehicleNew = map[string]string{
		"Nuovo": "SI",
		"Usato": "NO",
	}
	powerSupplyCodes = map[bool]string{
		true:  "E",
		false: "",
	}
	vehicleTypeCodes = map[string]string{
		"car":    "A",
		"truck":  "C",
		"camper": "P",
	}
)

func GapSogessurEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	now := time.Now()
	prevMonth := lib.GetPreviousMonth(now)
	from := lib.GetFirstDay(prevMonth)
	to := lib.GetFirstDay(now)

	filename := fmt.Sprintf(gapCsvFilenameFormat, prevMonth.Month(), prevMonth.Year())

	policies := getGapPolicies(from, to)
	if len(policies) == 0 {
		return "", nil, fmt.Errorf("[GapSogessurEmit] no policies found")
	}
	transactions := getGapTransactions(policies)
	if len(policies) != len(transactions) {
		return "", nil, fmt.Errorf("[GapSogessurEmit] number of transactions doesn't match number of policies")
	}
	csvRows := getGapCsv(policies, transactions)
	lib.WriteCsv(tmpPath+filename, csvRows, ';')
	source, err := os.ReadFile(tmpPath + filename)
	if err != nil {
		panic(err)
	}

	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), storagePath+filename, source)
	// TODO: SftUpload

	setCompanyEmitted(policies)

	return "", nil, e
}

func getGapCsv(policies []models.Policy, transactions []models.Transaction) [][]string {
	header := getGapHeader()
	// Space for header
	csvRows := make([][]string, len(policies)+1)
	csvRows[0] = header

	// Caching columns' position in array for faster search
	// given the column's name it returns its position in matrix
	columnsIdx := make(map[string]int)
	for i, column := range csvRows[0] {
		columnsIdx[column] = i
	}

	for i := range policies {
		rowMap := getGapRowMap(policies[i], transactions[i])
		if len(rowMap) > len(header) {
			panic(fmt.Errorf("the header (%d) has less columns than the row's fields (%d)", len(header), len(rowMap)))
		}

		if row, err := mapToSlice(rowMap, columnsIdx); err != nil {
			panic(err)
		} else {
			csvRows[i+1] = row
		}
	}
	return csvRows
}

func mapToSlice(values map[string]string, indices map[string]int) ([]string, error) {
	slice := make([]string, len(values))
	for key, value := range values {
		index, ok := indices[key]
		if !ok {
			return nil, fmt.Errorf("the key %q is not present in the header", key)
		}
		slice[index] = value
	}
	return slice, nil
}

func getGapHeader() []string {
	return []string{
		"NUMERO POLIZZA",
		"NUMERO CONTRATTO",
		"TIPO OPERAZIONE",
		"DATA OPERAZIONE",
		"COGNOME/RAGIONE SOCIALE ASSICURATO",
		"NOME ASSICURATO",
		"INDIRIZZO ASSICURATO",
		"COMUNE ASSICURATO",
		"CAP ASSICURATO",
		"PROVINCIA ASSICURATO",
		"NAZIONE ASSICURATO",
		"CODICE FISCALE ASSICURATO",
		"PARTITA IVA ASSICURATO",
		"TARGA",
		"TELAIO",
		"MODELLO",
		"MARCA",
		"CILINDRATA",
		"ALIMENTAZIONE",
		"ANTIFURTO SATELLITARE",
		"TIPO VEICOLO",
		"TIPO TARGA",
		"CAVALLI",
		"QUINTALI",
		"KW",
		"CODICE MODELLO",
		"VEICOLO NUOVO",
		"DATA CONSEGNA",
		"DATA IMMATRICOLAZIONE",
		"CODICE DEALER",
		"DEALER",
		"PACCHETTO",
		"CAPITALE ASSICURATO",
		"DATA INIZIO POLIZZA",
		"DATA FINE POLIZZA",
		"ORA EFFETTO",
		"ORA SCADENZA",
		"DATA SCADENZA RATE",
		"TIPO FRAZIONAMENTO",
		"NUMERO RATE",
		"DURATA COPERTURA",
		"PREMIO NETTO GAP",
		"IMPOSTE GAP",
		"PREMIO LORDO GAP",
		"PREMIO NETTO CVT",
		"IMPOSTE CVT (comprensivo dell'1%)",
		"PREMIO LORDO CVT",
		"PREMIO NETTO PAI",
		"IMPOSTE PAI",
		"PREMIO LORDO PAI",
		"PREMIO NETTO TL",
		"IMPOSTE TL",
		"PREMIO LORDO TL",
		"PROVVIGIONI GAP",
		"PROVVIGIONI PAI",
		"PROVVIGIONI TL",
		"PROVVIGIONE BK",
		"PROVVIGIONE BK PAI",
		"PROVVIGIONE BK TL",
		"FEE MGT",
		"TIPO PACCHETTO",
		"IBAN",
		"DATA INCASSO",
		"VINCOLO",
		"DATA VINCOLO",
		"ENTE VINCOLATARIO",
		"CODICE ZONA",
		"CL_SESSO",
		"CL_DATA_NASC",
		"CL_LUOGO_NASC",
		"CL_PROV_NASC",
		"COGNOME/RAGIONE SOCIALE PROPRIETARIO",
		"NOME PROPRIETARIO",
		"INDIRIZZO PROPRIETARIO",
		"COMUNE PROPRIETARIO",
		"CAP PROPRIETARIO",
		"PROVINCIA PROPRIETARIO",
		"NAZIONE PROPRIETARIO",
		"CODICE FISCALE PROPRIETARIO",
		"PARTITA IVA PROPRIETARIO",
		"PR_SESSO",
		"PR_DATA_NASC",
		"PR_LUOGO_NASC",
		"PR_PROV_NASC",
	}
}

func getOperationType(policy models.Policy) string {
	if policy.CompanyEmit && !policy.IsDeleted {
		return subscriptionOperation
	} else if policy.IsDeleted && policy.DeleteEmited {
		return withdrawalOperation
	}
	return ""
}

func getGapRowMap(policy models.Policy, transaction models.Transaction) map[string]string {
	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person
	offerName := policy.OfferlName

	// Assuming we have only this payment type
	offer := policy.OffersPrices[offerName][string(models.PaySingleInstallment)]
	vehicleWeight := ""
	if vehicle.Weight > 0 {
		vehicleWeight = strconv.Itoa(int(vehicle.Weight))
	}
	return map[string]string{
		"NUMERO POLIZZA":                       policy.CodeCompany,
		"NUMERO CONTRATTO":                     policy.CodeCompany,
		"TIPO OPERAZIONE":                      getOperationType(policy),
		"DATA OPERAZIONE":                      policy.StartDate.Format(gapDateFormat),
		"COGNOME/RAGIONE SOCIALE ASSICURATO":   contractor.Surname,
		"NOME ASSICURATO":                      contractor.Name,
		"INDIRIZZO ASSICURATO":                 getAddress(*contractor.Residence),
		"COMUNE ASSICURATO":                    contractor.Residence.Locality,
		"CAP ASSICURATO":                       contractor.Residence.PostalCode,
		"PROVINCIA ASSICURATO":                 contractor.Residence.CityCode,
		"NAZIONE ASSICURATO":                   "Italia",
		"CODICE FISCALE ASSICURATO":            contractor.FiscalCode,
		"PARTITA IVA ASSICURATO":               contractor.VatCode,
		"TARGA":                                vehicle.Plate,
		"TELAIO":                               "",
		"MODELLO":                              vehicle.Model,
		"MARCA":                                vehicle.Manufacturer,
		"CILINDRATA":                           "",
		"ALIMENTAZIONE":                        powerSupplyCodes[vehicle.IsElectric],
		"ANTIFURTO SATELLITARE":                boolAnswers[vehicle.HasSatellite],
		"TIPO VEICOLO":                         vehicleTypeCodes[vehicle.VehicleType],
		"TIPO TARGA":                           "",
		"CAVALLI":                              "",
		"QUINTALI":                             vehicleWeight,
		"KW":                                   "",
		"CODICE MODELLO":                       "",
		"VEICOLO NUOVO":                        isVehicleNew[vehicle.Condition],
		"DATA CONSEGNA":                        "",
		"DATA IMMATRICOLAZIONE":                vehicle.RegistrationDate.Format(gapDateFormat),
		"PACCHETTO":                            gapOptions[offerName],
		"CODICE DEALER":                        "",
		"DEALER":                               "",
		"CAPITALE ASSICURATO":                  floatToPrice(vehicle.PriceValue),
		"DATA INIZIO POLIZZA":                  policy.StartDate.Format(gapDateFormat),
		"DATA FINE POLIZZA":                    policy.EndDate.Format(gapDateFormat),
		"ORA EFFETTO":                          "24:00",
		"ORA SCADENZA":                         "24:00",
		"DATA SCADENZA RATE":                   "",
		"TIPO FRAZIONAMENTO":                   "",
		"NUMERO RATE":                          "",
		"DURATA COPERTURA":                     strconv.Itoa(lib.MonthsDifference(policy.StartDate, policy.EndDate)),
		"PREMIO NETTO GAP":                     floatToPrice(offer.Net),
		"IMPOSTE GAP":                          floatToPrice(offer.Tax),
		"PREMIO LORDO GAP":                     floatToPrice(offer.Gross),
		"PREMIO NETTO CVT":                     "",
		"IMPOSTE CVT (comprensivo dell'1%)":    "",
		"PREMIO LORDO CVT":                     "",
		"PREMIO NETTO PAI":                     "",
		"IMPOSTE PAI":                          "",
		"PREMIO LORDO PAI":                     "",
		"PREMIO NETTO TL":                      "",
		"IMPOSTE TL":                           "",
		"PREMIO LORDO TL":                      "",
		"PROVVIGIONI GAP":                      floatToPrice(transaction.Commissions),
		"PROVVIGIONI PAI":                      "",
		"PROVVIGIONI TL":                       "",
		"PROVVIGIONE BK":                       "",
		"PROVVIGIONE BK TL":                    "",
		"PROVVIGIONE BK PAI":                   "",
		"FEE MGT":                              "",
		"TIPO PACCHETTO":                       "",
		"IBAN":                                 "",
		"DATA INCASSO":                         transaction.PayDate.Format(gapDateFormat),
		"VINCOLO":                              "",
		"DATA VINCOLO":                         "",
		"ENTE VINCOLATARIO":                    "",
		"CODICE ZONA":                          vehicleOwner.Residence.Area,
		"CL_SESSO":                             contractor.Gender,
		"CL_DATA_NASC":                         stringToDateFormat(contractor.BirthDate, gapDateFormat),
		"CL_LUOGO_NASC":                        contractor.BirthCity,
		"CL_PROV_NASC":                         contractor.BirthProvince,
		"COGNOME/RAGIONE SOCIALE PROPRIETARIO": vehicleOwner.Surname,
		"NOME PROPRIETARIO":                    vehicleOwner.Name,
		"INDIRIZZO PROPRIETARIO":               getAddress(*vehicleOwner.Residence),
		"COMUNE PROPRIETARIO":                  vehicleOwner.Residence.Locality,
		"CAP PROPRIETARIO":                     vehicleOwner.Residence.PostalCode,
		"PROVINCIA PROPRIETARIO":               vehicleOwner.Residence.CityCode,
		"NAZIONE PROPRIETARIO":                 "Italia",
		"CODICE FISCALE PROPRIETARIO":          vehicleOwner.FiscalCode,
		"PARTITA IVA PROPRIETARIO":             "",
		"PR_SESSO":                             vehicleOwner.Gender,
		"PR_DATA_NASC":                         stringToDateFormat(vehicleOwner.BirthDate, gapDateFormat),
		"PR_LUOGO_NASC":                        vehicleOwner.BirthCity,
		"PR_PROV_NASC":                         vehicleOwner.BirthProvince,
	}
}

func getAddress(address models.Address) string {
	return fmt.Sprintf("%s %s %s %s", address.StreetName, address.StreetNumber, address.PostalCode, address.CityCode)
}

// The float numbers in the csv are divided by ',' and not by '.'
func floatToPrice(n float64) string {
	return humanize.FormatFloat("####,##", n)
}

// Some dates are saved as strings(time.RFC3339 format), this ensures that they will follow the csv date format.
func stringToDateFormat(date string, layout string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	return t.Format(layout)
}

func getGapPolicies(from time.Time, to time.Time) []models.Policy {
	queries := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "company",
				Operator:   "==",
				QueryValue: models.SogessurCompany,
			},
			{
				Field:      "companyEmit",
				Operator:   "==",
				QueryValue: true, // NOTE: Not working for testing env: the DB has every Gap policy with companyEmit to false
			},
			{
				Field:      "companyEmitted",
				Operator:   "==",
				QueryValue: false,
			},
			{
				Field:      "name",
				Operator:   "==",
				QueryValue: models.GapProduct,
			},
			{
				Field:      "startDate",
				Operator:   ">=",
				QueryValue: from,
			},
			{
				Field:      "startDate",
				Operator:   "<",
				QueryValue: to,
			},
		},
	}

	iter, err := queries.FirestoreWherefields(models.PolicyCollection)
	if err != nil {
		log.Println(err.Error())
	}

	policies := models.PolicyToListData(iter)
	return policies
}

// Returns all the transactions with the same order of their relative policies.
// That is, the transaction of policy[i] has the same index: policy[i] => transaction[i]
func getGapTransactions(policies []models.Policy) []models.Transaction {
	transactions := make([]models.Transaction, len(policies))
	for i, policy := range policies {
		iter := lib.WhereFirestore(models.TransactionsCollection, "policyUid", "==", policy.Uid)
		transactionsBuffer := models.TransactionToListData(iter)
		if len(transactionsBuffer) == 0 {
			panic(errors.New("no transcations found"))
		}
		transactions[i] = transactionsBuffer[0]
	}
	return transactions
}

func setCompanyEmitted(policies []models.Policy) {
	for _, policy := range policies {
		policy.CompanyEmitted = true
		policy.Updated = time.Now().UTC()
		lib.SetFirestore(models.PolicyCollection, policy.Uid, policy)
		policy.BigquerySave("")
	}
}
