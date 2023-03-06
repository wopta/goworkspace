package lib

import (
	"context"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func getBigqueryClient() *bigquery.Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	CheckError(err)
	return client
}
func InsertRowsBigQuery(datasetID string, tableID string, value interface{}) error {
	client := getBigqueryClient()
	defer client.Close()
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	e := inserter.Put(context.Background(), value)

	return e
}
func QueryRowsBigQuery[T any](datasetID string, tableID string, value any) error {
	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	query := client.Query(
		`SELECT
                CONCAT(
                        'https://stackoverflow.com/questions/',
                        CAST(id as STRING)) as url,
                view_count
        FROM ` + "`bigquery-public-data.stackoverflow.posts_questions`" + `
        WHERE tags like '%google-bigquery%'
        ORDER BY view_count DESC
        LIMIT 10;`)
	iter, e := query.Read(ctx)
	for {
		var row T
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}

	}
	return e
}
