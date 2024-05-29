package _script

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
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
	payByLinkFormatDev       = "pacewhitelabel/landingpage-web/pay-by-link/%s/modalita-addebito"
)

type rowStruct struct {
	PolicyUid           string    `json:"policyUid"`
	PolicyCode          string    `json:"policyCode"`
	ScheduleDate        string    `json:"scheduleDate"`
	ExternalId          string    `json:"externalId"`
	ProviderName        string    `json:"providerName"`
	ProviderId          string    `json:"providerId"`
	PaymentMethod       string    `json:"paymentMethod"`
	PaymentDate         time.Time `json:"paymentDate"`
	PayUrl              string    `json:"payUrl"`
	UserToken           string    `json:"userToken"`
	PaymentInstrumentId string    `json:"paymentInstrumentId"`
}

/*
Script to enrich all DB transactions with fabrick extracted data
*/
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
	filteredData := filterBy(data, statusCol, "OK")

	// filter valid policies checking externaId for the followinf format
	// 1. uid -> alpha-numeric 20 characters
	// 2. scheduleDate -> string in format time.DateOnly
	// the returned map is grouped by uid
	filteredData = filterByRegex(filteredData, externalIdCol, `^([a-zA-Z\d]){20}(_){1}((\d){4}-(\d){2}-(\d){2})`)

	// parse rows into struct
	parsedRows := parseRows(filteredData)

	var wg sync.WaitGroup

	// enrich transactions
	trToBeSaved := make([]models.Transaction, 0)
	notFound := make([]string, 0)
	noMatch := make([]rowStruct, 0)
	for _, rows := range parsedRows {
		wg.Add(1)
		go func(rows []rowStruct) {
			defer wg.Done()

			policyUid := rows[len(rows)-1].PolicyUid
			userToken := rows[len(rows)-1].UserToken

			transactions := transaction.GetPolicyActiveTransactions("", policyUid)
			if len(transactions) == 0 {
				notFound = append(notFound, policyUid)
				return
			}

			transactions = lib.SliceMap(transactions, func(tr models.Transaction) models.Transaction {
				if tr.UserToken == "" {
					tr.UserToken = userToken
					tr.UpdateDate = time.Now().UTC()
				}
				tr.BigQueryParse()
				return tr
			})

			for _, row := range rows {
				index := -1
				for trIndex, tr := range transactions {
					if tr.ScheduleDate == row.ScheduleDate || tr.ProviderId == row.ProviderId {
						index = trIndex
						break
					}
				}

				if index == -1 {
					// TODO: handle error
					noMatch = append(noMatch, row)
					continue
				}

				transactions[index].ProviderName = row.ProviderName
				if transactions[index].ProviderId == "" {
					transactions[index].ProviderId = row.ProviderId
				}
				if transactions[index].PaymentMethod == "" {
					transactions[index].PaymentMethod = row.PaymentMethod
				}
				if transactions[index].TransactionDate.IsZero() {
					transactions[index].TransactionDate = row.PaymentDate
				}
				transactions[index].PayUrl = row.PayUrl
				transactions[index].UpdateDate = time.Now().UTC()
				transactions[index].BigQueryParse()
			}

			trToBeSaved = append(trToBeSaved, transactions...)
		}(rows)
	}

	wg.Wait()

	// save transactions
	if err := saveData(trToBeSaved); err != nil {
		log.Fatalf("error saving data: %s", err)
	}

	// save report
	if err = saveReport(trToBeSaved, notFound, noMatch); err != nil {
		log.Printf("error saving report : %s", err.Error())
	}

	log.Println("Script done enriching transaction with fabrick data")
}

func groupBy(df dataframe.DataFrame, col int) map[string][][]string {
	res := make(map[string][][]string)
	for _, k := range df.Records() {
		res[k[col]] = append(res[k[col]], k)
	}
	return res
}

