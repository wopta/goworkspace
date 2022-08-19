package lib

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

//import "github.com/go-gota/gota/series"

func CsvToDataframe(filePath string) dataframe.DataFrame {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	df := dataframe.ReadCSV(f,
		dataframe.WithDelimiter(';'),
		dataframe.HasHeader(true))
	return df
}
