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
type FireGenericQueries[T any] struct {
	Queries []Firequery
	result  []T
}

func getFireClient() *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	return client
}

func NewDoc(collection string) string {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	ref := client.Collection(collection).NewDoc()
	return ref.ID
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

	if err != nil {
		return nil, err
	}

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
	// fmt.Println(dataMap)
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
	// fmt.Println(dataMap)
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

	// fmt.Println(dataMap)firestore.Query
}

func SetFirestoreErr(collection string, doc string, value interface{}) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return err
	}
	c := client.Collection(collection)
	col := c.Doc(doc)
	_, err = col.Set(ctx, value)

	return err
}

func FireUpdate(collection string, doc string, value interface{}) (*firestore.WriteResult, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	docsnap, err := col.Set(ctx, value, firestore.MergeAll)

	return docsnap, err
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

func DeleteFirestoreErr(collection string, doc string) (*firestore.WriteResult, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	c := client.Collection(collection)
	col := c.Doc(doc)
	res, err := col.Delete(ctx)
	return res, err
}

func WhereFirestore(collection string, field string, operator string, queryValue string) *firestore.DocumentIterator {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	query := client.Collection(collection).Where(field, operator, queryValue).Documents(ctx)

	return query
}

func QueryWhereFirestore(collection string, field string, operator string, queryValue string) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	query := client.Collection(collection).Where(field, operator, queryValue).Documents(ctx)

	return query, err
}

func (queries *Firequeries) FirestoreWherefields(collection string) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	var query firestore.Query
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	col := client.Collection(collection)
	query = col.Where(queries.Queries[0].Field, queries.Queries[0].Operator, queries.Queries[0].QueryValue)
	for i := 1; i <= len(queries.Queries)-1; i++ {
		query = query.Where(queries.Queries[i].Field, queries.Queries[i].Operator, queries.Queries[i].QueryValue)
	}

	return query.Documents(ctx), err
}

/*func (queries *FireGenericQueries[T]) FireQuery(collection string) ([]T, error) {
	ctx := context.Background()
	var query firestore.Query
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	col := client.Collection(collection)
	query = col.Where(queries.Queries[0].Field, queries.Queries[0].Operator, queries.Queries[0].QueryValue)
	for i := 1; i <= len(queries.Queries)-1; i++ {
		query = query.Where(queries.Queries[i].Field, queries.Queries[i].Operator, queries.Queries[i].QueryValue)
	}
	q := query.Documents(ctx)
	result := make([]T, 0)
	for {
		d, err := q.Next()
		if err != nil {
		}
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value T
		e := d.DataTo(&value)
		CheckError(e)
		result = append(result, value)
		log.Println(len(result))
	}
	return result, err
}*/

func (queries *Firequeries) FirestoreWhereLimitFields(collection string, limit int) (*firestore.DocumentIterator, error) {
	ctx := context.Background()
	var query firestore.Query
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	col := client.Collection(collection)
	query = col.Where(queries.Queries[0].Field, queries.Queries[0].Operator, queries.Queries[0].QueryValue)
	for i := 1; i <= len(queries.Queries)-1; i++ {
		query = query.Where(queries.Queries[i].Field, queries.Queries[i].Operator, queries.Queries[i].QueryValue)
	}

	return query.Limit(limit).Documents(ctx), err
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

func FireToListData[T interface{}](query *firestore.DocumentIterator) []T {
	var result []T
	var _struct T
	for {
		d, err := query.Next()
		log.Println("for")
		if err != nil {
			log.Println("error")
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}

		} else {
		}
		e := d.DataTo(&_struct)

		log.Println("todata")
		CheckError(e)
		result = append(result, _struct)

		log.Println(len(result))
	}
	return result
}

func SetBatchFirestoreErr[T any](operations map[string]map[string]T) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return err
	}

	batch := client.Batch()
	for collection, values := range operations {
		c := client.Collection(collection)

		for k, v := range values {
			col := c.Doc(k)
			batch.Set(col, v)
		}
	}

	_, err = batch.Commit(ctx)

	return err
}
