package _script

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	providerIdCol            = 0
	externalIdCol            = 1
	descriptionCol           = 3
	statusCol                = 8
	payDateCol               = 9
	customerIdCol            = 10
	paymentInstrumentIdCol   = 11
	paymentInstrumentTypeCol = 12
)

type rowStruct struct {
	policyUid           string
	policyCode          string
	scheduleDate        string
	externalId          string
	providerName        string
	providerId          string
	paymentMethod       string
	paymentDate         string
	payUrl              string
	userToken           string
	paymentInstrumentId string
}

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

	// grouping rows by Description
	data := groupBy(df, descriptionCol)

	// filter rows by status
	filteredData := make(map[string][][]string)
	for groupKey, rows := range data {
		outputRows := make([][]string, 0)
		for _, row := range rows {
			if row[statusCol] == "OK" {
				filteredData[groupKey] = append(filteredData[groupKey], row)
			}
		}

		if len(outputRows) > 0 {
			filteredData[groupKey] = outputRows
		}
	}

	// parse rows into struct
	parsedRows := make(map[string][]rowStruct)
	for key, rows := range filteredData {
		output := make([]rowStruct, 0)
		for _, row := range rows {
			var out rowStruct

			splittedExternalId := strings.Split(lib.TrimSpace(row[externalIdCol]), "_")

			if len(splittedExternalId) < 3 {
				continue
			}

			out.policyUid = lib.TrimSpace(splittedExternalId[0])
			out.policyCode = lib.TrimSpace(splittedExternalId[2])
			out.scheduleDate = lib.TrimSpace(splittedExternalId[1])
			out.externalId = lib.TrimSpace(row[externalIdCol])
			out.providerName = models.FabrickPaymentProvider
			out.providerId = lib.TrimSpace(row[providerIdCol])
			out.paymentMethod = lib.TrimSpace(row[paymentInstrumentTypeCol])
			out.paymentDate = lib.TrimSpace(row[payDateCol])
			// TODO: write payUrl builder
			out.payUrl = "www.wopta.it"
			out.userToken = lib.TrimSpace(row[customerIdCol])
			out.paymentInstrumentId = lib.TrimSpace(row[paymentInstrumentIdCol])

			output = append(output, out)
		}

		sort.Slice(output, func(i, j int) bool {
			return output[i].scheduleDate < output[j].scheduleDate
		})
		parsedRows[key] = output
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
