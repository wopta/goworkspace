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
		result [][]string
		now    time.Time
		refDay time.Time
		upload bool
	)
	log.Println("----------------BankAccountAxaInclusive-----------------")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(req))
	var obj BankAccountAxaInclusiveReq
	defer r.Body.Close()
	json.Unmarshal([]byte(req), &obj)
	now, upload = getRequestData(req)
	refDay = now.AddDate(0, 0, -1)

	log.Println("BankAccountAxaInclusive refMontly: ", refDay)
	//from, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 1))
	//query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME >'" + from.Format(layoutQuery) + " 00:00:00" + "' and _PARTITIONTIME <'" + to.Format(layoutQuery) + " 23:59:00" + "'"
	query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME ='" + refDay.Format(layoutQuery) + "'"
	log.Println("BankAccountAxaInclusive bigquery query: ", query)
	bankaccountlist, e := lib.QueryRowsBigQuery[inclusive.BankAccountMovement](query)
	log.Println("BankAccountAxaInclusive bigquery error: ", e)
	log.Println("BankAccountAxaInclusive len(bankaccountlist): ", len(bankaccountlist))
	//result = append(result, getHeader())
	result = append(result, getHeaderInclusiveBank())

	b, err := os.ReadFile(lib.GetAssetPathByEnv("companyData") + "/reverse-codes.json")
	var codes map[string]map[string]string
	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)
	lib.CheckError(err)
	for i, mov := range bankaccountlist {
		log.Println(i)
		result = append(result, setInclusiveRow(mov, codes)...)

	}

	filepath := "180623_" + strconv.Itoa(refDay.Year()) + fmt.Sprintf("%02d", int(refDay.Month())) + fmt.Sprintf("%02d", refDay.Day())
	CreateExcel(result, "../tmp/"+filepath+".xlsx")
	lib.WriteCsv("../tmp/"+filepath+".csv", result, ';')
	source, _ := ioutil.ReadFile("../tmp/" + filepath + ".xlsx")
	sourceCsv, _ := ioutil.ReadFile("../tmp/" + filepath + ".csv")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".xlsx", source)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".csv", sourceCsv)
	if upload {

		AxaSftpUpload(filepath+".xlsx", "HYPE/IN/")
	}
	log.Println("---------------------end------------------------------")
	return "", nil, e
}
func setInclusiveRow(mov inclusive.BankAccountMovement, codes map[string]map[string]string) [][]string {
	var (
		result [][]string
		user   models.User
	)

	if mov.FiscalCode != "" {
		_, user, _ = ExtractUserDataFromFiscalCode(mov.FiscalCode, codes)
	}
	birthDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", user.BirthDate)
	startDate, _ := time.Parse("2006-01-02", mov.BigStartDate.Date.String())
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
			"":          "A",
			"active":    "A",
			"insert":    "A",
			"delete":    "E",
			"suspended": "E",
		}), //    TIPO MOVIMENTO
	}
	layout := "2006-01-02"
	if mov.MovementType == "delete" || mov.MovementType == "suspended" {
		e = lib.UpdateRowBigQuery("wopta", dataBanckAccount, map[string]string{
			"status":  mov.Status,
			"endDate": mov.EndDate.Format(layout) + " 00:00:00",
		}, "fiscalCode='"+mov.FiscalCode+"' and guaranteesCode='"+mov.GuaranteesCode+"'")

	}

	result = append(result, row)

	return result
}
func mapEndDate(mov inclusive.BankAccountMovement) string {
	if mov.MovementType == "delete" {
		endDate, _ := time.Parse("2006-01-02", mov.BigEndDate.Date.String())
		return endDate.Format(layout)
	}
	return "31/12/9999"
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
		"DATA FINE VALIDITA' COPERTURA",        //    DATA FINE VALIDITA' COPERTURA
		"TIPO MOVIMENTO",
	}
}
