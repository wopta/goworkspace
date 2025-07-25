package inclusive

import (
	"context"

	"os"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"cloud.google.com/go/bigquery"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"google.golang.org/api/iterator"
)

func QueryRowsBigQuery[T any](datasetID string, tableID string, query string) []T {
	var (
		res  []T
		iter *bigquery.RowIterator
	)
	log.Println(query)
	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	queryi := client.Query(query)
	iter, _ = queryi.Read(ctx)
	for {
		var row T
		e := iter.Next(&row)

		if e == iterator.Done {
			return res
		}
		if e != nil {
			return res
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
	if e != nil {
		log.Error(e)
	}
	return e
}
