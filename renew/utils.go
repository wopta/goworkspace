package renew

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.dev.wopta.it/goworkspace/lib"
	env "gitlab.dev.wopta.it/goworkspace/lib/environment"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"google.golang.org/api/iterator"
)

func getProducts(policyType, quoteType string) map[string]map[string]models.Product {
	var (
		products = make(map[string]map[string]models.Product)
		channels = []string{models.MgaChannel, models.NetworkChannel, models.ECommerceChannel}
	)

	fileList := getProductsFileList()

	getProductsByChannel := func(channel string) {
		filteredFileList := lib.SliceFilter(fileList, func(file string) bool {
			filenameParts := strings.SplitN(file, "/", 4)
			return strings.HasPrefix(filenameParts[3], channel)
		})

		products[channel] = make(map[string]models.Product)
		channelProducts := getProductsFromFileList(filteredFileList)

		for _, prd := range channelProducts {
			if strings.EqualFold(prd.PolicyType, policyType) && strings.EqualFold(prd.QuoteType, quoteType) {
				key := fmt.Sprintf("%s-%s", prd.Name, prd.Version)
				products[channel][key] = prd
			}
		}
	}

	for _, channel := range channels {
		getProductsByChannel(channel)
	}

	return products
}

func getProductsFileList() []string {
	var (
		err      error
		fileList = make([]string, 0)
	)

	switch os.Getenv("env") {
	case env.Local, env.LocalTest:
		fileList, err = lib.ListLocalFolderContent(models.ProductsFolder)
	default:
		fileList, err = lib.ListGoogleStorageFolderContent(models.ProductsFolder)
	}

	if err != nil {
		log.ErrorF("error getting file list: %s", err.Error())
	}

	checkedList := lib.SliceFilter(fileList, checkSlashes)

	return checkedList
}

func checkSlashes(s string) bool {
	// Correct path is: products/{{product_dir}}/{{version_number}}/{{filename.extension}}
	// but this function supports further nesting
	return strings.Count(s, "/") >= 3
}

func getProductsFromFileList(fileList []string) []models.Product {
	var (
		err        error
		products   = make([]models.Product, 0)
		fileChunks = make([][]string, 0)
	)

	if len(fileList) == 0 {
		return products
	}

	// create subarrays for each different product
	for _, file := range fileList {
		filenameParts := strings.SplitN(file, "/", 4)
		productName := filenameParts[1]
		if len(fileChunks) == 0 {
			fileChunks = append(fileChunks, make([]string, 0))
		}
		if len(fileChunks[len(fileChunks)-1]) > 0 {
			chunkProductName := strings.SplitN(fileChunks[len(fileChunks)-1][0], "/", 3)[1]
			if chunkProductName != productName {
				fileChunks = append(fileChunks, make([]string, 0))
			}
		}
		fileChunks[len(fileChunks)-1] = append(fileChunks[len(fileChunks)-1], file)
	}

	// loop each product
	for _, chunk := range fileChunks {
		// sort them by the last version
		sort.Slice(chunk, func(i, j int) bool {
			return strings.SplitN(chunk[i], "/", 4)[2] > strings.SplitN(chunk[j], "/", 4)[2]
		})
		// loop each version
		for _, file := range chunk {
			var currentProduct models.Product
			// download file from bucket
			fileBytes := lib.GetFilesByEnv(file)
			err = json.Unmarshal(fileBytes, &currentProduct)
			if err != nil {
				continue
			}

			products = append(products, currentProduct)
		}
		if err != nil {
			break
		}
	}

	return products
}

func createDraftSaveBatch(policy models.Policy, transactions []models.Transaction) map[string]map[string]interface{} {
	var (
		polCollection = collectionPrefix + lib.RenewPolicyCollection
		trsCollection = collectionPrefix + lib.RenewTransactionCollection
	)

	policy.Updated = time.Now().UTC()
	policy.BigQueryParse()
	batch := map[string]map[string]interface{}{
		polCollection: {
			policy.Uid: policy,
		},
		trsCollection: {},
	}

	for idx, tr := range transactions {
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		batch[trsCollection][tr.Uid] = tr
		transactions[idx] = tr
	}

	return batch
}

