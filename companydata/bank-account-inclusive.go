package companydata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/bigquery"
	"gitlab.dev.wopta.it/goworkspace/inclusive"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/iterator"
)

const (
	dataMovement       = "inclusive_bank_account_movement"
	dataBanckAccount   = "inclusive_bank_account"
	dateString         = "2021-11-22"
	layout             = "02/01/2006"
	layoutQuery        = "2006-01-02"
	dataset            = "wopta_inclusive"
	movementTable      = "bank_account_movement"
	usersTableScalapay = "bank_account_users_scalapay"
)

type BankAccountAxaInclusiveReq struct {
	Day string `firestore:"-" json:"day,omitempty" bigquery:"-"`
}

func BankAccountInclusive(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	log.AddPrefix("BankAccountInclusive ")
	defer log.PopPrefix()
	body := lib.ErrorByte(io.ReadAll(r.Body))
	log.Println(r.Header)
	log.Println(string(body))

	defer r.Body.Close()
	now, upload, req := getCompanyDataReq(body)

	setPolicy("180623", now, upload, req)
	setPolicy("191123", now, upload, req)
	setScalapayPolicy("51114", now, upload, req)
	log.Println("---------------------end------------------------------")
	return "", nil, e
}

func setPolicy(code string, now time.Time, upload bool, req DataReq) {
	var (
		result [][]string
		query  string
		refDay time.Time
	)

	log.AddPrefix("BankAccountInclusive setPolicy " + code)
	defer log.PopPrefix()
	refDay = now.AddDate(0, 0, -1)
	log.Println("  refMontly: ", refDay)
	//from, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 1))
	//query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME >'" + from.Format(layoutQuery) + " 00:00:00" + "' and _PARTITIONTIME <'" + to.Format(layoutQuery) + " 23:59:00" + "'"
	if req.From != "" && req.To != "" {
		from, _ := time.Parse(layoutQuery, req.From)
		to, _ := time.Parse(layoutQuery, req.To)
		query = "select * from `wopta." + dataMovement + "` where startDate>'" + from.Format(layoutQuery) + "' and startDate<'" + to.Format(layoutQuery) + "' and policyNumber='" + code + "'"
	} else {
		query = "select * from `wopta." + dataMovement + "` where _PARTITIONTIME ='" + refDay.Format(layoutQuery) + "' and policyNumber='" + code + "'"
	}
	log.Println("  bigquery query: ", query)
	bankaccountlist, e := QueryRowsBigQuery[inclusive.BankAccountMovement](query)
	log.Println("  bigquery error: ", e)
	log.Println("  len(bankaccountlist): ", len(bankaccountlist))
	//result = append(result, getHeader())
	result = append(result, getHeaderInclusiveBank())
	b, err := os.ReadFile(lib.GetAssetPathByEnv("companyData") + "/reverse-codes.json")
	lib.CheckError(err)
	var codes map[string]map[string]string
	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)
	for i, mov := range bankaccountlist {
		log.Println(i)
		result = append(result, setInclusiveRow(mov, codes)...)

	}

	filepath := code + "_" + strconv.Itoa(refDay.Year()) + fmt.Sprintf("%02d", int(refDay.Month())) + fmt.Sprintf("%02d", refDay.Day())
	//CreateExcel(result, "../tmp/"+filepath+".xlsx", "Sheet1")
	lib.WriteCsv("../tmp/"+filepath+".csv", result, ';')
	source, _ := os.ReadFile("../tmp/" + filepath + ".xlsx")
	sourceCsv, _ := os.ReadFile("../tmp/" + filepath + ".csv")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".xlsx", source)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/hype/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".csv", sourceCsv)
	if upload {

		AxaSftpUpload(filepath+".xlsx", "HYPE/IN/")
	}

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
		mov.PolicyNumber,         // NUMERO POLIZZA
		"T",                      //    LOB
		"C",                      //    TIPOLOGIA POLIZZA
		"0548100",                //    CODICE CONFIGURAZIONE
		"1",                      //    TIPO OGGETTO ASSICURATO
		mov.HypeId,               //    IDENTIFICATIVO UNIVOCO APPLICAZIONE
		mov.FiscalCode,           //    CODICE FISCALE / P.IVA ASSICURATO
		mov.Surname,              //    COGNOME / RAGIONE SOCIALE ASSICURATO
		mov.Name,                 //    NOME ASSICURATO
		mov.Address,              //    INDIRIZZO RESIDENZA ASSICURATO
		mov.PostalCode,           //    CAP RESIDENZA ASSICURATO
		mov.City,                 //    CITTA' RESIDENZA ASSICURATO
		mov.CityCode,             //    PROVINCIA RESIDENZA ASSICURATO
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

	/*
		layout := "2006-01-02"
		if mov.MovementType == "delete" || mov.MovementType == "suspended" {
			e = lib.UpdateRowBigQuery("wopta", dataBanckAccount, map[string]string{
				"status":  mov.Status,
				"endDate": mov.EndDate.Format(layout) + " 00:00:00",
			}, "fiscalCode='"+mov.FiscalCode+"' and guaranteesCode='"+mov.GuaranteesCode+"'")

		}
	*/
	result = append(result, row)

	return result
}
func addMonthNow(mov inclusive.BankAccountMovement) string {
	now := time.Now().AddDate(0, 0, -1)
	if mov.MovementType == "delete" {
		endDate, _ := time.Parse("2006-01-02", mov.BigEndDate.Date.String())
		return endDate.Format(layout)
	} else {
		addMonth := now.AddDate(0, 1, 0)
		return addMonth.Format(layout)
	}
	return "31/12/9999"
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
		"CITTA' RESIDENZA ASSICURATO",          //    CITTA' RESIDENZA ASSICURATO
		"PROVINCIA RESIDENZA ASSICURATO",       //    PROVINCIA RESIDENZA ASSICURATO
		"DATA DI NASCITA ASSICURATO",           //    DATA DI NASCITA ASSICURATO
		"DATA INIZIO VALIDITA' COPERTURA",      //    DATA INIZIO VALIDITA' COPERTURA
		"DATA FINE VALIDITA' COPERTURA",        //    DATA FINE VALIDITA' COPERTURA
		"TIPO MOVIMENTO",
	}
}

