package companydata

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/quote"
)

const (
	gapProduct           = "gap"
	gapCompany           = "sogessur"
	gapDateFormat        = "02/01/2006"
	gapCsvFilenameFormat = "Contratti_GAP_%02d_%04d.csv"
	storagePath          = "track/" + gapCompany + "/" + gapProduct + "/"
	tmpPath              = "../tmp/"
)

func GapSogessurEmit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	prevMonth := getPreviousMonth()
	from := getFirstDay(prevMonth)
	to := getFirstDay(time.Now())

	filename := fmt.Sprintf(gapCsvFilenameFormat, prevMonth.Month(), prevMonth.Year())

	policies := getGapPolicies(from, to)
	if len(policies) == 0 {
		return "", nil, fmt.Errorf("no policy found")
	}
	transactions := getGapTransactions(policies)
	if len(policies) != len(transactions) {
		return "", nil, fmt.Errorf("number of transactions doesn't match number of policies")
	}
	csvRows := getGapCsv(policies, transactions)
	lib.WriteCsv(tmpPath+filename, csvRows, ';')
	source, err := os.ReadFile(tmpPath + filename)
	if err != nil {
		panic(err)
	}

	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), storagePath+filename, source)
	// TODO: SftUpload

	//setCompanyEmitted(policies)

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