func createPromoteSaveBatch(policy models.Policy, transactions []models.Transaction) map[string]map[string]interface{} {
	var (
		polCollection string = collectionPrefix + lib.PolicyCollection
		trsCollection string = collectionPrefix + lib.TransactionsCollection
	)

	policy.Updated = time.Now().UTC()
	policy.BigQueryParse()
	batch := map[string]map[string]interface{}{
		polCollection: {
			policy.Uid: policy,
		},
		trsCollection: {},
	}

	for idx, tr := range transactions {
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		batch[trsCollection][tr.Uid] = tr
		transactions[idx] = tr
	}

	return batch
}

func createPromoteProcessedBatch(policy models.Policy, transactions []models.Transaction) map[string]map[string]interface{} {
	var (
		polCollection = collectionPrefix + lib.RenewPolicyCollection
		trsCollection = collectionPrefix + lib.RenewTransactionCollection
	)

	policy.IsDeleted = true
	policy.Updated = time.Now().UTC()
	policy.BigQueryParse()
	batch := map[string]map[string]interface{}{
		polCollection: {
			policy.Uid: policy,
		},
		trsCollection: {},
	}

	for idx, tr := range transactions {
		tr.IsDelete = true
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		batch[trsCollection][tr.Uid] = tr
		transactions[idx] = tr
	}

	return batch
}

func saveToDatabases(data map[string]map[string]interface{}) error {
	err := lib.SetBatchFirestoreErr(data)
	if err != nil {
		return err
	}

	for collection, values := range data {
		dataToSave := make([]interface{}, 0)
		for _, value := range values {
			dataToSave = append(dataToSave, value)
		}
		err = lib.InsertRowsBigQuery(models.WoptaDataset, collection, dataToSave)
		if err != nil {
			return err
		}
	}

	return nil
}

type firestoreQuery struct {
	field      string
	operator   string
	queryValue interface{}
}

func firestoreWhere[T any](collection string, queries []firestoreQuery) (documents []T, err error) {
	var (
		client *firestore.Client
		query  firestore.Query
		ctx    context.Context = context.Background()
	)

	if client, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID")); err != nil {
		return nil, err
	}

	colRef := client.Collection(collection)

	for idx, q := range queries {
		if idx == 0 {
			query = colRef.Where(q.field, q.operator, q.queryValue)
			continue
		}
		query = query.Where(q.field, q.operator, q.queryValue)
	}

	docIterator := query.Documents(ctx)

	for {
		var (
			snapshot *firestore.DocumentSnapshot
			document T
		)
		if snapshot, err = docIterator.Next(); err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = snapshot.DataTo(&document)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func sendReportMail(date time.Time, report RenewResp, isDraft bool) {
	var (
		message string = fmt.Sprintf(`
		<p>Quietanzamento del %s</p>
		<p>Con successo: %d</p>
		<p>Con errori: %d</p>
		<p>Per report completo vedere file json in allegato.</p>
		`, date.Format(time.DateOnly), len(report.Success), len(report.Failure))
		title   string = "Report quietanzamento"
		subject string = fmt.Sprintf("Report quietanzamento del %s", date.Format(time.DateOnly))
	)

	if isDraft {
		title = "Report quietanzamento provvisorio"
		subject = fmt.Sprintf("Report quietanzamento provvisorio del %s", date.Format(time.DateOnly))
		message = fmt.Sprintf(`
		<p>Quietanzamento provvisorio del %s</p>
		<p>Con successo: %d</p>
		<p>Con errori: %d</p>
		<p>Per report completo vedere file json in allegato.</p>
		`, date.Format(time.DateOnly), len(report.Success), len(report.Failure))
	}

	responseJson, _ := json.Marshal(report)

	mail.SendMail(mail.MailRequest{
		FromAddress:  mail.AddressAnna,
		To:           []string{mail.AddressTechnology.Address},
		Cc:           mail.AddressOperations.Address,
		Message:      message,
		Title:        title,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: true,
		Attachments: &[]models.Attachment{{
			Name:        fmt.Sprintf("report-%s-%d.json", date.Format(time.DateOnly), time.Now().Unix()),
			Byte:        base64.StdEncoding.EncodeToString(responseJson),
			FileName:    fmt.Sprintf("report-%s-%d.json", date.Format(time.DateOnly), time.Now().Unix()),
			ContentType: "application/json",
		}},
	})
}