func QueryRowsBigQuery[T any](query string) ([]T, error) {
	var (
		res  []T
		e    error
		iter *bigquery.RowIterator
	)
	log.Println(query)

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	lib.CheckError(err)
	defer client.Close()
	queryi := client.Query(query)
	iter, e = queryi.Read(ctx)

	for {
		var row T
		e = iter.Next(&row)

		if e == iterator.Done {
			log.Println(e)
			return res, nil
		}
		if e != nil {
			log.Println(e)
			return res, e
		}

		res = append(res, row)

	}

}
func setScalapayPolicy(code string, now time.Time, upload bool, req DataReq) {
	var (
		result [][]string
		refDay time.Time
		query  string
	)

	log.AddPrefix("BankAccountInclusive setScalapayPolicy " + code)
	defer log.PopPrefix()
	header := []string{
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
		"CITTA' RESIDENZA ASSICURATO",          //    CITTA' RESIDENZA ASSICURATO
		"PROVINCIA RESIDENZA ASSICURATO",       //    PROVINCIA RESIDENZA ASSICURATO          //    DATA DI NASCITA ASSICURATO
		"DATA INIZIO VALIDITA' COPERTURA",      //    DATA INIZIO VALIDITA' COPERTURA
		"DATA FINE VALIDITA' COPERTURA",
		"DATA VENDITA", //    DATA FINE VALIDITA' COPERTURA
		"TIPO MOVIMENTO",
	}
	if code == "51114" {
		code = "051114"
	}
	refDay = now.AddDate(0, 0, -1)
	log.Println("  refMontly: ", refDay)
	//from, e = time.Parse("2006-01-02", strconv.Itoa(now.Year())+"-"+fmt.Sprintf("%02d", int(now.Month()))+"-"+fmt.Sprintf("%02d", 1))
	//query := "select * from `wopta." + dataMovement + "` where _PARTITIONTIME >'" + from.Format(layoutQuery) + " 00:00:00" + "' and _PARTITIONTIME <'" + to.Format(layoutQuery) + " 23:59:00" + "'"
	if req.From != "" && req.To != "" {
		from, _ := time.Parse(layoutQuery, req.From)
		to, _ := time.Parse(layoutQuery, req.To)
		query = "select * from `wopta." + dataMovement + "` where startDate>'" + from.Format(layoutQuery) + "' and startDate<'" + to.Format(layoutQuery) + "' and policyNumber='" + code + "' and status='active'"
	} else {
		query = "select * from `wopta_inclusive." + usersTableScalapay + "` where daystart ='" + strconv.Itoa(refDay.Day()) + "' and policyNumber='" + code + "' and status='active'"
	}
	log.Println("  bigquery query: ", query)
	bankaccountlist, e := QueryRowsBigQuery[inclusive.BankAccountMovement](query)
	log.Println("  bigquery error: ", e)
	log.Println("  len(bankaccountlist): ", len(bankaccountlist))
	//result = append(result, getHeader())
	result = append(result, header)

	for i, mov := range bankaccountlist {
		log.Println(i)
		startDate, _ := time.Parse("2006-01-02", mov.BigStartDate.Date.String())
		row := []string{
			mov.PolicyNumber,         // NUMERO POLIZZA
			"H",                      //    LOB
			"C",                      //    TIPOLOGIA POLIZZA
			"",                       //    CODICE CONFIGURAZIONE
			"1",                      //    TIPO OGGETTO ASSICURATO
			mov.Id,                   //    IDENTIFICATIVO UNIVOCO APPLICAZIONE
			mov.FiscalCode,           //    CODICE FISCALE / P.IVA ASSICURATO
			mov.Surname,              //    COGNOME / RAGIONE SOCIALE ASSICURATO
			mov.Name,                 //    NOME ASSICURATO
			mov.Address,              //    INDIRIZZO RESIDENZA ASSICURATO
			mov.PostalCode,           //    CAP RESIDENZA ASSICURATO
			mov.City,                 //    CITTA' RESIDENZA ASSICURATO
			mov.CityCode,             //    PROVINCIA RESIDENZA ASSICURATO
			startDate.Format(layout), //    DATA INIZIO VALIDITA' COPERTURA
			addMonthNow(mov),         //    DATA FINE VALIDITA' COPERTURA
			startDate.Format(layout), //    DATA INIZIO VALIDITA' COPERTURA
			StringMapping(mov.MovementType, map[string]string{
				"":          "A",
				"active":    "A",
				"insert":    "A",
				"delete":    "E",
				"suspended": "E",
			})}

		result = append(result, row)

	}

	filepath := code + "_" + strconv.Itoa(refDay.Year()) + fmt.Sprintf("%02d", int(refDay.Month())) + fmt.Sprintf("%02d", refDay.Day())
	CreateExcel(result, "../tmp/"+filepath+".xlsx", "Sheet1")
	lib.WriteCsv("../tmp/"+filepath+".csv", result, ';')
	source, _ := os.ReadFile("../tmp/" + filepath + ".xlsx")
	sourceCsv, _ := os.ReadFile("../tmp/" + filepath + ".csv")
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/scalapay/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".xlsx", source)
	lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/axa/inclusive/scalapay/"+strconv.Itoa(refDay.Year())+"/"+fmt.Sprintf("%02d", int(refDay.Month()))+"/"+filepath+".csv", sourceCsv)
	if upload {

		AxaSftpUpload(filepath+".xlsx", "SCALAPAY/IN/")
	}

}
