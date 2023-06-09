package lib

import (
	"encoding/csv"
	"os"
)

func WriteCsv(path string, table [][]string, delimiter rune) error {
	file, err := os.Create(path)
	defer file.Close()
	w := csv.NewWriter(file)
	w.Comma = delimiter
	defer w.Flush()
	// Using Write
	err = w.WriteAll(table)
	return err
}