func getGapRowMap(policy models.Policy, transaction models.Transaction) map[string]string {
	vehicle := policy.Assets[0].Vehicle
	contractor := policy.Contractor
	vehicleOwner := policy.Assets[0].Person
	offerName := policy.OfferlName

	genders := []string{"F", "M", "G"} // For validation
	// Maps for replacing the values in policy into valid values for the gap csv
	boolAnswers := map[bool]string{
		true:  "SI",
		false: "NO",
	}
	gapOptions := map[string]string{
		"base":     "OPTION 1",
		"complete": "OPTION 2",
	}
	isVehicleNew := map[string]string{
		"Nuovo": "SI",
		"Usato": "NO",
	}
	powerSupplyCodes := map[bool]string{
		true:  "E",
		false: "",
	}
	vehicleTypeCodes := map[string]string{
		"car":    "A",
		"truck":  "C",
		"camper": "P",
	}
	// Assuming we have only this payment type
	offer := policy.OffersPrices[offerName][string(models.PaySingleInstallment)]
	vehicleWeight := ""
	if vehicle.Weight > 0 {
		vehicleWeight = strconv.Itoa(int(vehicle.Weight))
	}
	zoneCode := quote.GetAreaByProvince(vehicleOwner.Residence.CityCode)
	return map[string]string{
		"NUMERO POLIZZA":                       CheckIfIsAlphaNumeric(policy.CodeCompany),
		"NUMERO CONTRATTO":                     CheckIfIsAlphaNumeric(policy.CodeCompany),
		"TIPO OPERAZIONE":                      "A", // 'A' = Subscription
		"DATA OPERAZIONE":                      policy.StartDate.Format(gapDateFormat),
		"COGNOME/RAGIONE SOCIALE ASSICURATO":   vehicleOwner.Surname,
		"NOME ASSICURATO":                      vehicleOwner.Name,
		"INDIRIZZO ASSICURATO":                 getAddress(*vehicleOwner.Residence),
		"COMUNE ASSICURATO":                    vehicleOwner.Residence.Locality,
		"CAP ASSICURATO":                       CheckIfIsNumeric(vehicleOwner.Residence.PostalCode),
		"PROVINCIA ASSICURATO":                 CheckIfIsAlphaNumeric(vehicleOwner.Residence.CityCode),
		"NAZIONE ASSICURATO":                   "Italia",
		"CODICE FISCALE ASSICURATO":            CheckIfIsAlphaNumeric(vehicleOwner.FiscalCode),
		"PARTITA IVA ASSICURATO":               CheckIfIsNumeric(vehicleOwner.VatCode),
		"TARGA":                                CheckIfIsAlphaNumeric(vehicle.Plate),
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
		"DURATA COPERTURA":                     strconv.Itoa(ElapsedMonths(policy.StartDate, policy.EndDate)),
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
		"DATA INCASSO":                         policy.StartDate.Format(gapDateFormat),
		"VINCOLO":                              "",
		"DATA VINCOLO":                         "",
		"ENTE VINCOLATARIO":                    "",
		"CODICE ZONA":                          CheckIfIsAlphaNumeric(zoneCode),
		"CL_SESSO":                             CheckIfIsWithin(contractor.Gender, genders),
		"CL_DATA_NASC":                         stringToDateFormat(contractor.BirthDate, gapDateFormat),
		"CL_LUOGO_NASC":                        contractor.BirthCity,
		"CL_PROV_NASC":                         CheckIfIsAlphaNumeric(contractor.BirthProvince),
		"COGNOME/RAGIONE SOCIALE PROPRIETARIO": vehicleOwner.Surname,
		"NOME PROPRIETARIO":                    vehicleOwner.Name,
		"INDIRIZZO PROPRIETARIO":               getAddress(*vehicleOwner.Residence),
		"COMUNE PROPRIETARIO":                  vehicleOwner.Residence.Locality,
		"CAP PROPRIETARIO":                     CheckIfIsNumeric(vehicleOwner.Residence.PostalCode),
		"PROVINCIA PROPRIETARIO":               CheckIfIsAlphaNumeric(vehicleOwner.Residence.CityCode),
		"NAZIONE PROPRIETARIO":                 "Italia",
		"CODICE FISCALE PROPRIETARIO":          CheckIfIsAlphaNumeric(vehicleOwner.FiscalCode),
		"PARTITA IVA PROPRIETARIO":             "",
		"PR_SESSO":                             CheckIfIsWithin(vehicleOwner.Gender, genders),
		"PR_DATA_NASC":                         stringToDateFormat(vehicleOwner.BirthDate, gapDateFormat),
		"PR_LUOGO_NASC":                        vehicleOwner.BirthCity,
		"PR_PROV_NASC":                         CheckIfIsAlphaNumeric(vehicleOwner.BirthProvince),
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
				QueryValue: gapCompany,
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
				Operator:   ">",
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

func getPreviousMonth() time.Time {
	return time.Now().AddDate(0, -1, 0)
}

func getFirstDay(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func getLastDay(t time.Time) time.Time {
	t = getFirstDay(t)
	year, month, _ := t.Date()
	lastDay := t.AddDate(0, 1, -1).Day()
	return time.Date(year, month, lastDay, 0, 0, 0, 0, time.UTC)
}

func ElapsedMonths(t1 time.Time, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}

	t1y, t1m, t1d := t1.Date()
	date1 := time.Date(t1y, t1m, t1d, 0, 0, 0, 0, time.UTC)

	t2y, t2m, t2d := t2.Date()

	months := (t2y - t1y) * 12
	anniversary := date1.AddDate(0, months, 0)
	months += int(t2m - anniversary.Month())

	if t2d < t1d {
		months--
	}

	return months
}

// ----------------------------------------------------
// ----------------VALIDATION--------------------------
// ----------------------------------------------------

func CheckIfIsWithin(value string, values []string) string {
	if value == "" {
		return value
	}
	if !lib.SliceContains(values, value) {
		panic(errors.New("value not in slice"))
	}
	return value
}

func CheckIfIsDate(value string) string {
	if value == "" {
		return value
	}
	RegexPanicOnFail(value, `^\d{2}/\d{2}/\d{4}$`, "the value is not a date")
	return value
}

func CheckIfIsAlphaNumeric(value string) string {
	RegexPanicOnFail(value, "^[A-Za-z0-9]*$", "the value is not alphanumeric")
	return value
}

func CheckIfIsNumeric(value string) string {
	RegexPanicOnFail(value, `^\d*$`, "the value is not an integer")
	return value
}

func RegexPanicOnFail(value string, pattern string, noMatchMsg string) {
	isMatching, err := regexp.Match(pattern, []byte(value))
	if err != nil {
		panic(err)
	}
	if !isMatching {
		panic(errors.New(noMatchMsg + ": " + value))
	}
}