func filterBy(data map[string][][]string, col int, value string) map[string][][]string {
	filteredData := make(map[string][][]string)
	for groupKey, rows := range data {
		outputRows := make([][]string, 0)
		for _, row := range rows {
			if row[col] == value {
				outputRows = append(outputRows, row)
			}
		}

		if len(outputRows) > 0 {
			filteredData[groupKey] = outputRows
		}
	}
	return filteredData
}

func filterByRegex(data map[string][][]string, col int, regex string) map[string][][]string {
	filteredData := make(map[string][][]string)
	for _, rows := range data {
		outputRows := make([][]string, 0)
		key := ""
		for _, row := range rows {
			if matched, _ := regexp.MatchString(regex, row[col]); matched {
				key = strings.Split(row[col], "_")[0]
				outputRows = append(outputRows, row)
			}
		}

		if len(outputRows) > 0 {
			filteredData[key] = outputRows
		}
	}
	return filteredData
}

func parseRows(data map[string][][]string) map[string][]rowStruct {
	parsedRows := make(map[string][]rowStruct)
	for key, rows := range data {
		output := make([]rowStruct, 0)
		for _, row := range rows {
			var out rowStruct

			splittedExternalId := strings.Split(lib.TrimSpace(row[externalIdCol]), "_")

			if len(splittedExternalId) < 3 {
				log.Printf("[parseRows] not one of ours: %s", row[externalIdCol])
				continue
			}

			// check if second value is time.DateOnly
			if _, err := time.Parse(time.DateOnly, splittedExternalId[1]); err != nil {
				log.Printf("[parseRows] not one of ours: %s", row[externalIdCol])
				continue
			}

			payDate, err := time.Parse(time.DateOnly, lib.TrimSpace(row[payDateCol]))
			if err != nil {
				log.Printf("[parseRows] error: %v", err)
				continue
			}

			out.PolicyUid = lib.TrimSpace(splittedExternalId[0])
			out.PolicyCode = lib.TrimSpace(splittedExternalId[2])
			out.ScheduleDate = lib.TrimSpace(splittedExternalId[1])
			out.ExternalId = lib.TrimSpace(row[externalIdCol])
			out.ProviderName = models.FabrickPaymentProvider
			out.ProviderId = lib.TrimSpace(row[providerIdCol])
			out.PaymentMethod = lib.ToLower(row[paymentInstrumentTypeCol])
			out.PaymentDate = payDate
			out.PayUrl = os.Getenv("FABRICK_BASEURL") + fmt.Sprintf(payByLinkFormatDev, lib.TrimSpace(row[providerIdCol]))
			out.UserToken = lib.TrimSpace(row[customerIdCol])
			out.PaymentInstrumentId = lib.TrimSpace(row[paymentInstrumentIdCol])

			output = append(output, out)
		}

		if len(output) == 0 {
			continue
		}

		sort.Slice(output, func(i, j int) bool {
			return output[i].ScheduleDate < output[j].ScheduleDate
		})
		parsedRows[key] = output
	}

	return parsedRows
}

func saveData(trToBeSaved []models.Transaction) error {
	trMap := make(map[string]models.Transaction)
	for _, tr := range trToBeSaved {
		trMap[tr.Uid] = tr
	}
	firestoreBatch := map[string]map[string]models.Transaction{
		lib.TransactionsCollection: trMap,
	}
	if err := lib.SetBatchFirestoreErr(firestoreBatch); err != nil {
		return err
	}
	return lib.InsertRowsBigQuery(lib.WoptaDataset, lib.TransactionsCollection, trToBeSaved)
}

func saveReport(trToBeSaved []models.Transaction, notFound []string, noMatch []rowStruct) error {
	var (
		rawOutput []byte
		err       error
	)

	report := struct {
		Success  []models.Transaction `json:"success"`
		NotFound []string             `json:"notFound"`
		NoMatch  []rowStruct          `json:"noMatch"`
	}{
		Success:  trToBeSaved,
		NotFound: notFound,
		NoMatch:  noMatch,
	}

	if rawOutput, err = json.Marshal(report); err != nil {
		return err
	}

	return os.WriteFile("./import_fabrick_data_report.json", rawOutput, 777)
}
