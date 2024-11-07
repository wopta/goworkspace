package lib

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
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

func GetFromStorageV2(bucket string, file string, keyPath string) ([]byte, error) {
	var err error
	log.Println("start GetFromStorageV2")
	log.Println("File: " + file)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	slurp, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	err = rc.Close()
	if err != nil {
		return nil, err
	}
	return slurp, err
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

func PutToStorageIfNotExists(bucketname string, path string, file []byte) (string, error) {
	log.Println("start PutToStorage")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucket := client.Bucket(bucketname)
	obj := bucket.Object(path)

	// Check if the object already exists
	_, err = obj.Attrs(ctx)
	if err == nil {
		// Object already exists, return an error
		return "", fmt.Errorf("file already exists")
	}
	// check if the error is because the object does not exist
	if err != storage.ErrObjectNotExist {
		return "", err
	}

	// Object does not exist, create a new writer and writer the file
	writer := obj.NewWriter(ctx)
	defer writer.Close()
	if _, err := writer.Write(file); err != nil {
		return "", err
	}

	return "gs://" + bucketname + "/" + path, nil
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
	app, err := firebase.NewApp(ctx, config, nil)

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
	case "local-test":
		res1 = ErrorByte(os.ReadFile("../../function-data/dev/" + file))
	case "dev":
		res1 = GetFromStorage("function-data", file, "")
	case "prod":
		res1 = GetFromStorage("core-350507-function-data", file, "")
	}

	return res1
}

func GetFilesByEnvV2(file string) ([]byte, error) {
	var (
		res []byte
		err error
	)

	switch os.Getenv("env") {
	case "local":
		res, err = os.ReadFile("../function-data/dev/" + file)
	case "local-test":
		res, err = os.ReadFile("../../function-data/dev/" + file)
	case "dev":
		res, err = GetFromStorageV2("function-data", file, "")
	case "prod":
		res, err = GetFromStorageV2("core-350507-function-data", file, "")
	}

	return res, err
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
		path = "./serverless_function_source_code/tmp/assets/"
	}

	return path
}

func CheckFileExistence(filePath string) bool {
	if IsLocal() {
		_, err := os.OpenFile("../function-data/dev/"+filePath, os.O_RDWR, 0755)
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return true
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	_, err = client.Bucket(os.Getenv("GOOGLE_STORAGE_BUCKET")).Object(filePath).Attrs(ctx)
	if err != nil {
		return false
	}
	return true
}

func ListLocalFolderContent(folderPath string) ([]string, error) {
	var (
		res      []string
		basePath = "../function-data/dev/"
	)

	err := filepath.WalkDir(basePath+folderPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		res = append(res, strings.ReplaceAll(path, basePath, ""))
		return nil
	})

	return res, err
}
