package models

import (
	"log"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/firestore"
	"github.com/wopta/goworkspace/lib"
	"google.golang.org/api/iterator"
)

type Transaction struct {
	Id                 string         `firestore:"id,omitempty" json:"id,omitempty" bigquery:"-"`
	Amount             float64        `firestore:"amount,omitempty" json:"amount,omitempty" bigquery:"amount" `
	Commissions        float64        `firestore:"commissions,omitempty" json:"commissions,omitempty" bigquery:"commissions"`
	CommissionsCompany float64        `firestore:"commissionsCompany,omitempty" json:"commissionsCompany,omitempty" bigquery:"commissionsCompany"`
	Status             string         `firestore:"status ,omitempty" json:"status ,omitempty" bigquery:"status"`
	PolicyName         string         `firestore:"policyName,omitempty" json:"policName,omitempty" bigquery:"policyName"`
	Name               string         `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	Commission         float64        `firestore:"commission,omitempty" json:"commission,omitempty" bigquery:"commission"`
	ScheduleDate       string         `firestore:"scheduleDate,omitempty" json:"scheduleDate,omitempty" bigquery:"scheduleDate"`
	PayDate            time.Time      `firestore:"payDate,omitempty" json:"payDate,omitempty" bigquery:"-"`
	CreationDate       time.Time      `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	BigPayDate         civil.DateTime `firestore:"-" json:"-" bigquery:"payDate"`
	BigCreationDate    civil.DateTime `firestore:"-" json:"-" bigquery:"creationDate"`
	Uid                string         `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	PolicyUid          string         `firestore:"policyUid,omitempty" json:"policyUid,omitempty" bigquery:"policyUid"`
	Company            string         `firestore:"company,omitempty" json:"company,omitempty" bigquery:"company"`
	NumberCompany      string         `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty" bigquery:"numberCompany"`
	StatusHistory      []string       `firestore:"statusHistory,omitempty" json:"statusHistory ,omitempty" bigquery:"-"`
	BigStatusHistory   string         `firestore:"-" json:"-" bigquery:"statusHistory"`
	IsPay              bool           `firestore:"isPay" json:"isPay,omitempty" bigquery:"isPay"`
	IsEmit             bool           `firestore:"isEmit" json:"isEmit,omitempty" bigquery:"isEmit"`
	IsDelete           bool           `json:"isDelete" firestore:"isDelete" bigquery:"-"`
	ProviderId         string         `json:"providerId" firestore:"providerId" bigquery:"-"`
	UserToken          string         `json:"userToken" firestore:"userToken" bigquery:"-"`
	ProviderName       string         `json:"providerName" firestore:"providerName" bigquery:"-"`
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
