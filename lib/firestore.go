package lib

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Firequery struct {
	Field      string
	Operator   string
	QueryValue interface{}
}
type Firequeries struct {
	Queries []Firequery
}

func getFireClient() *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	return client
}
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
func GetFirestoreErr(collection string, doc string) (*firestore.DocumentSnapshot, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Get(ctx)
	return docsnap, err
}
func GetFirestoreData(collection string, doc string, i *interface{}) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Get(ctx)
	CheckError(err)
	docsnap.DataTo(i)

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
func PutFirestoreErr(collection string, value interface{}) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	ref, result, err := c.Add(ctx, value)
	CheckError(err)
	log.Println(ref)
	return ref, result, err
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
func UpdateFirestoreErr(collection string, doc string, values map[string]interface{}) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)

	var updateValues []firestore.Update

	for key, val := range values {
		updateValues = append(updateValues, firestore.Update{Path: key, Value: val})
	}

	docsnap, err := col.Update(ctx, updateValues)
	log.Println(docsnap)

	return err
}
func WhereFirestore(collection string, field string, operator string, queryValue string) *firestore.DocumentIterator {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).Where(field, operator, queryValue).Documents(ctx)
	query.GetAll()

	return query
}
func QueryWhereFirestore(collection string, field string, operator string, queryValue string) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	query := client.Collection(collection).Where(field, operator, queryValue).Documents(ctx)

	return query, err
}
func (queries *Firequeries) FirestoreWherefields(collection string) *firestore.DocumentIterator {
	ctx := context.Background()
	var query firestore.Query
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	col := client.Collection(collection)
	query = col.Where(queries.Queries[0].Field, queries.Queries[0].Operator, queries.Queries[0].QueryValue)
	for i := 1; i < len(queries.Queries)-1; i++ {
		query = col.Where(queries.Queries[i].Field, queries.Queries[i].Operator, queries.Queries[i].QueryValue)
	}

	return query.Documents(ctx)
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
	query := client.Collection(collection).Where(field, operator, queryValue).Limit(limit).Documents(ctx)

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
func ToListData(query *firestore.DocumentIterator, v interface{}) []interface{} {
	var result []interface{}
	for {
		d, err := query.Next()
		log.Println("for")
		if err != nil {
			log.Println("error")
		}
		if err != nil {
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}

		}
		e := d.DataTo(&v)

		log.Println("todata")
		CheckError(e)
		result = append(result, v)

		log.Println(len(result))
	}
	return result
}
