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
func PutFirestore(collection string, value interface{}) (*firestore.DocumentRef, *firestore.WriteResult) {
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

	//fmt.Println(dataMap)firestore.Query
}
func WhereFirestore(collection string, field string, operator string, queryValue string) *firestore.DocumentIterator {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).Where(field, operator, queryValue).Documents(ctx)
	query.GetAll()

	return query
}
func OrderFirestore(collection string, field string, value firestore.Direction) *firestore.DocumentIterator {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).OrderBy(field, value).Documents(ctx)

	return query
}
func WhereLimitFirestore(collection string, field string, operator string, queryValue string, limit int) *firestore.DocumentIterator {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).Where(field, operator, queryValue).LimitToLast(limit).Documents(ctx)
	query.GetAll()

	return query
}
func OrderLimitFirestore(collection string, field string, value firestore.Direction, limit int) *firestore.DocumentIterator {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).OrderBy(field, value).Limit(limit).Documents(ctx)

	return query
}
func OrderLimitFirestoreErr(collection string, field string, value firestore.Direction, limit int) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	query := client.Collection(collection).OrderBy(field, value).Limit(limit).Documents(ctx)

	return query, err
}
func OrderWhereLimitFirestoreErr(collection string, field string, fieldOrder string, operator string, queryValue string, value firestore.Direction, limit int) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	query := client.Collection(collection).Where(field, operator, queryValue).OrderBy(fieldOrder, value).Limit(limit).Documents(ctx)

	return query, err
}
