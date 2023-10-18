package lib

import (
	"context"
	"encoding/csv"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
		res = append(res, strconv.FormatInt(f.Size(), 10))

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

func ReadLocalDirContent(folderPath string) [][]byte {
	var (
		res [][]byte
	)
	dir, err := os.ReadDir(folderPath)
	CheckError(err)

	for _, contentDir := range dir {
		if contentDir.IsDir() {
			res = append(res, ReadLocalDirContent(folderPath+"/"+contentDir.Name())...)
			continue
		}
		fileByte := ErrorByte(os.ReadFile(folderPath + "/" + contentDir.Name()))
		res = append(res, fileByte)
	}

	return res
}

func GetFromStorage(bucket string, file string, keyPath string) []byte {
	//var credential models.Credential
	log.Println("start GetFromStorage")
	log.Println("File: " + file)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	CheckError(err)
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	CheckError(err)
	slurp, err := io.ReadAll(rc)
	rc.Close()
	CheckError(err)
	return slurp
}

func GetFromStorageErr(bucket string, file string, keyPath string) ([]byte, error) {
	//var credential models.Credential
	log.Println("start GetFromStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	return slurp, err
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

func ReadStorageDirContent(bucketName, folderPath string) [][]byte {
	var (
		res [][]byte
	)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{
		Prefix:    folderPath,
		Delimiter: "/",
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		// this checks if blob correspond to a directory
		if attrs.ContentType == "" {
			continue
		}
		file := GetFromStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), attrs.Name, "")
		res = append(res, file)
	}

	return res
}

func PutToStorage(bucketname string, path string, file []byte) string {

	log.Println("start PutToStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	bucket := client.Bucket(bucketname)
	write := bucket.Object(path).NewWriter(ctx)
	defer write.Close()
	write.Write(file)
	CheckError(err)
	return "gs://" + bucketname + "/" + path

}

func PutGoogleStorage(bucketname string, path string, file []byte, contentType string) (string, error) {
	// some process request msg, decode base64 to image byte
	// create image file in current directory with os.create()
	log.Println("start PutToStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	bucket := client.Bucket(bucketname)
	write := bucket.Object(path).NewWriter(ctx)
	write.ContentType = contentType
	write.Close()
	objAttrs, err := client.Bucket(bucketname).Object(path).Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: contentType,
	})
	fmt.Println(objAttrs)
	CheckError(err)
	return "gs://" + bucketname + "/" + path, err

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
		res1 = ErrorByte(os.ReadFile("../function-data/dev/" + file))
	case "dev":
		res1 = GetFromStorage("function-data", file, "")
	case "prod":
		res1 = GetFromStorage("core-350507-function-data", file, "")

	default:

	}
	return res1
}

func GetFolderContentByEnv(folderName string) [][]byte {
	var res [][]byte

	switch os.Getenv("env") {
	case "local":
		res = ReadLocalDirContent("../function-data/dev/" + folderName)
	case "dev":
		res = ReadStorageDirContent("function-data", folderName)
	case "prod":
		res = ReadStorageDirContent("core-350507-function-data", folderName)
	}

	return res
}

func GetByteByEnv(file string, isLocal bool) []byte {
	var res1 []byte
	switch os.Getenv("env") {

	case "local":
		res1 = ErrorByte(os.ReadFile("../function-data/dev/" + file))
	case "dev":
		if isLocal {
			res1 = ErrorByte(os.ReadFile("./serverless_function_source_code/" + file))
		} else {
			res1 = GetFromStorage("function-data", file, "")
		}
	case "prod":
		if isLocal {
			res1 = ErrorByte(os.ReadFile("./serverless_function_source_code/" + file))
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

	return res1
}

func GetAssetPathByEnvV2() string {
	var path string
	switch os.Getenv("env") {
	case "local":
		path = "../function-data/dev/assets/documents/"
	case "dev", "prod":
		path = "./serverless_function_source_code/../tmp/assets/"
	}

	return path
}
