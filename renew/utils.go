package renew

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"google.golang.org/api/iterator"
)

// TODO: remove me
const (
	policyRenewedTestCollection      string = "policyRenewedTest"
	transactionRenewedTestCollection string = "transactionRenewedTest"
)

func createSaveBatch(policy models.Policy, transactions []models.Transaction) map[string]map[string]interface{} {
	policy.BigQueryParse()
	batch := map[string]map[string]interface{}{
		policyRenewedTestCollection: { // TODO: change to lib.PolicyCollection
			policy.Uid: policy,
		},
		transactionRenewedTestCollection: {}, // TODO: change to lib.TransactionCollection
	}

	for idx, tr := range transactions {
		tr.UpdateDate = time.Now().UTC()
		tr.BigQueryParse()
		batch[transactionRenewedTestCollection][tr.Uid] = tr // TODO: change to lib.TransactionCollection
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

// In case we need to get the data from BigQuery. Shouldn't be used now
// because bigquery does not have all data
func GetTransactionsByPolicyAnnuity(policyUid string, annuity int) ([]models.Transaction, error) {
	var (
		query  bytes.Buffer
		params = make(map[string]interface{})
	)

	params["policyUid"] = policyUid
	params["annuity"] = annuity

	query.WriteString(fmt.Sprintf("SELECT * FROM `%s.%s` WHERE "+
		"policyUid = '@policyUid' AND "+
		"annuity = @annuity",
		models.WoptaDataset,
		renewTransactionCollection))

	return lib.QueryParametrizedRowsBigQuery[models.Transaction](query.String(), params)
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
		Message:      message,
		Title:        title,
		Subject:      subject,
		IsHtml:       true,
		IsAttachment: true,
		Attachments: &[]mail.Attachment{{
			Name:        fmt.Sprintf("report-%s-%d.json", date.Format(time.DateOnly), time.Now().Unix()),
			Byte:        base64.StdEncoding.EncodeToString(responseJson),
			FileName:    fmt.Sprintf("report-%s-%d.json", date.Format(time.DateOnly), time.Now().Unix()),
			ContentType: "application/json",
		}},
	})
}
