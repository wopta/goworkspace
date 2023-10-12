package lib

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

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
func GetFromGoogleStorage(bucket string, file string) ([]byte, error) {
	//var credential models.Credential
	log.Println("GetFromGoogleStorage")
	client, ctx, err := GetGoogleStorageClient()
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
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
	log.Printf("[ReadFileFromGoogleStorage] gsLink: %s", gsLink)
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
	log.Printf("[ReadFileFromGoogleStorage] invalid gsLink: %s", gsLink)
	return nil, errors.New("invalid gsLink")
}

func ListGoogleStorageFolderContent(folderPath string) ([]string, error) {
	filesList := make([]string, 0)
	bucket := os.Getenv("GOOGLE_STORAGE_BUCKET")

	log.Println("[ListGoogleStorageFolderContent] function start ---------------")

	log.Printf("[ListGoogleStorageFolderContent] bucket: %s ---- folderPath: %s", bucket, folderPath)

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

	log.Printf("[ListGoogleStorageFolderContent] found %d files", len(filesList))

	return filesList, nil
}
