package companydata

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oliveagle/jsonpath"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

const localBasePath = "../tmp/"

func ProductTrackOutFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from          time.Time
		to            time.Time
		procuctTrack  Track
		result        [][]string
		event         []Column
		db            Database
		policies      []models.Policy
		transactions  []models.Transaction
		eventFilename string
	)
	log.SetPrefix("ProductTrackOutFx ")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	now, upload, reqData := getCompanyDataReq(req)

	procuctTrackByte := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+reqData.Name+"/v1/track.json", "")

	err := json.Unmarshal(procuctTrackByte, &procuctTrack)
	lib.CheckError(err)
	log.Println("Product track: ", procuctTrack)
	for _, ev := range reqData.Event {
		switch ev {

		case "emit":
			db = procuctTrack.Emit.Database
			event = procuctTrack.Emit.Event
			eventFilename = procuctTrack.Emit.FileName
		case "payment":
			db = procuctTrack.Payment.Database
			event = procuctTrack.Payment.Event
			eventFilename = procuctTrack.Payment.FileName
		case "delete":
			db = procuctTrack.Delete.Database
			event = procuctTrack.Delete.Event
			eventFilename = procuctTrack.Delete.FileName

		}

		from, to = procuctTrack.frequency(now)
		switch db.Dataset {
		case "policy":
			policies = query[models.Policy](from, to, db)
			result = procuctTrack.PolicyProductTrack(policies, event)
		case "transaction":
			transactions = query[models.Transaction](from, to, db)
			result = procuctTrack.TransactionProductTrack(transactions, event)
		}
		result = procuctTrack.makeHeader(event, result, procuctTrack.HasHeader)
		filename, byteArray := procuctTrack.saveFile(result, from, to, now, eventFilename)
		log.Println("len(result): ", len(result))
		if len(result) > 0 {

			if upload {
				procuctTrack.upload(filename)
			}
			if procuctTrack.SendMail {
				procuctTrack.sendMail(filename, byteArray)
			}
		}
	}
	return "", nil, e
}
func (track Track) PolicyProductTrack(policies []models.Policy, event []Column) [][]string {
	var (
		result [][]string
		err    error
	)

	for _, policy := range policies {
		if track.IsAssetFlat {
			result = track.policyAssetRow(&policy, event)
		} else {
			result = track.policyAssetRow(&policy, event)
		}
		//docsnap := lib.GetFirestore("policy", transaction.PolicyUid)
		//docsnap.DataTo(&policy)
		//result = append(result, procuctTrack.ProductTrack(policy, event)...)

	}

	lib.CheckError(err)
	return result
}
func (track Track) TransactionProductTrack(transactions []models.Transaction, event []Column) [][]string {
	var (
		result    [][]string
		json_data interface{}
		err       error
	)

	for indexG, _ := range transactions {
		log.Println("index transactions: ", indexG)
		var cells []string
		for i, column := range event {
			log.Println("index event:", i)
			var resPaths []interface{}

			for _, value := range column.Values {
				log.Println("column value", value)
				if strings.Contains(value, "$.") {

					log.Println("column value: ", value)
					resPath, err := jsonpath.JsonPathLookup(json_data, value)
					resPaths = append(resPaths, resPath)
					log.Println(err)
					log.Println(resPath)
				} else {
					resPaths = append(resPaths, value)
				}

			}
			resdata := checkMap(column, resPaths)
			cells = append(cells, checkType(resdata))

		}
		result = append(result, cells)
	}

	lib.CheckError(err)
	return result
}
func (track Track) saveFile(matrix [][]string, from time.Time, to time.Time, now time.Time, baseFilename string) (string, []byte) {
	var byteArray []byte
	filename := stringTimeToken(baseFilename, from, to, now)
	filepath := "track/" + track.Name + "/" + strconv.Itoa(from.Year()) + "/" + filename
	log.Println("filepath: ", filepath)
	switch track.Type {
	case "csv":

		sep := []rune(track.CsvConfig.Separator)
		e := lib.WriteCsv(localBasePath+filename, matrix, sep[0])
		lib.CheckError(e)
		source, e := os.ReadFile(localBasePath + filename)
		lib.CheckError(e)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath, source)
	case "excel":
		byteArray, e := CreateExcel(matrix, localBasePath+filename, track.ExcelConfig.SheetName)
		lib.CheckError(e)
		//source, e := os.ReadFile(localBasePath + filename)
		lib.CheckError(e)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath, byteArray)

	}
	return filename, byteArray
}

