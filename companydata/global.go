package companydata

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var (
	config = lib.SftpConfig{
		Username:     os.Getenv("GLOBAL_SFTP_USER"),
		Password:     os.Getenv("GLOBAL_SFTP_PSW"), // required only if password authentication is to be used
		PrivateKey:   "",                           // required only if private key authentication is to be used
		Server:       "ftps.globalassistance.it:222",
		KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
		Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
	}
	location, e   = time.LoadLocation("Europe/Rome")
	executiondate = time.Now().In(location)
	from          = time.Date(executiondate.Year(), executiondate.Month(), executiondate.Day(), 0, 0, 0, 0, location)
	to            = time.Date(executiondate.Year(), executiondate.Month(), executiondate.Day(), 8, 0, 0, 0, location)
)

func GlobalSftpDownload(filename string, bucket string, folder string) ([]byte, io.ReadCloser, error) {

	if executiondate.After(from) && executiondate.Before(to) {

		localPath := "../tmp/" + filename
		client, e := lib.NewSftpclient(config)
		client.ListFiles(".")
		println("folder +filename: ", folder+filename)
		println("GlobalSftpDownload error: ", fmt.Errorf("unable to open remote file: %v", e))
		srcFile, e := client.Download(folder + filename)
		if e != nil {
			log.Println(fmt.Errorf("unable to open remote file: %v", e))
		}
		defer srcFile.Close()

		dstFile, err := os.Create(localPath)
		if err != nil {
			log.Println(fmt.Errorf("unable to open local file: %v", err))
		}
		defer dstFile.Close()

		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			log.Println(fmt.Errorf("unable to download remote file: %v", err))
		}
		log.Printf("%d bytes copied to %v", bytes, localPath)

		//sourceByte, e := io.ReadAll(srcFile)
		log.Println(e)
		//buf := new(bytes.Buffer)
		//_, e = buf.ReadFrom(reader)
		sourceByte, _ := ioutil.ReadFile("../tmp/" + filename)
		//excelsource, _ := lib.ExcelRead(reader)
		lib.PutToStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "track/in/global/emit/"+filename, sourceByte)
		return sourceByte, srcFile, e
	}
	return nil, nil, e
}
func GlobalSftpUpload(filename string, folder string) error {

	if executiondate.After(from) && executiondate.Before(to) && os.Getenv("env") == "prod" {

		client, e := lib.NewSftpclient(config)
		println("filename: ", filename)
		defer client.Close()
		log.Println("Open local file for reading.:")
		source, e := os.Open("../tmp/" + filename)
		lib.CheckError(e)
		//defer source.Close()
		log.Println("Create remote file for writing:")
		// Create remote file for writing.
		lib.Files("../tmp")
		//destination, e := client.Create(filePath)
		destination, e := client.Create(folder + filename)

		defer destination.Close()
		log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
		info, e := source.Stat()
		log.Println(info.Size())
		// Upload local file to a remote location as in 1MB (byte) chunks.
		e = client.Upload(source, destination, int(info.Size()))
	}
	return e
}
func GlobalSftpDelete(filename string) error {
	if executiondate.After(from) && executiondate.Before(to) && os.Getenv("env") == "prod" {
		client, e := lib.NewSftpclient(config)

		println("filename: ", filename)
		defer client.Close()
		log.Println("Open local file for reading.:")
		e = client.Remove(filename)
		defer client.Close()
		return e
	}
	return e
}

