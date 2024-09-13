package companydata

import (
	"encoding/json"
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
		query        []Query
	)
	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	now, upload, reqData := getCompanyDataReq(req)

	procuctTrackByte := lib.GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+reqData.Name+"/v1/track.json", "")
	json.Unmarshal(procuctTrackByte, &procuctTrack)

	switch reqData.Event {

	case "emit":
		query = procuctTrack.Emit.Query
		event = procuctTrack.Emit.Event
	case "payment":
		query = procuctTrack.Payment.Query
		event = procuctTrack.Payment.Event
	case "delete":
		query = procuctTrack.Delete.Query
		event = procuctTrack.Delete.Event

	}

	from, to = procuctTrack.frequency(now)
	for _, transaction := range procuctTrack.query(from, to, query) {
		var (
			policy *models.Policy
		)

		docsnap := lib.GetFirestore("policy", transaction.PolicyUid)
		docsnap.DataTo(&policy)
		result = append(result, procuctTrack.ProductTrack(policy, event)...)

	}
	filepath := procuctTrack.saveFile(result, from, to, now)
	if upload {
		procuctTrack.upload(filepath)
	}
	return "", nil, e
}
func (track Track) ProductTrack(policy *models.Policy, event []Column) [][]string {
	var (
		result [][]string

		err error
	)
	if track.IsAssetFlat {
		result = track.assetRow(policy, event)
	} else {
		result = track.assetRow(policy, event)
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
func (track Track) assetRow(policy *models.Policy, event []Column) [][]string {
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

				if column.MapFx != "" {
					resPath = GetMapFx(column.MapFx, column.Value)
				}
				if column.MapStatic != nil {
					resPath = column.MapStatic[column.Value]
				}
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
func (track Track) query(from time.Time, to time.Time, queryEvent []Query) []models.Transaction {
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
	for _, qe := range queryEvent {

		firequery.Queries = append(firequery.Queries,
			lib.Firequery{
				Field:      qe.Field,    //
				Operator:   qe.Operator, //
				QueryValue: qe.QueryValue,
			})
	}
	query, e := firequery.FirestoreWherefields("transactions")
	lib.CheckError(e)
	transactions := TransactionToListData(query)
	return transactions
}

func (track Track) formatFilename(filename string, from time.Time, to time.Time, now time.Time) string {
	filename = strings.Replace(filename, "fdd", string(from.Day()), 1)
	filename = strings.Replace(filename, "fmm", string(from.Month()), 1)
	filename = strings.Replace(filename, "fyyyy", string(from.Year()), 1)

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
