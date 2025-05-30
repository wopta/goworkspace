package lib

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func GetGoogleStorageClient() (*storage.Client, context.Context, error) {

	log.Println("GetGoogleStorageClient")
	ctx := context.Background()
	client, e := storage.NewClient(ctx)
	return client, ctx, e

}
func PutToGoogleStorage(bucketname string, path string, file []byte) (string, error) {
	log.Println("PutToGoogleStorage")
	client, ctx, e := GetGoogleStorageClient()
	bucket := client.Bucket(bucketname)
	write := bucket.Object(path).NewWriter(ctx)
	defer write.Close()
	write.Write(file)
	return "gs://" + bucketname + "/" + path, e

}

func PutToGoogleStorageWithSpecificContentType(bucketName string, path string, file []byte, contentType string) (str string, err error) {
	client, ctx, err := GetGoogleStorageClient()
	if err != nil {
		return "", fmt.Errorf("unable to get google storage client: %v", err)
	}
	bucket := client.Bucket(bucketName)
	write := bucket.Object(path).NewWriter(ctx)
	defer func() {
		err = write.Close()
	}()
	write.ContentType = contentType
	_, _ = write.Write(file)
	return "gs://" + bucketName + "/" + path, err
}

func GetFromGoogleStorage(bucket string, file string) ([]byte, error) {
	//var credential models.Credential
	log.Println("GetFromGoogleStorage")
	client, ctx, err := GetGoogleStorageClient()
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	slurp, err := ioutil.ReadAll(rc)

	return slurp, err
}

func PromoteFile(fromPath, toPath string) (string, error) {
	fileBytes, err := GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), fromPath)
	if err != nil {
		return "", err
	}

	return PutToGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), toPath, fileBytes)
}

func ReadFileFromGoogleStorage(gsLink string) ([]byte, error) {
	log.AddPrefix("ReadFileFromGoogleStorage")
	defer log.PopPrefix()
	log.Printf("gsLink: %s", gsLink)
	components := strings.Split(gsLink, "://")

	// Check if there are two components (scheme and path)
	if len(components) == 2 {
		scheme := components[0]
		path := components[1]

		// Split the path into bucket and object name
		pathComponents := strings.SplitN(path, "/", 2)
		if len(pathComponents) == 2 {
			bucketName := pathComponents[0]
			objectName := pathComponents[1]

			log.Println("Scheme:", scheme)
			log.Println("Bucket Name:", bucketName)
			log.Println("Object Name:", objectName)

			return GetFromStorageErr(bucketName, objectName, "")
		}
	}
	log.Printf("invalid gsLink: %s", gsLink)
	return nil, errors.New("invalid gsLink")
}

// Get file from bucket, either with a GSlink or a path
func ReadFileFromGoogleStorageEitherGsOrNot(pathOrGs string) ([]byte, error) {
	if strings.HasPrefix(pathOrGs, "gs://") {
		return ReadFileFromGoogleStorage(pathOrGs)
	} else {
		return GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), pathOrGs)
	}
}

func ListGoogleStorageFolderContent(folderPath string) ([]string, error) {
	filesList := make([]string, 0)
	bucket := os.Getenv("GOOGLE_STORAGE_BUCKET")
	log.AddPrefix("ListGoogleStorageFolderContent")
	defer log.PopPrefix()
	log.Println("function start ---------------")

	log.Printf("bucket: %s ---- folderPath: %s", bucket, folderPath)

	client, ctx, err := GetGoogleStorageClient()
	if err != nil {
		return filesList, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := &storage.Query{
		Prefix: folderPath,
	}

	it := client.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return filesList, fmt.Errorf("Bucket(%q).Objects: %w", bucket, err)
		}
		filesList = append(filesList, attrs.Name)
	}

	log.Printf("found %d files", len(filesList))

	log.Println("function end --------------")

	return filesList, nil
}
