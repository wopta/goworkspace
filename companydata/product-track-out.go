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

func ProductTrackOutFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		from         time.Time
		to           time.Time
		procuctTrack Track
		result       [][]string
		event        []Column
		db           Database
		policies     []models.Policy
		transactions []models.Transaction
	)
	log.SetPrefix("ProductTrackOutFx ")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	now, upload, reqData := getCompanyDataReq(req)

	procuctTrackByte := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+reqData.Name+"/v1/track.json", "")

	err := json.Unmarshal(procuctTrackByte, &procuctTrack)
	lib.CheckError(err)
	log.Println("Product track: ", procuctTrack)

	switch reqData.Event {

	case "emit":
		db = procuctTrack.Emit.Database
		event = procuctTrack.Emit.Event
	case "payment":
		db = procuctTrack.Payment.Database
		event = procuctTrack.Payment.Event
	case "delete":
		db = procuctTrack.Delete.Database
		event = procuctTrack.Delete.Event

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

	filename := procuctTrack.saveFile(result, from, to, now)

	if upload {
		procuctTrack.upload(filename)
	}
	if procuctTrack.SendMail {
		procuctTrack.sendMail(filename)
	}
	return "", nil, e
}
func (track Track) PolicyProductTrack(policies []models.Policy, event []Column) [][]string {
	var (
		result [][]string

		err error
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

		value string
		cells []string
		err   error
	)

	for _, tr := range transactions {
		b, err := json.Marshal(tr)
		lib.CheckError(err)
		json.Unmarshal(b, &json_data)
		log.Println(string(b))
		for _, column := range event {

			log.Println(column)
			log.Println(value)
			var resPaths []interface{}
			for _, value := range column.Values {
				log.Println("column value", value)
				if strings.Contains(value, "$.") {
					resPath, err := jsonpath.JsonPathLookup(json_data, value)
					resPaths = append(resPaths, resPath)
					log.Println(err)
					log.Println(resPath)
				}
			}
			resdata := checkMap(column, resPaths)
			cells = append(cells, resdata.(string))

		}

		result = append(result, cells)

	}

	lib.CheckError(err)
	return result
}
func (track Track) saveFile(matrix [][]string, from time.Time, to time.Time, now time.Time) string {
	const localBasePath = "../tmp/"
	filename := track.formatFilename(track.FileName, from, to, now)
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
		_, e := lib.CreateExcel(matrix, filename, "Risultato")
		lib.CheckError(e)
		source, e := os.ReadFile(localBasePath + filename)
		lib.CheckError(e)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), filepath, source)

	}
	return filepath
}

func (track Track) policyAssetRow(policy *models.Policy, event []Column) [][]string {
	var (
		json_data interface{}
		result    [][]string
		cells     []string
		err       error
	)
	log.Println("policyAssetRow")
	b, err := json.Marshal(policy)
	lib.CheckError(err)
	json.Unmarshal(b, &json_data)
	log.Println(string(b))
	for indexAsset, asset := range policy.Assets {
		for indexG, _ := range asset.Guarantees {
			log.Println("index Guarantees: ", indexG)
			for i, column := range event {
				log.Println("index event: ", i)
				var resPaths []interface{}
				for _, value := range column.Values {
					log.Println("column value", value)
					if strings.Contains(value, "$.") {

						value = strings.Replace(value, "guarantees[*]", "guarantees["+strconv.Itoa(indexG)+"]", 1)
						value = strings.Replace(value, "assets[*]", "assets["+strconv.Itoa(indexAsset)+"]", 1)
						resPath, err := jsonpath.JsonPathLookup(json_data, value)
						resPaths = append(resPaths, resPath)
						log.Println(err)
						log.Println(resPath)
					}
				}
				resdata := checkMap(column, resPaths)
				cells = append(cells, resdata.(string))

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
func (track Track) frequency(now time.Time) (time.Time, time.Time) {
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

func (track Track) formatFilename(filename string, from time.Time, to time.Time, now time.Time) string {
	filename = strings.Replace(filename, "fdd", fmt.Sprint(from.Day()), 1)
	filename = strings.Replace(filename, "fmm", fmt.Sprint(from.Month()), 1)
	filename = strings.Replace(filename, "fyyyy", fmt.Sprint(from.Year()), 1)

	return filename
}
func (track Track) sftp(filePath string) {

	pk := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "env/axa-life.ppk", "")
	config := lib.SftpConfig{
		Username:     os.Getenv("AXA_LIFE_SFTP_USER"),
		Password:     "",                                                                                                          // required only if password authentication is to be used
		PrivateKey:   string(pk),                                                                                                  //                           // required only if private key authentication is to be used
		Server:       os.Getenv("AXA_LIFE_SFTP_HOST") + ":10026",                                                                  //
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,
		KeyPsw:       "", // 0 for not timeout
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
	destination, e := client.Create("To_CLP/" + filePath)
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
func (track Track) sendMail(filename string) {
	var contentType string
	source, _ := os.ReadFile("../tmp/" + filename)
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