func getInstallamentDate(p models.Policy, layout string) string {
	var res string
	res = p.EndDate.Format(layout)
	if p.PaymentSplit == "monthly" {
		res = p.StartDate.AddDate(0, 1, 0).Format(layout)
	}

	return res
}
func getInstallament(key string, price float64) float64 {
	var res float64
	res = price
	if key == "monthly" {
		res = price / 12
	}
	return res
}
func getYesNo(key bool) string {
	var res string
	mapGarante := map[bool]string{true: "SI", false: "NO"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getOneTwo(key bool) string {
	var res string
	mapGarante := map[bool]string{true: "1", false: "2"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingYear(key string) string {
	var res string
	mapGarante := map[string]string{"before1972": "1", "1972between2009": "2", "after2009": "3"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingFloor(key string) string {
	var res string
	mapGarante := map[string]string{"ground_floor": "1", "first": "2", "second": "3", "greater_than_second": "4"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapBuildingMaterial(key string) string {
	var res string
	mapGarante := map[string]string{"masonry": "1", "reinforcedConcrete": "2", "antiSeismicLaminatedTimber": "3", "steel": "4"}

	if seconds, ok := mapGarante[key]; ok { // will be false if person is not in the map
		res = seconds
	}
	return res
}
func getMapSplit(key string) string {
	var res string
	res = "1"
	if key == "monthly" {
		res = "12"
	}
	return res
}
func getBuildingType(key string) string {
	var res string
	res = "1"
	if key == "montly" {
		res = "12"
	}
	return res
}
func getMapRevenue(key int) string {
	var res int

	if key <= 200000 { // will be false if person is not in the map
		res = 1
	}
	if key > 200000 && key <= 500000 { // will be false if person is not in the map
		res = 2
	}
	if key > 500000 && key <= 1000000 { // will be false if person is not in the map
		res = 3
	}
	if key > 1000000 && key <= 1500000 { // will be false if person is not in the map
		res = 4
	}
	if key > 1500000 && key <= 5000000 { // will be false if person is not in the map
		res = 5
	}
	if key > 5000000 && key <= 7500000 { // will be false if person is not in the map
		res = 6
	}
	if key > 7500000 && key <= 10000000 { // will be false if person is not in the map
		res = 7
	}
	return strconv.Itoa(res)
}
func getMapSelfInsurance(key string) string {
	var res int

	if key == "5% - minimo € 500" { // will be false if person is not in the map
		res = 1
	}
	if key == "5% - minimo € 1.000" { // will be false if person is not in the map
		res = 2
	}
	if key == "5% - minimo € 1.500" { // will be false if person is not in the map
		res = 3
	}
	if key == "10% - minimo € 500" { // will be false if person is not in the map
		res = 4
	}
	if key == "10% - minimo € 1.000" { // will be false if person is not in the map
		res = 5
	}
	if key == "10% - minimo € 1.500" { // will be false if person is not in the map
		res = 6
	}
	if key == "10% - minimo € 2.000" { // will be false if person is not in the map
		res = 7
	}
	if key == "10% - minimo € 3.000" { // will be false if person is not in the map
		res = 8
	}
	if key == "10% - minimo € 5.000" { // will be false if person is not in the map
		res = 9
	}
	if key == "15% - minimo € 5.000" { // will be false if person is not in the map
		res = 10
	}
	if key == "10% - minimo € 10.000" { // will be false if person is not in the map
		res = 11
	}
	if key == "10% - minimo € 20.000" { // will be false if person is not in the map
		res = 12
	}
	if key == "10% - minimo € 25.000" { // will be false if person is not in the map
		res = 13
	}
	if key == "10% - minimo € 30.000" { // will be false if person is not in the map
		res = 14
	}

	return strconv.Itoa(res)
}
func getSumLimit(sumlimitContentBuilding float64, g models.Guarante) (string, string) {
	var (
		sum, perc string
	)
	sum = strconv.Itoa(int(g.SumInsuredLimitOfIndemnity))

	if g.SumInsuredLimitOfIndemnity <= 1 {
		percF := g.SumInsuredLimitOfIndemnity * 100
		sumF := sumlimitContentBuilding * g.SumInsuredLimitOfIndemnity
		perc = strconv.Itoa(int(percF))
		sum = strconv.Itoa(int(sumF))
	}

	return sum, perc
}
