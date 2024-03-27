package form

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func AxaFleet(w http.ResponseWriter, r *http.Request,spreadsheetSource string , spreadsheetArch string) ( [][]interface{}, bool, error) {
	log.Println("AxaFleetEmit")
	log.Println(os.Getenv("env"))
	var (
		path                  []byte
		excel                 [][]interface{}
		e                     error
		toEmit                bool
		typeMov               string
		DATAIMMATRICOLAZIONE  string
		MODELLO               string
		TARGA                 string
		DATAVALIDITACOPERTURA string
		DATAFINECOPERTURA     string
		deleteList            []string
		insertList            []string
		sequence              int
	)
	const (
		satusCol = "J"
	)

	switch os.Getenv("env") {
	case "local":
		path = lib.ErrorByte(os.ReadFile("function-data/sa/positive-apex-350507-33284d6fdd55.json"))
	case "dev":
		path = lib.GetFromStorage("function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")
	case "prod":
		path = lib.GetFromStorage("core-350507-function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")

	default:

	}
	ctx := context.Background()
	srv, e := sheets.NewService(ctx, option.WithCredentialsJSON(path), option.WithScopes(sheets.SpreadsheetsScope))
	fmt.Println(e)
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := spreadsheetSource
	spreadsheettotId := spreadsheetArch
	tway, e := srv.Spreadsheets.Values.Get(spreadsheetId, "A:J").Do()
	fmt.Println(e)
	axa, e := srv.Spreadsheets.Values.Get(spreadsheettotId, "A:W").Do()
	fmt.Println(e)
	excelhead := []interface{}{"NUMERO POLIZZA", "LOB", "	TIPOLOGIA POLIZZA", "CODICE CONFIGURAZIONE", "IDENTIFICATIVO UNIVOCO APPLICAZIONE", "	TIPO OGGETTO ASSICURATO", "	CODICE FISCALE / P.IVA ASSICURATO", "COGNOME / RAGIONE SOCIALE ASSICURATO", "	NOME ASSICURATO", "	INDIRIZZO RESIDENZA ASSICURATO", "	CAP RESIDENZA ASSICURATO", "	CITTA’ RESIDENZA ASSICURATO", "	PROVINCIA RESIDENZA ASSICURATO", "	TARGA VEICOLO", "	TELAIO VEICOLO	", "MARCA VEICOLO", "	MODELLO VEICOLO	TIPOLOGIA VEICOLO", "PESO VEICOLO", "	DATA IMMATRICOLAZIONE", "	DATA INIZIO VALIDITA' COPERTURA", "	DATA FINE VALIDITA' COPERTURA", "TIPO MOVIMENTO"}
	excel = append(excel, excelhead)
	toEmit = false
	sequence = 0
	if len(tway.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range tway.Values {
			isError := false
			founded := false
			fmt.Println(axa.Values[len(axa.Values)-1][4])
			fmt.Println(axa.Values[len(axa.Values)-1][21])
			lenTableDelta := 1
			for i := 1; i < 100; i++ {
				if axa.Values[len(axa.Values)-i][22] != "E" {
					lenTableDelta = i
					break
				}
			}
			marks, _ := strconv.Atoi(axa.Values[len(axa.Values)-lenTableDelta][4].(string)[2:10])
			progressive := marks + 1 + sequence

			progressiveFormatted := fmt.Sprintf("%08d", progressive)
			progressiveFormattedpre := "WR" + progressiveFormatted

			if len(row) <= 9 && i != 0 {
				toEmit = true
				//mailsource = row[1].(string)
				fmt.Println("Enter in No EMESSO")
				fmt.Println(row)
				if row[6] == "Inserimento" {
					typeMov = "A"
					DATAVALIDITACOPERTURA = row[5].(string)
					DATAIMMATRICOLAZIONE = row[4].(string)
					TARGA = row[2].(string)
					MODELLO = row[3].(string)
					now:=time.Now()
					DATAFINECOPERTURA = "31/12/"+strconv.Itoa(now.Year())
					sequence++

				} else {
					for x, axarow := range axa.Values {
						// var t string
						if strings.ToUpper(axarow[13].(string)) == strings.ToUpper(row[7].(string)) {
							fmt.Println("axarow[13] == row[2]: ",x)
							progressiveFormattedpre = axarow[4].(string)
							founded = true
							DATAVALIDITACOPERTURA = axarow[20].(string)
							DATAIMMATRICOLAZIONE = axarow[19].(string)
							TARGA = axarow[13].(string)
							MODELLO = axarow[16].(string)
							DATAFINECOPERTURA = row[8].(string)

						}

					}
					if !founded {
						mail.SendMail((getMailObj("<p>Opss.. qualcosa è andato storto</p><p>Il servizio di aggiornamento copertura flotte di Wopta per T-way non è stato in grado di trovare la targa: "+row[2].(string)+"</p><p>Non ti preoccupare questa operazione è stata gia annullata devi solo rieffetuare la richiesta dall' apposito form con la targa corretta</p>",
							row[1].(string))))
						fmt.Println("ERROR save, :")
						celindex := strconv.Itoa(i + 1)
						cel := &sheets.ValueRange{
							Values: [][]interface{}{{"ERRATO"}},
						}
						_, e = srv.Spreadsheets.Values.Update(spreadsheetId, satusCol+celindex+":"+satusCol+celindex, cel).ValueInputOption("USER_ENTERED").Context(ctx).Do()
						isError = true
					}
					typeMov = "E"
				}
				if !isError {
					if founded {
						deleteList = append(deleteList, TARGA)
					} else {
						insertList = append(insertList, TARGA)
					}
					fmt.Println("!isError")
					celindex := strconv.Itoa(i + 1)
					excelRow := []interface{}{"191222", "A", "C", "00001", progressiveFormattedpre, "2", "03682240043", "T-WAY SPA", "", "Piazza Walther Von Der Vogelweide", "39100", "Bolzano", "BZ", TARGA, "", "", MODELLO, 3, 4, DATAIMMATRICOLAZIONE, DATAVALIDITACOPERTURA, DATAFINECOPERTURA, typeMov}
					cel := &sheets.ValueRange{
						Values: [][]interface{}{{"EMESSO"}},
					}
					row := &sheets.ValueRange{
						Values: [][]interface{}{excelRow},
					}

					excel = append(excel, excelRow)
					fmt.Println("first save, :")
					_, e = srv.Spreadsheets.Values.Update(spreadsheetId, satusCol+celindex+":"+satusCol+celindex, cel).ValueInputOption("USER_ENTERED").Context(ctx).Do()
					fmt.Println(e)
					fmt.Println("second save:")
					_, e = srv.Spreadsheets.Values.Append(spreadsheettotId, "Foglio1", row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
				}
			}

		}
	}
	

	return excel, toEmit, e

}
