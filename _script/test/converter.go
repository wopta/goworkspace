package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("nazioni_estere_cessate.csv") // Replace "values.csv" with the actual path to your CSV file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	data := make(map[string]map[string]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Use the second column as the key and the seventh column as the value
		key := strings.ToUpper(record[4])

		data[key] = map[string]string{
			//"codFisc": strings.ToUpper(record[4]),
			//"province": record[2],
			//"cap":      record[5],
			"city": record[7],
		}
	}

	// Print the generated map
	out, _ := json.Marshal(data)
	fmt.Println(string(out))
}
