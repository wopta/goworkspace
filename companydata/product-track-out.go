package companydata

import (
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
	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	now, upload, reqData := getCompanyDataReq(req)

	procuctTrackByte := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+reqData.Name+"/v1/track.json", "")
	json.Unmarshal(procuctTrackByte, &procuctTrack)

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

	filepath := procuctTrack.saveFile(result, from, to, now)

	if upload {
		procuctTrack.upload(filepath)
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

		cells []string
		err   error
	)

	for _, tr := range transactions {
		b, err := json.Marshal(tr)
		lib.CheckError(err)
		json.Unmarshal(b, &json_data)
		log.Println(string(b))
		for _, column := range event {
			var (
				resPath interface{}
				value   string
			)
			value = column.Value
			resPath = column.Value
			log.Println(column.Value)
			log.Println(value)
			resPath, err = jsonpath.JsonPathLookup(json_data, value)
			lib.CheckError(err)
			log.Println(resPath)

			resPath = checkMap(column, value)
			cells = append(cells, resPath.(string))

		}

		result = append(result, cells)

	}

	lib.CheckError(err)
	return result
}
func (track Track) saveFile(matrix [][]string, from time.Time, to time.Time, now time.Time) string {
	filepath := track.formatFilename(track.FileName, from, to, now)
	switch track.Type {
	case "csv":

		sep := []rune(track.CsvConfig.Separator)
		lib.WriteCsv("../tmp/"+filepath, matrix, sep[0])
		source, _ := os.ReadFile("../tmp/" + filepath)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/"+track.Name+"/"+strconv.Itoa(from.Year())+"/"+filepath, source)
	case "excel":

		_, e := lib.CreateExcel(matrix, filepath, "Risultato")
		lib.CheckError(e)
		source, _ := os.ReadFile("../tmp/" + filepath)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/"+track.Name+"/"+strconv.Itoa(from.Year())+"/"+filepath, source)

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
	//json.Unmarshal([]byte(getpolicymock()), &json_data)

	b, err := json.Marshal(policy)
	lib.CheckError(err)
	json.Unmarshal(b, &json_data)
	log.Println(string(b))
	for indexAsset, asset := range policy.Assets {
		for indexG, _ := range asset.Guarantees {

			for _, column := range event {
				var (
					resPath interface{}
					value   string
				)
				value = column.Value
				resPath = column.Value
				log.Println(column.Value)
				if strings.Contains(column.Value, "$.") {
					if strings.Contains(column.Value, "$.assets[*].guarantees[*]") {
						value = strings.Replace(value, "guarantees[*]", "guarantees["+strconv.Itoa(indexG)+"]", 1)
						value = strings.Replace(value, "assets[*]", "assets["+strconv.Itoa(indexAsset)+"]", 1)
					}
					log.Println(value)
					resPath, err = jsonpath.JsonPathLookup(json_data, value)
					lib.CheckError(err)
					log.Println(resPath)
				}
				resPath = checkMap(column, value)
				cells = append(cells, resPath.(string))

			}
		}
		result = append(result, cells)
	}
	return result

}
func (track Track) upload(filePath string) {
	switch track.UploadType {
	case "sftp":

		track.sftp(filePath)
	}

}
func checkMap(column Column, value string) interface{} {
	var res interface{}
	if column.MapFx != "" {
		res = GetMapFx(column.MapFx, value)
	}
	if column.MapStatic != nil {
		res = column.MapStatic[value]
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
func (track Track) query(from time.Time, to time.Time, queryEvent Database) []models.Transaction {
	firequery := lib.Firequeries{
		Queries: []lib.Firequery{

			{
				Field:      "policyName", //
				Operator:   "==",         //
				QueryValue: track.Name,
			},
			{
				Field:      "effectiveDate", //
				Operator:   ">=",            //
				QueryValue: from,
			},
			{
				Field:      "effectiveDate", //
				Operator:   "<=",            //
				QueryValue: to,
			},
		},
	}
	for _, qe := range queryEvent.Query {

		firequery.Queries = append(firequery.Queries,
			lib.Firequery{
				Field:      qe.Field,    //
				Operator:   qe.Operator, //
				QueryValue: qe.QueryValue,
			})
	}
	query, e := firequery.FirestoreWherefields(queryEvent.Dataset)
	lib.CheckError(e)
	transactions := TransactionToListData(query)

	return transactions
}
func query[T any](from time.Time, to time.Time, db Database) []T {
	var res []T
	switch db.Name {

	case "firestore":
		res = firestoreQuery[T](from, to, db)

	}

	return res
}
func firestoreQuery[T any](from time.Time, to time.Time, queryEvent Database) []T {

	firequery := lib.FireGenericQueries[T]{
		Queries: []lib.Firequery{
			{
				Field:      "effectiveDate", //
				Operator:   ">=",            //
				QueryValue: from,
			},
			{
				Field:      "effectiveDate", //
				Operator:   "<=",            //
				QueryValue: to,
			},
		},
	}
	for _, qe := range queryEvent.Query {

		firequery.Queries = append(firequery.Queries,
			lib.Firequery{
				Field:      qe.Field,    //
				Operator:   qe.Operator, //
				QueryValue: qe.QueryValue,
			})
	}
	res, e := firequery.FireQuery(queryEvent.Dataset)
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
