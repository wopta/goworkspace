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
	Amount             float64        `firestore:"amount,omitempty" json:"amount,omitempty" bigquery:"amount" `
	Commissions        float64        `firestore:"commissions,omitempty" json:"commissions,omitempty" bigquery:"commissions"`
	CommissionsCompany float64        `firestore:"commissionsCompany,omitempty" json:"commissionsCompany,omitempty" bigquery:"commissionsCompany"`
	Status             string         `firestore:"status ,omitempty" json:"status ,omitempty" bigquery:"status"`
	PolicyName         string         `firestore:"policyName,omitempty" json:"policName,omitempty" bigquery:"policyName"`
	Name               string         `firestore:"name,omitempty" json:"name,omitempty" bigquery:"name"`
	ScheduleDate       string         `firestore:"scheduleDate,omitempty" json:"scheduleDate,omitempty" bigquery:"scheduleDate"`
	PayDate            time.Time      `firestore:"payDate,omitempty" json:"payDate,omitempty" bigquery:"-"`
	CreationDate       time.Time      `firestore:"creationDate,omitempty" json:"creationDate,omitempty" bigquery:"-"`
	BigPayDate         civil.DateTime `bigquery:"payDate"`
	BigCreationDate    civil.DateTime `bigquery:"creationDate"`
	Uid                string         `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	PolicyUid          string         `firestore:"policyUid,omitempty" json:"policyUid,omitempty" bigquery:"policyUid"`
	Company            string         `firestore:"company,omitempty" json:"company,omitempty" bigquery:"company"`
	NumberCompany      string         `firestore:"numberCompany,omitempty" json:"numberCompany,omitempty" bigquery:"numberCompany"`
	StatusHistory      []string       `firestore:"statusHistory,omitempty" json:"statusHistory ,omitempty" bigquery:"-"`
	BigStatusHistory   string         `bigquery:"statusHistory"`
	IsPay              bool           `firestore:"isPay,omitempty" json:"isPay,omitempty" bigquery:"isPay"`
	IsEmit             bool           `firestore:"isEmit,omitempty" json:"isEmit,omitempty" bigquery:"isEmit"`
}

func TransactionToListData(query *firestore.DocumentIterator) []Transaction {
	var result []Transaction
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

func SetTransactionPolicy(policy Policy, amount float64, schedule string) Transaction {

	return Transaction{
		Amount:        amount,
		PolicyName:    policy.Name,
		PolicyUid:     policy.Uid,
		CreationDate:  time.Now(),
		Status:        TransactionStatusToPay,
		StatusHistory: []string{TransactionStatusToPay},
		ScheduleDate:  schedule,
		NumberCompany: policy.CodeCompany,
		IsPay:         false,
		Name:          policy.Contractor.Name + " " + policy.Contractor.Surname,
	}
}
