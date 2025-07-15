package lib

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
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
func WriteCsvByte(path string, table [][]string, delimiter rune) error {
	file, err := os.Create(path)
	defer file.Close()
	w := csv.NewWriter(file)
	w.Comma = delimiter
	defer w.Flush()
	// Using Write
	err = w.WriteAll(table)

	return err
}
func GetCsvByte(table [][]string, delimiter rune) ([]byte, error) {
	var res bytes.Buffer
	for _, row := range table {
		res.WriteString(strings.Join(row, ",") + string(delimiter))
	}
	return res.Bytes(), nil
}
