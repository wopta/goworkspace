package companydata

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wopta/goworkspace/inclusive"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	dataMovement     = "inclusive_bank_account_movement"
	dataBanckAccount = "inclusive_bank_account"
	dateString       = "2021-11-22"
	layout           = "02/01/2006"
	layoutQuery      = "2006-01-02"
)

func BankAccountAxaInclusive(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from   time.Time
		to     time.Time
		result [][]string
	)
	log.Println("----------------BankAccountAxaInclusive-----------------")
	now := time.Now()
	M := time.Now().AddDate(0, 0, -2)
	date, _ := time.Parse("2006-01-02", dateString)
	log.Println(date)
	from, e = time.Parse("2006-01-02", strconv.Itoa(M.Year())+"-"+fmt.Sprintf("%02d", int(M.Month()))+"-"+fmt.Sprintf("%02d", 1))
	to, e = time.Parse("2006-01-02", strconv.Itoa(M.Year())+"-"+fmt.Sprintf("%02d", int(M.Month()))+"-"+fmt.Sprintf("%02d", M.Day()))
	log.Println(from)
	log.Println(to)
	query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME >'" + from.Format(layoutQuery) + " 00:00:00" + "' and _PARTITIONTIME <'" + to.Format(layoutQuery) + " 23:59:00" + "'"
	log.Println(query)
	bankaccountlist, e := lib.QueryRowsBigQuery[inclusive.BankAccountMovement](query)

	log.Println("len(bankaccountlist): ", len(bankaccountlist))
	//result = append(result, getHeader())
	result = append(result, getHeaderInclusiveBank())
	for _, mov := range bankaccountlist {

		result = append(result, setInclusiveRow(mov)...)

	}
	refMontly := now.AddDate(0, -1, 0)
	filepath := "180623_" + strconv.Itoa(refMontly.Year()) + fmt.Sprintf("%02d", int(refMontly.Month())) + "_" + fmt.Sprintf("%02d", now.Day())
	lib.CreateExcel(result, "../tmp/"+filepath+".xlsx", "")
	lib.WriteCsv("../tmp/"+filepath+".csv", result, ';')
	source, _ := ioutil.ReadFile("../tmp/" + filepath + ".xlsx")
	sourceCsv, _ := ioutil.ReadFile("../tmp/" + filepath + ".csv")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refMontly.Year())+"/"+fmt.Sprintf("%02d", int(refMontly.Month()))+"/"+filepath, source)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refMontly.Year())+"/"+fmt.Sprintf("%02d", int(refMontly.Month()))+"/"+filepath, sourceCsv)
	//AxaSftpUpload("/HYPE/IN/" + filepath + ".xlsx")
	return "", nil, e
}
func setInclusiveRow(mov inclusive.BankAccountMovement) [][]string {
	var (
		result [][]string
		user   models.User
	)
	log.Println(mov.FiscalCode)
	if mov.FiscalCode != "" {
		_, user, _ = ExtractUserDataFromFiscalCode(mov.FiscalCode)
	}
	birthDate, _ := time.Parse("1980-12-09T00:00:00Z", user.BirthDate)
	log.Println(mov.StartDate)
	row := []string{
		"180623",                     // NUMERO POLIZZA
		"T",                          //    LOB
		"C",                          //    TIPOLOGIA POLIZZA
		"0548100",                    //    CODICE CONFIGURAZIONE
		"1",                          //    TIPO OGGETTO ASSICURATO
		mov.HypeId,                   //    IDENTIFICATIVO UNIVOCO APPLICAZIONE
		mov.FiscalCode,               //    CODICE FISCALE / P.IVA ASSICURATO
		mov.Surname,                  //    COGNOME / RAGIONE SOCIALE ASSICURATO
		mov.Name,                     //    NOME ASSICURATO
		"",                           //    INDIRIZZO RESIDENZA ASSICURATO
		"",                           //    CAP RESIDENZA ASSICURATO
		"",                           //    CITTA' RESIDENZA ASSICURATO
		"",                           //    PROVINCIA RESIDENZA ASSICURATO
		birthDate.Format(layout),     //    DATA DI NASCITA ASSICURATO 1980-12-09T00:00:00Z
		mov.StartDate.Format(layout), //    DATA INIZIO VALIDITA' COPERTURA
		mapEndDate(mov),              //    DATA FINE VALIDITA' COPERTURA
		StringMapping(mov.MovementType, map[string]string{
			"insert": "A",
			"delete": "E",
		}), //    TIPO MOVIMENTO
	}

	result = append(result, row)

	return result
}
func mapEndDate(mov inclusive.BankAccountMovement) string {
	if mov.MovementType == "delete" {
		return mov.EndDate.Format(layout)
	}
	return ""
}
func getHeaderInclusiveBank() []string {
	return []string{
		"NUMERO POLIZZA",                       // NUMERO POLIZZA
		"LOB",                                  //    LOB
		"TIPOLOGIA POLIZZA",                    //    TIPOLOGIA POLIZZA
		"CODICE CONFIGURAZIONE",                //    CODICE CONFIGURAZIONE
		"TIPO OGGETTO ASSICURATO",              //    TIPO OGGETTO ASSICURATO
		"IDENTIFICATIVO UNIVOCO APPLICAZIONE",  //    IDENTIFICATIVO UNIVOCO APPLICAZIONE
		"CODICE FISCALE / P.IVA ASSICURATO",    //    CODICE FISCALE / P.IVA ASSICURATO
		"COGNOME / RAGIONE SOCIALE ASSICURATO", //    COGNOME / RAGIONE SOCIALE ASSICURATO
		"NOME ASSICURATO",                      //    NOME ASSICURATO
		"INDIRIZZO RESIDENZA ASSICURATO",       //    INDIRIZZO RESIDENZA ASSICURATO
		"CAP RESIDENZA ASSICURATO",             //    CAP RESIDENZA ASSICURATO
		" CITTA' RESIDENZA ASSICURATO",         //    CITTA' RESIDENZA ASSICURATO
		"PROVINCIA RESIDENZA ASSICURATO",       //    PROVINCIA RESIDENZA ASSICURATO
		"DATA DI NASCITA ASSICURATO",           //    DATA DI NASCITA ASSICURATO
		"DATA INIZIO VALIDITA' COPERTURA",      //    DATA INIZIO VALIDITA' COPERTURA
		" DATA FINE VALIDITA' COPERTURA",       //    DATA FINE VALIDITA' COPERTURA
		"TIPO MOVIMENTO",
	}
}
