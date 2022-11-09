package lib

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

func GetFirestore(collection string, doc string) *firestore.DocumentSnapshot {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Get(ctx)
	CheckError(err)
	return docsnap
}
func PutFirestore(collection string, doc string, value interface{}) (*firestore.DocumentRef, *firestore.WriteResult) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	ref, result, err := c.Add(ctx, value)
	CheckError(err)
	log.Println(ref)
	return ref, result
	//fmt.Println(dataMap)
}
func SetFirestore(collection string, doc string, value interface{}) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Set(ctx, value)
	CheckError(err)
	log.Println(docsnap)

	//fmt.Println(dataMap)
}
