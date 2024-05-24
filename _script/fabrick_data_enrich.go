package _script

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"log"
	"os"
)

func FabrickDataEnrich() {
	rawDoc, err := os.ReadFile("./_script/fabrick_data.csv")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	df, err := lib.CsvToDataframeV2(rawDoc, ';', true)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("#rows: %d", df.Nrow())
	log.Printf("#cols: %d", df.Ncol())

	data := groupBy(df, 3)

	for _, row := range data {
		log.Printf("row: %v", row)
	}

}

func groupBy(df dataframe.DataFrame, col int) map[string][][]string {
	res := make(map[string][][]string)
	for _, k := range df.Records() {
		if _, found := res[k[col]]; found {
			res[k[col]] = append(res[k[col]], k)
		} else {
			res[k[col]] = [][]string{k}
		}
	}
	return res
}
