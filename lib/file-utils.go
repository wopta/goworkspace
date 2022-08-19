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
func readDir() {
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
func GetFromStorage() {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := client.Bucket("cloud-build-data").Object("data-functions/Riclassificazione%20_Ateco.csv").NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	slurp, err := ioutil.ReadAll(rc)

	rc.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(slurp)

}
