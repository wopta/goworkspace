package lib

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func GetFirestore(collection string, doc string) *firestore.DocumentSnapshot {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	CheckError(err)
	c := client.Collection("States")
	col := c.Doc("NewYork")
	docsnap, err := col.Get(ctx)
	CheckError(err)
	return docsnap
}
func PutFirestore(collection string, doc string, value interface{}) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Update(ctx, []firestore.Update{{Value: value}})
	log.Println(docsnap)
	CheckError(err)
	//fmt.Println(dataMap)
}
