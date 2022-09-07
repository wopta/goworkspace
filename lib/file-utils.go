package lib

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func Files(path string) {
	if path == "" {
		path = "./"
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
func ReadDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
func GetFromStorage(bucket string, file string) []byte {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("SA_KEY"))))
	if err != nil {
		log.Fatal(err)
	}
	rc, err := client.Bucket(bucket).Object(file).NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	slurp, err := ioutil.ReadAll(rc)

	rc.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(slurp)
	return slurp
}
