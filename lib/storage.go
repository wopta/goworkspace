package lib

import (
	"context"
	"io/ioutil"
	"log"

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
