package form

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/xuri/excelize/v2"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// {"responses":{"TIPO MOVIMENTO":"Inserimento","Targa Inserimento":"test","MODELLO VEICOLO":"test mod","DATA IMMATRICOLAZIONE":"1212-12-02","DATA INIZIO VALIDITA' COPERTURA":"1212-12-12"}}
// {"responses":{"TIPO MOVIMENTO":"Annullo","Targa Annullo":"targa","DATA FINE VALIDITA' COPERTURA":"0009-09-09"},"mail":"test@gmail.com"}

func AxaFleetTway(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("AxaFleetEmit")
	log.Println(os.Getenv("env"))
	var (
		path                  []byte
		excel                 [][]interface{}
		e                     error
		toEmit                bool
		sourcest              []byte
		typeMov               string
		mailsource            string
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

	path = lib.GetFilesByEnv("sa/positive-apex-350507-33284d6fdd55.json")

	ctx := context.Background()
	srv, e := sheets.NewService(ctx, option.WithCredentialsJSON(path), option.WithScopes(sheets.SpreadsheetsScope))
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := "1S3ELzCRXHMT0xarthgUWof9PMhFolmj9wE0G9I9ItWw"
	spreadsheettotId := "1UtYiPt7fJ8FAZQRpmwRZpyqGpOT26Q2-qNidLAAPFlQ"
	tway, e := srv.Spreadsheets.Values.Get(spreadsheetId, "A:J").Do()
	axa, e := srv.Spreadsheets.Values.Get(spreadsheettotId, "A:W").Do()
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
				mailsource = row[1].(string)
				fmt.Println("Enter in No EMESSO")
				fmt.Println(row)
				if row[6] == "Inserimento" {
					typeMov = "A"
					DATAVALIDITACOPERTURA = row[5].(string)
					DATAIMMATRICOLAZIONE = row[4].(string)
					TARGA = row[2].(string)
					MODELLO = row[3].(string)
					DATAFINECOPERTURA = "31/12/2024"
					sequence++

				} else {
					for x, axarow := range axa.Values {
						// var t string
						if strings.ToUpper(axarow[13].(string)) == strings.ToUpper(row[7].(string)) {
							fmt.Println("axarow[13] == row[2]")
							fmt.Println(x)
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
					fmt.Println("second save:")
					_, e = srv.Spreadsheets.Values.Append(spreadsheettotId, "Foglio1", row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
				}
			}

		}
	}

	if toEmit {

		now := time.Now()
		layout2 := "2006-01-02"
		filepath := now.Format(layout2) + "-" + strconv.Itoa(time.Now().Nanosecond()) + ".xlsx"

		if os.Getenv("env") != "local" {
			//./serverless_function_source_code/
			//filepath = "../tmp/" + filepath
		}
		sourcest, e = CreateExcel(excel, "../tmp/"+filepath)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "tway-fleet-axa/"+filepath, sourcest)
		SftpUpload(filepath)
		mail.SendMail((getMailObj("<p>inserite in copertura: </p><p>"+strings.Join(insertList, ", ")+"</p><p>escluse dalla copertura: </p><p>"+strings.Join(deleteList, ", ")+"</p>",
			mailsource)))
	}

	return "", nil, e

}

func SftpUpload(filePath string) {
	var (
		pk []byte
		e  error
	)

	pk = lib.GetFilesByEnv("env/twayserviceKey.ssh")

	lib.CheckError(e)

	//ssh: handshake failed: ssh: no common algorithm for key exchange; client offered: [diffie-hellman-group-exchange-sha256 diffie-hellman-group14-sha256 ext-info-c], server offered: [diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1]
	//diffie-hellman-group-exchange-sha1 diffie-hellman-group1-sha1 diffie-hellman-group14-sha1
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_SFTP_USER"),
		Password:     os.Getenv("AXA_SFTP_PSW"), // required only if password authentication is to be used
		PrivateKey:   string(pk),                // required only if private key authentication is to be used
		Server:       "ftp.ip-assistance.it:22",
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       os.Getenv("AXA_SFTP_PSW"), // 0 for not timeout
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
	destination, e := client.Create("IN/" + filePath)
	lib.CheckError(e)
	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	lib.CheckError(e)
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

func CreateExcel(sheet [][]interface{}, filePath string) ([]byte, error) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	alfabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	lib.CheckError(err)
	for x, row := range sheet {
		for i, cel := range row {

			fmt.Println(cel)
			f.SetCellValue("Sheet1", alfabet[i]+""+strconv.Itoa(x+1), cel)
		}
	}
	//Set active sheet of the workbook.
	f.SetActiveSheet(index)

	//Save spreadsheet by the given path.
	err = f.SaveAs(filePath)

	resByte, err := f.WriteToBuffer()

	return resByte.Bytes(), err
}

func getMailObj(msg string, mailsource string) mail.MailRequest {
	//link := "https://storage.googleapis.com/documents-public-dev/information-set/information-sets//v1/Precontrattuale.pdf "michele.lomazzi@wopta.it","
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{
		"assunzione@wopta.it",
		"luca.barbieri@wopta.it",
		"beatrice.sala@wopta.it",
		mailsource,
	}
	obj.Message = msg
	obj.Subject = " Wopta T-Way Axa Fleet"
	obj.IsHtml = true
	obj.IsAttachment = false
	obj.IsLink = true
	obj.Link = "https://docs.google.com/spreadsheets/d/1UtYiPt7fJ8FAZQRpmwRZpyqGpOT26Q2-qNidLAAPFlQ/edit#gid=0"
	obj.LinkLabel = "Archivio completo"
	obj.Title = ""
	obj.SubTitle = ""

	return obj
}
