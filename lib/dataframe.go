package lib

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

//import "github.com/go-gota/gota/series"

func CsvToDataframe(data []byte) dataframe.DataFrame {
	reader := bytes.NewReader(data)
	df := dataframe.ReadCSV(reader,
		dataframe.WithDelimiter(';'),
		dataframe.HasHeader(true),
		dataframe.NaNValues(nil))
	log.Println(df.Error())
	return df
}
func FileToDf(path string, delimiter rune, header bool) dataframe.DataFrame {
	log.Println("Opening a file ")
	var file, err = os.Open(path)
	CheckError(err)
	defer file.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	reader := bytes.NewReader(buf.Bytes())
	df := dataframe.ReadCSV(reader,
		dataframe.WithDelimiter(delimiter),
		dataframe.HasHeader(header))
	return df
}
func FilesToDf(filenames []string) dataframe.DataFrame {
	log.Println("Opening a file ")
	buf := bytes.NewBuffer(nil)
	for _, filename := range filenames {
		f, err := os.Open(filename)
		CheckError(err)
		io.Copy(buf, f) // Error handling elided for brevity.
		defer f.Close()
	}
	//s := string(buf.Bytes())

	reader := bytes.NewReader(buf.Bytes())
	df := dataframe.ReadCSV(reader,
		dataframe.WithDelimiter(';'),
		dataframe.HasHeader(true))
	return df
}

// ExportToCSV exports a Dataframe to a CSV file.
/*func ExportToCSV(ctx context.Context, w io.Writer, df *dfgo.DataFrame, options ...CSVExportOptions) error {

	df.Lock()
	defer df.Unlock()

	header := []string{}

	var r dfgo.Range

	nullString := "NaN" // Default will be "NaN"

	cw := csv.NewWriter(w)

	if len(options) > 0 {
		cw.Comma = options[0].Separator
		cw.UseCRLF = options[0].UseCRLF
		r = options[0].Range
		if options[0].NullString != nil {
			nullString = *options[0].NullString
		}
	}

	for _, aSeries := range df.Series {
		header = append(header, aSeries.Name())
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	nRows := df.NRows(dfgo.DontLock)

	if nRows > 0 {

		s, e, err := r.Limits(nRows)
		if err != nil {
			return err
		}

		flushCount := 0
		for row := s; row <= e; row++ {

			if err := ctx.Err(); err != nil {
				return err
			}

			flushCount++
			// flush after every 100 writes
			if flushCount > 100 { // flush in the 101th count
				cw.Flush()
				if err := cw.Error(); err != nil {
					return err
				}
				flushCount = 1
			}

			sVals := []string{}
			for _, aSeries := range df.Series {
				val := aSeries.Value(row)
				if val == nil {
					sVals = append(sVals, nullString)
				} else {
					sVals = append(sVals, aSeries.ValueString(row, dfgo.DontLock))
				}
			}

			// Write every row
			if err := cw.Write(sVals); err != nil {
				return err
			}
		}

	}

	// flush before exit
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}*/
