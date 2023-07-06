package lib

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
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
	log.Println(e)
	return e
}
func QueryRowsBigQuery[T any](datasetID string, tableID string, query string) ([]T, error) {
	var (
		res  []T
		e    error
		iter *bigquery.RowIterator
	)
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
			return res, e
		}
		if e != nil {
			return res, e
		}

	}

}
func UpdateRowBigQuery(datasetID string, tableID string, id string, params map[string]string, condiction string) error {
	var (
		e error
		b bytes.Buffer
	)

	b.WriteString("UPDATE")
	b.WriteString(" ")
	b.WriteString(datasetID + "." + tableID)
	b.WriteString(" ")
	b.WriteString("SET")
	b.WriteString(" ")
	for k, v := range params {
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString("'" + v + "'")
		b.WriteString(" ")

	}
	b.WriteString("WHERE")
	b.WriteString(" ")
	b.WriteString(condiction)

	log.Println(b.String())

	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	q := client.Query(b.String())
	job, err := q.Run(ctx)
	status, err := job.Wait(ctx)
	if err != nil {
		log.Println(e)
		return err
	}

	if err := status.Err(); err != nil {
		log.Println(e)
		return err
	}
	return e
}

func GetBigQueryNullDateTime(date time.Time) bigquery.NullDateTime {
	nilTime := time.Time{}
	return bigquery.NullDateTime{
		DateTime: civil.DateTimeOf(date),
		Valid:    date != nilTime,
	}
}
