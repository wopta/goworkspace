package lib

import (
	"bytes"

	"github.com/go-gota/gota/dataframe"
)

//import "github.com/go-gota/gota/series"

func CsvToDataframe(file []byte) dataframe.DataFrame {
	reader := bytes.NewReader(file)
	df := dataframe.ReadCSV(reader,
		dataframe.WithDelimiter(';'),
		dataframe.HasHeader(true))
	return df
}
