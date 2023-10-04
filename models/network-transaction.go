package models

import (
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

const (
	AccountTypeActive                 string = "Active"
	AccountTypePassive                string = "Passive"
	PaymentTypeRemittanceCompany      string = "RemittanceCompany"
	PaymentTypeRemittanceMga          string = "RemittanceMga"
	PaymentTypeCommission             string = "Commission"
	NetworkTransactionStatusCreated   string = "Created"
	NetworkTransactionStatusToPay     string = "ToPay"
	NetworkTransactionStatusPaid      string = "Paid"
	NetworkTransactionStatusConfirmed string = "Confirmed"
)

type NetworkTransaction struct {
	Uid              string                `json:"uid" bigquery:"uid"`
	PolicyUid        string                `json:"policyUid" bigquery:"policyUid"`
	TransactionUid   string                `json:"transactionUid" bigquery:"transactionUid"`
	NetworkUid       string                `json:"networkUid" bigquery:"networkUid"`
	NetworkNodeUid   string                `json:"networkNodeUid" bigquery:"networkNodeUid"`
	NetworkNodeType  string                `json:"networkNodeType" bigquery:"networkNodeType"`
	AccountType      string                `json:"accountType" bigquery:"accountType"` // AccountTypeActive | AccountTypePassive
	PaymentType      string                `json:"paymentType" bigquery:"paymentType"` // PaymentTypeRemittance | PaymentTypeCommission
	Amount           float64               `json:"amount" bigquery:"amount"`
	AmountNet        float64               `json:"amountNet" bigquery:"amountNet"`
	Name             string                `json:"name" bigquery:"name"`
	Status           string                `json:"status" bigquery:"status"`
	StatusHistory    []string              `json:"statusHistory" bigquery:"statusHistory"`
	IsPay            bool                  `json:"isPay" bigquery:"isPay"`
	IsConfirmed      bool                  `json:"isConfirmed" bigquery:"isConfirmed"`
	CreationDate     bigquery.NullDateTime `json:"creationDate" bigquery:"creationDate"`
	PayDate          bigquery.NullDateTime `json:"payDate" bigquery:"payDate"`
	TransactionDate  bigquery.NullDateTime `json:"transactionDate" bigquery:"transactionDate"`
	ConfirmationDate bigquery.NullDateTime `json:"confirmationDate" bigquery:"confirmationDate"`
}

func (nt *NetworkTransaction) SaveBigQuery() error {
	log.Println("[NetworkTransaction.SaveBigQuery]")

	var (
		err       error
		datasetId = "test1" // WoptaDataset
		tableId   = NetworkTransactionCollection
	)

	baseQuery := fmt.Sprintf("SELECT * FROM `%s.%s` WHERE ", datasetId, tableId)
	whereClause := fmt.Sprintf("uid = '%s'", nt.Uid)
	query := fmt.Sprintf("%s %s", baseQuery, whereClause)

	result, err := lib.QueryRowsBigQuery[NetworkTransaction](query)
	if err != nil {
		log.Printf("[NetworkTransaction.SaveBigQuery] error querying db with query %s: %s", query, err.Error())
		return err
	}

	if len(result) == 0 {
		log.Printf("[NetworkTransaction.SaveBigQuery] creating new NetworkTransaction %s", nt.Uid)
		err = lib.InsertRowsBigQuery(datasetId, tableId, nt)
	} else {
		log.Printf("[NetworkTransaction.SaveBigQuery] updating NetworkTransaction %s", nt.Uid)
		updatedFields := make(map[string]interface{})
		updatedFields["status"] = nt.Status
		updatedFields["statusHistory"] = nt.StatusHistory
		updatedFields["isPay"] = nt.IsPay
		updatedFields["isConfirmed"] = nt.IsConfirmed
		updatedFields["payDate"] = nt.PayDate
		updatedFields["transactionDate"] = nt.TransactionDate
		updatedFields["confirmationDate"] = nt.ConfirmationDate

		err = lib.UpdateRowBigQueryV2(datasetId, tableId, updatedFields, whereClause)
	}

	if err != nil {
		log.Printf("[NetworkTransaction.SaveBigQuery] error saving to db: %s", err.Error())
	} else {
		log.Println("[NetworkTransaction.SaveBigQuery] NetworkTransaction saved!")
	}

	return err
}
