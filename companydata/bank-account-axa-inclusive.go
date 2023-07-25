package companydata

import (
	"encoding/json"
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

type BankAccountAxaInclusiveReq struct {
	Day string `firestore:"-" json:"day,omitempty" bigquery:"-"`
}

func BankAccountAxaInclusive(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from      time.Time
		to        time.Time
		result    [][]string
		now       time.Time
		refMontly time.Time
	)
	log.Println("----------------BankAccountAxaInclusive-----------------")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(req))
	var obj BankAccountAxaInclusiveReq
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	if obj.Day == "" {
		refMontly = now.AddDate(0, 0, -1)
		now = time.Now().AddDate(0, 0, -1)

		from, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 1))
		to, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", now.Day()))
	} else {
		date, _ := time.Parse("2006-01-02", obj.Day)
		refMontly = date
		log.Println(date)
		from = date
		to = date
	}
	log.Println(from)
	log.Println(to)
	//query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME >'" + from.Format(layoutQuery) + " 00:00:00" + "' and _PARTITIONTIME <'" + to.Format(layoutQuery) + " 23:59:00" + "'"
	query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME ='" + from.Format(layoutQuery) + "'"
	log.Println(query)
	bankaccountlist, e := lib.QueryRowsBigQuery[inclusive.BankAccountMovement](query)
	log.Println("len(bankaccountlist): ", len(bankaccountlist))
	//result = append(result, getHeader())
	result = append(result, getHeaderInclusiveBank())
	for _, mov := range bankaccountlist {

		result = append(result, setInclusiveRow(mov)...)

	}

	filepath := "180623_" + strconv.Itoa(refMontly.Year()) + fmt.Sprintf("%02d", int(refMontly.Month())) + fmt.Sprintf("%02d", refMontly.Day())
	CreateExcel(result, "../tmp/"+filepath+".xlsx")
	lib.WriteCsv("../tmp/"+filepath+".csv", result, ';')
	source, _ := ioutil.ReadFile("../tmp/" + filepath + ".xlsx")
	sourceCsv, _ := ioutil.ReadFile("../tmp/" + filepath + ".csv")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refMontly.Year())+"/"+fmt.Sprintf("%02d", int(refMontly.Month()))+"/"+filepath+".xlsx", source)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refMontly.Year())+"/"+fmt.Sprintf("%02d", int(refMontly.Month()))+"/"+filepath+".csv", sourceCsv)
	AxaSftpUpload(filepath+".xlsx", "HYPE/IN/")
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
	birthDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", user.BirthDate)
	startDate, _ := time.Parse("2006-01-02", mov.BigStartDate.Date.String())
	log.Println(mov.StartDate)
	row := []string{
		"180623",                 // NUMERO POLIZZA
		"T",                      //    LOB
		"C",                      //    TIPOLOGIA POLIZZA
		"0548100",                //    CODICE CONFIGURAZIONE
		"1",                      //    TIPO OGGETTO ASSICURATO
		mov.HypeId,               //    IDENTIFICATIVO UNIVOCO APPLICAZIONE
		mov.FiscalCode,           //    CODICE FISCALE / P.IVA ASSICURATO
		mov.Surname,              //    COGNOME / RAGIONE SOCIALE ASSICURATO
		mov.Name,                 //    NOME ASSICURATO
		"",                       //    INDIRIZZO RESIDENZA ASSICURATO
		"",                       //    CAP RESIDENZA ASSICURATO
		"",                       //    CITTA' RESIDENZA ASSICURATO
		"",                       //    PROVINCIA RESIDENZA ASSICURATO
		birthDate.Format(layout), //    DATA DI NASCITA ASSICURATO 1980-12-09T00:00:00Z
		startDate.Format(layout), //    DATA INIZIO VALIDITA' COPERTURA
		mapEndDate(mov),          //    DATA FINE VALIDITA' COPERTURA
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
		endDate, _ := time.Parse("2006-01-02", mov.BigEndDate.Date.String())
		return endDate.Format(layout)
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
