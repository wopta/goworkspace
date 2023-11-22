package lib

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
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

func QueryRowsBigQuery[T any](query string) ([]T, error) {
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

	for {
		var row T
		e = iter.Next(&row)

		if e == iterator.Done {
			log.Println(e)
			return res, nil
		}
		if e != nil {
			log.Println(e)
			return res, e
		}

		res = append(res, row)

	}

}

func QueryParametrizedRowsBigQuery[T any](query string, params map[string]interface{}) ([]T, error) {
	var (
		res  []T
		e    error
		iter *bigquery.RowIterator
	)
	client := getBigqueryClient()
	ctx := context.Background()
	defer client.Close()
	queryBigQuery := client.Query(query)

	for name, value := range params {
		queryBigQuery.Parameters = append(queryBigQuery.Parameters, bigquery.QueryParameter{Name: name, Value: value})
	}

	iter, e = queryBigQuery.Read(ctx)

	if e != nil {
		return res, e
	}
	for {
		var row T
		e := iter.Next(&row)
		if e == iterator.Done {
			return res, e
		}
		if e != nil {
			return res, e
		}
		res = append(res, row)

	}

}

func UpdateRowBigQuery(datasetID string, tableID string, params map[string]string, condiction string) error {
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
	count := 1
	for k, v := range params {
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString("'" + v + "'")
		if len(params) > count {
			b.WriteString(", ")
		}
		count = count + 1

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

// use biquery query constructors
func UpdateRowBigQueryV2(datasetId, tableId string, params map[string]interface{}, condition string) error {
	var (
		bytes bytes.Buffer
		err   error
	)

	if len(params) == 0 {
		return fmt.Errorf("no params to update")
	}

	tableSelection := fmt.Sprintf("UPDATE `%s.%s` SET ", datasetId, tableId)
	bytes.WriteString(tableSelection)

	count := 1
	length := len(params)
	for key, value := range params {
		var query string
		if reflect.TypeOf(value).String() == "string" {
			query = key + "=" + "'" + reflect.ValueOf(value).String() + "'"
		} else if reflect.TypeOf(value).String() == "bool" {
			query = key + "=" + strconv.FormatBool(value.(bool))
		} else if reflect.TypeOf(value).String() == "[]string" {
			query = key + "=" + "ARRAY ["
			for index, s := range value.([]string) {
				query = query + "'" + s + "'"
				if index < len(value.([]string))-1 {
					query = query + ", "
				}
			}
			query = query + "]"
		} else if reflect.TypeOf(value).String() == "bigquery.NullDateTime" {
			query = key + "=" + "'" + bigquery.CivilDateTimeString(value.(bigquery.NullDateTime).DateTime) + "'"
		} else if reflect.TypeOf(value).String() == "float64" {
			query = key + "=" + strconv.FormatFloat(value.(float64), 'f', 2, 64)
		}
		if count < length {
			query = query + ", "
		} else {
			query = query + " "
		}
		bytes.WriteString(query)
		count = count + 1
	}
	bytes.WriteString(condition)
	queryString := bytes.String()
	log.Printf("[UpdateRowBigQueryV2] query: %s", queryString)

	client := getBigqueryClient()
	defer client.Close()
	ctx := context.Background()
	query := client.Query(queryString)
	job, err := query.Run(ctx)
	if err != nil {
		log.Printf("[UpdateRowBigQueryV2] error running query: %s", err.Error())
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		log.Printf("[UpdateRowBigQueryV2] error waiting for job to run: %s", err.Error())
		return err
	}
	if err := status.Err(); err != nil {
		log.Printf("[UpdateRowBigQueryV2] error on job: %s", err.Error())
		return err
	}

	return err
}

func GetBigQueryNullDateTime(date time.Time) bigquery.NullDateTime {
	nilTime := time.Time{}
	return bigquery.NullDateTime{
		DateTime: civil.DateTimeOf(date),
		Valid:    date != nilTime,
	}
}

func GetBigQueryNullGeography(latitude, longitude float64) bigquery.NullGeography {
	// TODO: Check if correct: Geography type uses the WKT format for geometry
	return bigquery.NullGeography{
		GeographyVal: fmt.Sprintf("POINT (%f %f)", longitude, latitude),
		Valid:        true,
	}
}
