package models

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Transaction struct {
	Id                 string                `firestore:"id,omitempty" json:"id,omitempty" bigquery:"-"`
	Amount             float64               `firestore:"amount,omitempty" json:"amount,omitempty" bigquery:"amount" `
	AmountNet          float64               `json:"amountNet,omitempty" firestore:"amountNet,omitempty" bigquery:"amountNet"`
	AgentUid           string                `json:"agentUid,omitempty" firestore:"agentUid,omitempty" bigquery:"agentUid"`    // DEPRECATED
	AgencyUid          string                `json:"agencyUid,omitempty" firestore:"agencyUid,omitempty" bigquery:"agencyUid"` // DEPRECATED
	Commissions        float64               `firestore:"commissions,omitempty" json:"commissions,omitempty" bigquery:"commissions"`
	CommissionsCompany float64               `firestore:"commissionsCompany,omitempty" json:"commissionsCompany,omitempty" bigquery:"commissionsCompany"` // DEPRECATED
	CommissionsAgent   float64               `firestore:"commissionsAgent,omitempty" json:"commissionsAgent,omitempty" bigquery:"commissionsAgent"`       // DEPRECATED
	CommissionsAgency  float64               `firestore:"commissionsAgency,omitempty" json:"commissionsAgency,omitempty" bigquery:"commissionsAgency"`    // DEPRECATED
	Status             string                `firestore:"status,omitempty" json:"status,omitempty" bigquery:"status"`
	PolicyName         string                `firestore:"policyName,omitempty" json:"policName,omitempty" bigquery:"policyName"`
	Name               string                `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	Commission         float64               `firestore:"commission,omitempty" json:"commission,omitempty" bigquery:"commission"` // DEPRECATED
	ScheduleDate       string                `firestore:"scheduleDate,omitempty" json:"scheduleDate,omitempty" bigquery:"scheduleDate"`
	ExpirationDate     string                `json:"expirationDate,omitempty" firestore:"expirationDate,omitempty" bigquery:"expirationDate"`
	PayDate            time.Time             `firestore:"payDate,omitempty" json:"payDate,omitempty" bigquery:"-"`
	CreationDate       time.Time             `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	TransactionDate    time.Time             `firestore:"transactionDate,omitempty" json:"transactionDate,omitempty" bigquery:"-"`
	BigPayDate         bigquery.NullDateTime `firestore:"-" json:"-" bigquery:"payDate"`
	BigCreationDate    civil.DateTime        `firestore:"-" json:"-" bigquery:"creationDate"`
	BigTransactionDate bigquery.NullDateTime `firestore:"-" json:"-" bigquery:"transactionDate"`
	Uid                string                `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	PolicyUid          string                `firestore:"policyUid,omitempty" json:"policyUid,omitempty" bigquery:"policyUid"`
	Company            string                `firestore:"company,omitempty" json:"company,omitempty" bigquery:"company"`
	NumberCompany      string                `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty" bigquery:"numberCompany"`
	StatusHistory      []string              `firestore:"statusHistory,omitempty" json:"statusHistory,omitempty" bigquery:"-"`
	BigStatusHistory   string                `firestore:"-" json:"-" bigquery:"statusHistory"`
	IsPay              bool                  `firestore:"isPay" json:"isPay,omitempty" bigquery:"isPay"`
	IsEmit             bool                  `firestore:"isEmit" json:"isEmit,omitempty" bigquery:"isEmit"`
	IsDelete           bool                  `json:"isDelete" firestore:"isDelete" bigquery:"isDelete"`
	ProviderId         string                `json:"providerId" firestore:"providerId" bigquery:"-"`
	UserToken          string                `json:"userToken" firestore:"userToken" bigquery:"-"`
	ProviderName       string                `json:"providerName" firestore:"providerName" bigquery:"-"`
	PaymentMethod      string                `firestore:"paymentMethod,omitempty" json:"paymentMethod,omitempty" bigquery:"paymentMethod"`
	PaymentNote        string                `firestore:"paymentNote,omitempty" json:"paymentNote,omitempty" bigquery:"paymentNote"`
	NetworkCommissions map[string]float64    `json:"networkCommissions,omitempty" firestore:"networkCommissions,omitempty" bigquery:"-"` // DEPRECATED
	UpdateDate         time.Time             `json:"updateDate" firestore:"updateDate" bigquery:"-"`
	BigUpdateDate      bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"updateDate"`
	EffectiveDate      time.Time             `json:"effectiveDate,omitempty" firestore:"effectiveDate,omitempty" bigquery:"-"`
	BigEffectiveDate   bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"effectiveDate"`
}

func TransactionToListData(query *firestore.DocumentIterator) []Transaction {
	result := make([]Transaction, 0)
	for {
		d, err := query.Next()
		if err != nil {
		}
		if err != nil {
			if err == iterator.Done {
				break
			}
		}
		var value Transaction
		e := d.DataTo(&value)
		value.Uid = d.Ref.ID
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}

func SetTransactionPolicy(policy Policy, id string, amount float64, schedule string, Commissions float64, company string) Transaction {
	return Transaction{
		Amount:        amount,
		Id:            id,
		PolicyName:    policy.Name,
		PolicyUid:     policy.Uid,
		CreationDate:  time.Now(),
		Status:        TransactionStatusToPay,
		StatusHistory: []string{TransactionStatusToPay},
		ScheduleDate:  schedule,
		NumberCompany: policy.CodeCompany,
		Commissions:   Commissions,
		IsPay:         false,
		Name:          policy.Contractor.Name + " " + policy.Contractor.Surname,
		Company:       company,
	}
}

func (t *Transaction) Normalize() {
	t.Name = lib.TrimSpace(t.Name)
	t.PaymentNote = lib.ToUpper(t.PaymentNote)
}

func (t *Transaction) BigQuerySave(origin string) {
	fireTransactions := lib.GetDatasetByEnv(origin, TransactionsCollection)
	transactionJson, err := json.Marshal(t)
	if err != nil {
		log.Println("ERROR Transaction "+t.Uid+" Marshal: ", err)
		return
	}
	log.Println("Transaction: "+t.Uid, string(transactionJson))
	t.BigPayDate = lib.GetBigQueryNullDateTime(t.PayDate)
	t.BigTransactionDate = lib.GetBigQueryNullDateTime(t.TransactionDate)
	t.BigCreationDate = civil.DateTimeOf(t.CreationDate)
	t.BigStatusHistory = strings.Join(t.StatusHistory, ",")
	t.BigUpdateDate = lib.GetBigQueryNullDateTime(t.UpdateDate)
	t.BigEffectiveDate = lib.GetBigQueryNullDateTime(t.EffectiveDate)
	log.Println("Transaction save BigQuery: " + t.Uid)

	err = lib.InsertRowsBigQuery(WoptaDataset, fireTransactions, t)
	if err != nil {
		log.Println("ERROR Transaction "+t.Uid+" save BigQuery: ", err)
		return
	}
	log.Println("Transaction BigQuery saved!")
}

func (t *Transaction) IsLate(limit time.Time) bool {
	return t.PayDate.After(limit)
}
