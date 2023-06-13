package models

import (
	"time"
)

type Agency struct {
	Uid          string    `json:"uid" firestore:"uid" bigquery:"-"`
	Name         string    `json:"name" firestore:"name" bigquery:"-"`
	Manager      User      `json:"manager" firestore:"manager" bigquery:"-"`
	Portfolio    []string  `json:"portfolio" firestore:"portfolio" bigquery:"-"`                           // will contain users UIDs
	ParentAgency string    `json:"parentAgency,omitempty" firestore:"parentAgency,omitempty" bigquery:"-"` // parent Agency UID
	Agencies     []string  `json:"agencies" firestore:"agencies" bigquery:"-"`                             // will contain agencies UIDs
	Agents       []string  `json:"agents" firestore:"agents" bigquery:"-"`                                 // will contain agents UIDs
	IsActive     bool      `json:"isActive" firestore:"isActive" bigquery:"-"`
	Products     []Product `json:"products" firestore:"products" bigquery:"-"`
	Policies     []string  `json:"policies" firestore:"policies" bigquery:"-"` // will contain policies UIDs
	Steps        []Step    `json:"steps" firestore:"steps" bigquery:"-"`
	CreationDate time.Time `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	UpdateDate   time.Time `json:"updateDate" firestore:"updateDate" bigquery:"-"`
}

func SetAgencyPolicy(policy Policy, amount float64, schedule string) Transaction {

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
