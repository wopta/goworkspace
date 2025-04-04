package inclusive

import (
	"context"

	"log"
	"os"

	"cloud.google.com/go/bigquery"
	lib "github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func QueryRowsBigQuery[T any](datasetID string, tableID string, query string) ([]T, error) {
	var (
		res  []T
		e    error
		iter *bigquery.RowIterator
	)
	log.Println(query)
	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	queryi := client.Query(query)
	iter, e = queryi.Read(ctx)
	log.Println(e)
	for {
		var row T
		e := iter.Next(&row)

		if e == iterator.Done {
			log.Println(e)
			return res, e
		}
		if e != nil {
			log.Println(e)
			return res, e
		}

		res = append(res, row)

	}

}

func getBigqueryClient() *bigquery.Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	lib.CheckError(err)
	return client
}

func InsertRowsBigQuery(datasetID string, tableID string, value interface{}) error {
	client := getBigqueryClient()
	defer client.Close()
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	e := inserter.Put(context.Background(), value)
	log.Println(e)
	return e
}
