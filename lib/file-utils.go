package lib

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	fireStorage "firebase.google.com/go"
	firebase "firebase.google.com/go"
)

func Files(path string) []string {
	var res []string
	if path == "" {
		path = "./"
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		res = append(res, f.Name())
	}
	return res
}
func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
func ReadDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
func GetFromStorage(bucket string, file string, keyPath string) []byte {
	//var credential models.Credential
	log.Println("start GetFromStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	CheckError(err)
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	CheckError(err)
	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	CheckError(err)
	return slurp
}
func GetReaderGCS(bucket string, file string, keyPath string) io.Reader {
	//var credential models.Credential
	log.Println("start GetFromStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	CheckError(err)
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	CheckError(err)
	slurp := rc
	rc.Close()
	CheckError(err)
	return slurp
}
func deleteFiles() {

}
func PutToStorage(bucketname string, path string, file []byte) string {
	// some process request msg, decode base64 to image byte
	// create image file in current directory with os.create()
	log.Println("start PutToStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	CheckError(err)
	bucket := client.Bucket(bucketname)
	write := bucket.Object(path).NewWriter(ctx)

	defer write.Close()
	write.Write(file)
	log.Println("write.MediaLink: " + write.MediaLink)

	CheckError(err)
	log.Println("write.MediaLink: ")
	return "gs://" + bucketname + "/" + path

}
func PutToFireStorage(bucketname string, path string, file []byte) string {
	// some process request msg, decode base64 to image byte
	// create image file in current directory with os.create()
	log.Println("start PutToStorage")
	ctx := context.Background()
	config := &firebase.Config{
		StorageBucket: "positive-apex-350507.appspot.com",
	}
	app, err := fireStorage.NewApp(ctx, config, nil)

	CheckError(err)
	client, e := app.Storage(ctx)
	CheckError(e)

	bucket, e := client.DefaultBucket()
	wc := bucket.Object(path).NewWriter(ctx)

	wc.Write(file)

	CheckError(e)
	defer wc.Close()

	CheckError(err)
	log.Println("write.MediaLink: ")
	return "gs://positive-apex-350507.appspot.com/UID:0g10fVw0fdOM5Ugho1tQJcOcRVD3/image_profile/profileImage"

}
func GetFilesByEnv(file string) []byte {
	var res1 []byte
	switch os.Getenv("env") {

	case "local":
		res1 = ErrorByte(ioutil.ReadFile("function-data/" + file))

	case "dev":

		res1 = GetFromStorage("function-data", file, "")

	case "prod":
		res1 = GetFromStorage("core-350507-function-data", file, "")

	default:

	}
	return res1
}
func GetByteByEnv(file string, isLocal bool) []byte {
	var res1 []byte
	switch os.Getenv("env") {

	case "local":
		res1 = ErrorByte(ioutil.ReadFile("function-data/" + file))

	case "dev":
		if isLocal {
			res1 = ErrorByte(ioutil.ReadFile("./serverless_function_source_code/" + file))
		} else {
			res1 = GetFromStorage("function-data", file, "")
		}

	case "prod":

		if isLocal {
			res1 = ErrorByte(ioutil.ReadFile("./serverless_function_source_code/" + file))
		} else {
			res1 = GetFromStorage("core-350507-function-data", file, "")
		}
	default:

	}
	return res1
}
func GetAssetPathByEnv(base string) string {
	var res1 string
	switch os.Getenv("env") {

	case "local":
		res1 = base + "/assets"

	case "dev":
		res1 = "./serverless_function_source_code/assets"
	case "prod":
		res1 = "./serverless_function_source_code/assets"

	default:
	}
	log.Println(res1)
	return res1
}
