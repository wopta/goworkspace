package _script

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/transaction"
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
	payByLinkFormatDev       = "https://pre.fabrick.com/pacewhitelabel/landingpage-web/pay-by-link/%s/modalita-addebito"
)

type rowStruct struct {
	policyUid           string
	policyCode          string
	scheduleDate        string
	externalId          string
	providerName        string
	providerId          string
	paymentMethod       string
	paymentDate         time.Time
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
		if len(rows) == 0 {
			continue
		}

		output := make([]rowStruct, 0)
		for _, row := range rows {
			var out rowStruct

			splittedExternalId := strings.Split(lib.TrimSpace(row[externalIdCol]), "_")

			if len(splittedExternalId) < 3 {
				continue
			}

			payDate, err := time.Parse("2006-01-02", lib.TrimSpace(row[payDateCol]))
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			out.policyUid = lib.TrimSpace(splittedExternalId[0])
			out.policyCode = lib.TrimSpace(splittedExternalId[2])
			out.scheduleDate = lib.TrimSpace(splittedExternalId[1])
			out.externalId = lib.TrimSpace(row[externalIdCol])
			out.providerName = models.FabrickPaymentProvider
			out.providerId = lib.TrimSpace(row[providerIdCol])
			out.paymentMethod = lib.TrimSpace(row[paymentInstrumentTypeCol])
			out.paymentDate = payDate
			// TODO: write payUrl builder
			out.payUrl = fmt.Sprintf(payByLinkFormatDev, lib.TrimSpace(row[providerIdCol]))
			out.userToken = lib.TrimSpace(row[customerIdCol])
			out.paymentInstrumentId = lib.TrimSpace(row[paymentInstrumentIdCol])

			output = append(output, out)
		}

		if len(output) == 0 {
			continue
		}

		sort.Slice(output, func(i, j int) bool {
			return output[i].scheduleDate < output[j].scheduleDate
		})
		parsedRows[key] = output
	}

	for _, rows := range parsedRows {
		policyUid := rows[len(rows)-1].policyUid
		userToken := rows[len(rows)-1].userToken

		//log.Printf("PolicyUid: %s - UserToken: %s", policyUid, userToken)

		transactions := transaction.GetPolicyActiveTransactions("", policyUid)
		if len(transactions) == 0 {
			log.Printf("no transactions found for policy %s", policyUid)
			continue
		}

		transactions = lib.SliceMap(transactions, func(tr models.Transaction) models.Transaction {
			if tr.UserToken == "" {
				tr.UserToken = userToken
			}
			return tr
		})

		for _, row := range rows {
			index := -1
			for trIndex, tr := range transactions {
				if tr.ScheduleDate == row.scheduleDate {
					index = trIndex
					break
				}
			}

			if index == -1 {
				// TODO: handle error
				continue
			}

			transactions[index].ProviderName = row.providerName
			if transactions[index].ProviderId == "" {
				transactions[index].ProviderId = row.providerId
			}
			if transactions[index].PaymentMethod == "" {
				transactions[index].PaymentMethod = row.paymentMethod
			}
			if transactions[index].TransactionDate.IsZero() {
				transactions[index].TransactionDate = row.paymentDate
			}
			transactions[index].PayUrl = row.payUrl
			transactions[index].UpdateDate = time.Now().UTC()
		}

		log.Printf("PolicyUid: %s", policyUid)
		lib.SliceMap(transactions, func(tr models.Transaction) models.Transaction {
			log.Printf("Transaction: %v", tr)
			return tr
		})

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
