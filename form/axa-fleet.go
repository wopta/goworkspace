package form

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func GetFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("AxaFleetEmit")
	log.Println(os.Getenv("env"))
	var path []byte
	var excel [][]interface{}
	var e error

	switch os.Getenv("env") {
	case "local":
		path = lib.ErrorByte(ioutil.ReadFile("function-data/sa/positive-apex-350507-33284d6fdd55.json"))
	case "dev":
		path = lib.GetFromStorage("function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")
	case "prod":
		path = lib.GetFromStorage("core-350507-function-data", "sa/positive-apex-350507-33284d6fdd55.json", "")

	default:

	}
	ctx := context.Background()

	log.Println(string(path))
	srv, e := sheets.NewService(ctx, option.WithCredentialsJSON(path), option.WithScopes(sheets.SpreadsheetsScope))
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := "1S3ELzCRXHMT0xarthgUWof9PMhFolmj9wE0G9I9ItWw"
	spreadsheettotId := "1UtYiPt7fJ8FAZQRpmwRZpyqGpOT26Q2-qNidLAAPFlQ"
	tway, e := srv.Spreadsheets.Values.Get(spreadsheetId, "A:H").Do()
	axa, e := srv.Spreadsheets.Values.Get(spreadsheettotId, "A:W").Do()

	if len(tway.Values) == 0 {
		fmt.Println("No data found.")
	} else {

		for i, row := range tway.Values {

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
			progressive := marks + 1
			progressiveFormatted := fmt.Sprintf("%08d", progressive)
			progressiveFormattedpre := "WR" + progressiveFormatted
			//"NUMERO POLIZZA","LOB","	TIPOLOGIA POLIZZA","CODICE CONFIGURAZIONE","IDENTIFICATIVO UNIVOCO APPLICAZIONE","	TIPO OGGETTO ASSICURATO","	CODICE FISCALE / P.IVA ASSICURATO","COGNOME / RAGIONE SOCIALE ASSICURATO","	NOME ASSICURATO","	INDIRIZZO RESIDENZA ASSICURATO","	CAP RESIDENZA ASSICURATO","	CITTA’ RESIDENZA ASSICURATO","	PROVINCIA RESIDENZA ASSICURATO","	TARGA VEICOLO","	TELAIO VEICOLO	","MARCA VEICOLO","	MODELLO VEICOLO	TIPOLOGIA VEICOLO","PESO VEICOLO","	DATA IMMATRICOLAZIONE","	DATA INIZIO VALIDITA' COPERTURA","	DATA FINE VALIDITA' COPERTURA","TIPO MOVIMENTO"

			if len(row) == 7 && i != 0 {
				excelhead := []interface{}{"NUMERO POLIZZA", "LOB", "	TIPOLOGIA POLIZZA", "CODICE CONFIGURAZIONE", "IDENTIFICATIVO UNIVOCO APPLICAZIONE", "	TIPO OGGETTO ASSICURATO", "	CODICE FISCALE / P.IVA ASSICURATO", "COGNOME / RAGIONE SOCIALE ASSICURATO", "	NOME ASSICURATO", "	INDIRIZZO RESIDENZA ASSICURATO", "	CAP RESIDENZA ASSICURATO", "	CITTA’ RESIDENZA ASSICURATO", "	PROVINCIA RESIDENZA ASSICURATO", "	TARGA VEICOLO", "	TELAIO VEICOLO	", "MARCA VEICOLO", "	MODELLO VEICOLO	TIPOLOGIA VEICOLO", "PESO VEICOLO", "	DATA IMMATRICOLAZIONE", "	DATA INIZIO VALIDITA' COPERTURA", "	DATA FINE VALIDITA' COPERTURA", "TIPO MOVIMENTO"}
				excel = append(excel, excelhead)
				var typeMov string
				isError := false
				if row[5] == "Inserimento" {
					typeMov = "A"
				} else {
					founded := false

					for x, axarow := range axa.Values {
						// var t string
						if axarow[13] == row[2] {
							fmt.Println("axarow[13] == row[2]")
							fmt.Println(x)
							progressiveFormattedpre = axarow[4].(string)
							founded = true
						}

					}
					if !founded {
						mail.SendMail((getMailObj("<p>Opss.. qualcosa è andato storto</p><p>Il servizio di aggiornamento copertura flotte di Wopta per T-way non è stato in grado di trovare la targa: " + row[2].(string) + "</p><p>Non ti preoccupare questa operazione è stata gia annullata devi solo rieffetuare la richiesta dall' apposito form con la targa corretta</p>")))
						isError = true
					}
					typeMov = "E"
				}
				if !isError {
					celindex := strconv.Itoa(i + 1)
					excelRow := []interface{}{"191222", "A", "C", "00001", progressiveFormattedpre, "2", "03682240043", "T-WAY SPA", "", "Piazza Walther Von Der Vogelweide", "39100", "Bolzano", "BZ", row[1], "", "", row[2], 3, 4, row[3], row[4], "31/12/2023", typeMov}
					cel := &sheets.ValueRange{
						Values: [][]interface{}{{"EMESSO"}},
					}
					row := &sheets.ValueRange{
						Values: [][]interface{}{{"191222", "A", "C", "00001", progressiveFormattedpre, "2", "03682240043", "T-WAY SPA", "", "Piazza Walther Von Der Vogelweide", "39100", "Bolzano", "BZ", row[1], "", "", row[2], 3, 4, row[3], row[4], "31/12/2023", typeMov}},
					}
					excel = append(excel, excelRow)
					fmt.Println("first save, :")
					_, e = srv.Spreadsheets.Values.Update(spreadsheetId, "H"+celindex+":H"+celindex, cel).ValueInputOption("USER_ENTERED").Context(ctx).Do()
					fmt.Println("second save:")
					_, e = srv.Spreadsheets.Values.Append(spreadsheettotId, "Foglio1", row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
				}
			}

		}
	}
	now := time.Now()
	layout2 := "2006-01-02"
	filepath := now.Format(layout2) + "-" + strconv.Itoa(time.Now().Nanosecond()) + ".xlsx"

	if os.Getenv("env") != "local" {
		//./serverless_function_source_code/

		//filepath = "../tmp/" + filepath
	}
	lib.CreateExcel(excel, "../tmp/"+filepath)
	//root = path.dirname(path.abspath(__file__))
	log.Println("tempdir")
	lib.Files("../tmp")
	sourcest, e := ioutil.ReadFile("../tmp/" + filepath)
	lib.PutToStorage("function-data", "tway-fleet-axa/"+filepath, sourcest)
	SftpUpload(filepath)
	return "", nil, e
}
func SftpUpload(filePath string) {
	var (
		pk []byte
		e  error
	)
	switch os.Getenv("env") {
	case "local":
		pk = lib.ErrorByte(ioutil.ReadFile("function-data/env/twayserviceKey.ssh"))
	case "dev":
		pk = lib.GetFromStorage("function-data", "env/twayserviceKey.ssh", "")

	case "prod":
		pk = lib.GetFromStorage("core-350507-function-data", "env/twayserviceKey.ssh", "")
	default:

	}

	lib.CheckError(e)

	//ssh: handshake failed: ssh: no common algorithm for key exchange; client offered: [diffie-hellman-group-exchange-sha256 diffie-hellman-group14-sha256 ext-info-c], server offered: [diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1]
	//diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1
	config := lib.Config{
		Username:     os.Getenv("AXA_SFTP_USER"),
		Password:     os.Getenv("AXA_SFTP_PSW"), // required only if password authentication is to be used
		PrivateKey:   string(pk),                // required only if private key authentication is to be used
		Server:       "ftp.ip-assistance.it:22",
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	lib.CheckError(e)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open("../tmp/" + filePath)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files("../tmp")
	//destination, e := client.Create(filePath)
	destination, e := client.OpenFile("./" + filePath)
	lib.CheckError(e)
	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))

	lib.CheckError(e)
	/*
		// Download remote file.
		file, err := client.Download("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		// Read downloaded file.
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))

		// Get remote file stats.
		info, err := client.Info("tmp/file.txt")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", info)*/
}
func getMailObj(msg string) mail.MailRequest {

	//link := "https://storage.googleapis.com/documents-public-dev/information-set/information-sets//v1/Precontrattuale.pdf "michele.lomazzi@wopta.it","
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{"luca.barbieri@wopta.it"}
	obj.Message = msg
	obj.Subject = " Wopta T-Way Axa Fleet"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
