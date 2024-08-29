package _script

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
)

func EditLifePolicy(policyUid string) {
	rawData, err := os.ReadFile("./_script/policy_80.csv")
	if err != nil {
		log.Fatal(err)
	}

	df, err := lib.CsvToDataframeV2(rawData, ';', true)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range df.Records() {
		for index, col := range v {
			log.Printf("index: %02d, value: %v", index, col)
		}
	}

}

func groupByColumn(df dataframe.DataFrame, col int) map[string][][]string {
	res := make(map[string][][]string)
	for _, k := range df.Records() {
		res[k[col]] = append(res[k[col]], k)
	}
	return res
}