func (track Track) policyAssetRow(policy *models.Policy, event []Column) [][]string {
	var (
		json_data interface{}
		result    [][]string

		err error
	)
	log.Println("policyAssetRow")
	b, err := json.Marshal(policy)
	lib.CheckError(err)
	json.Unmarshal(b, &json_data)
	log.Println(string(b))
	for indexAsset, asset := range policy.Assets {
		for indexG, _ := range asset.Guarantees {
			log.Println("index Guarantees: ", indexG)
			var cells []string
			for i, column := range event {
				log.Println("index event: ", i)
				var resPaths []interface{}

				for _, value := range column.Values {
					log.Println("column value", value)
					if strings.Contains(value, "$.") {
						value = strings.Replace(value, "guarantees[*]", "guarantees["+strconv.Itoa(indexG)+"]", 1)
						value = strings.Replace(value, "assets[*]", "assets["+strconv.Itoa(indexAsset)+"]", 1)
						log.Println("column value guarantee: ", value)
						resPath, err := jsonpath.JsonPathLookup(json_data, value)
						resPaths = append(resPaths, resPath)
						log.Println(err)
						log.Println(resPath)
					} else {
						resPaths = append(resPaths, value)
					}

				}
				resdata := checkMap(column, resPaths)
				cells = append(cells, checkType(resdata))

			}
			result = append(result, cells)
		}

	}
	return result

}
func (track Track) upload(filePath string) {
	switch track.UploadType {
	case "sftp":
		track.sftp(filePath)
	}

}
func checkMap(column Column, value []interface{}) interface{} {
	var res interface{}
	res = value[0]
	log.Println("column.MapFx: ", column.MapFx)
	if column.MapFx != "" {
		res = GetMapFx(column.MapFx, value)
	}
	if column.MapStatic != nil && value[0] != nil {
		res = column.MapStatic[value[0].(string)]
	}
	if res == nil {
		res = ""
	}
	return res

}
func (track *Track) frequency(now time.Time) (time.Time, time.Time) {
	location, e := time.LoadLocation("Europe/Rome")
	lib.CheckError(e)
	switch track.Frequency {
	case "monthly":
		prevMonth := lib.GetPreviousMonth(now)
		from = lib.GetFirstDay(prevMonth)
		to = lib.GetFirstDay(now)
	case "daily":
		prevDay := now.AddDate(0, 0, -1)
		from = time.Date(prevDay.Year(), prevDay.Month(), prevDay.Day(), 0, 0, 0, 0, location)
		to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	}
	log.Println(from, to)
	track.from = from
	track.to = to
	track.now = now
	return from, to
}

func query[T any](from time.Time, to time.Time, db Database) []T {
	var res []T
	switch db.Name {

	case "firestore":
		res = firestoreQuery[T](from, to, db)

	}

	return res
}
func firestoreQuery[T any](from time.Time, to time.Time, db Database) []T {
	var value interface{}
	firequery := lib.FireGenericQueries[T]{
		Queries: []lib.Firequery{},
	}
	for _, qe := range db.Query {
		value = qe.QueryValue
		if qe.QueryValue == "from" {
			value = from
		}
		if qe.QueryValue == "to" {
			value = to
		}

		firequery.Queries = append(firequery.Queries,
			lib.Firequery{
				Field:      qe.Field,    //
				Operator:   qe.Operator, //
				QueryValue: value,
			})
	}
	res, _, e := firequery.FireQueryUid(db.Dataset)
	lib.CheckError(e)
	return res
}
func BigQuery[T any](from time.Time, to time.Time, db Database) []T {
	var value interface{}
	firequery := lib.FireGenericQueries[T]{
		Queries: []lib.Firequery{},
	}

	for _, qe := range db.Query {
		value = qe.QueryValue
		if qe.QueryValue == "from" {
			value = from
		}
		if qe.QueryValue == "to" {
			value = to
		}

		firequery.Queries = append(firequery.Queries,
			lib.Firequery{
				Field:      qe.Field,    //
				Operator:   qe.Operator, //
				QueryValue: value,
			})
	}
	res, _, e := firequery.FireQueryUid(db.Dataset)
	lib.CheckError(e)
	return res
}
func stringTimeToken(filename string, from time.Time, to time.Time, now time.Time) string {
	filename = strings.Replace(filename, "fdd", fmt.Sprintf("%02d", from.Day()), 1)
	filename = strings.Replace(filename, "fmm", fmt.Sprintf("%02d", int(from.Month())), 1)
	filename = strings.Replace(filename, "fyyyy", fmt.Sprint(from.Year()), 1)

	return filename
}
func (track Track) sftp(filePath string) {

	pk := lib.GetFromStorage(os.Getenv(track.FtpConfig.PrivateKey), "env/axa-life.ppk", "")
	config := lib.SftpConfig{
		Username:     track.FtpConfig.Username,
		Password:     track.FtpConfig.Password,                                                                                    // required only if password authentication is to be used
		PrivateKey:   string(pk),                                                                                                  //                           // required only if private key authentication is to be used
		Server:       track.FtpConfig.Server,                                                                                      //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       "", // 0 for not timeout
	}
	client, e := lib.NewSftpclient(config)
	lib.CheckError(e)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open(localBasePath + filePath)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files(localBasePath)
	destination, e := client.Create(track.FtpConfig.Path + "/" + filePath)
	lib.CheckError(e)
	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	lib.CheckError(e)
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))
	lib.CheckError(e)

}
func (track *Track) sendMail(filename string, byteArray []byte) {
	log.Println("sendMail: ")
	var source []byte
	var contentType string
	if byteArray != nil {
		source = byteArray
	} else {
		source, _ = os.ReadFile("../tmp/" + filename)
	}
	if track.Type == "csv" {
		contentType = "text/csv"
	}
	if track.Type == "excel" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	at := &[]mail.Attachment{{
		Byte:        base64.StdEncoding.EncodeToString(source),
		ContentType: contentType,
		FileName:    filename,
		Name:        strings.ReplaceAll(filename, "_", " "),
	}}

	mail.SendMail(mail.MailRequest{
		From:         track.MailConfig.From,
		To:           track.MailConfig.To,
		Cc:           track.MailConfig.Cc,
		Bcc:          track.MailConfig.Bcc,
		Message:      track.MailConfig.Message,
		Subject:      track.MailConfig.Subject,
		IsAttachment: true,
		Attachments:  at,
	})
}
func (track *Track) makeHeader(event []Column, data [][]string, hasHeader bool) [][]string {

	var res [][]string
	var header [][]string
	var row []string
	if hasHeader {
		for _, colums := range event {
			row = append(row, colums.Name)
		}
		header = append(res, row)
		res = append(header, data...)
	} else {
		res = data
	}
	return res
}
func checkType(i interface{}) string {

	var res string
	switch v := i.(type) {
	case int:
		res = strconv.Itoa(i.(int))
	case string:
		res = i.(string)
	case float64:
		res = fmt.Sprint(i.(float64))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	return res
}
