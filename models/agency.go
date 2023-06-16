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
	Skin         Skin      `json:"skin" firestore:"skin" bigquery:"-"`
	RuiCode      string    `json:"ruiCode" firestore:"ruiCode" bigquery:"-"`
	CreationDate time.Time `json:"creationDate" firestore:"creationDate" bigquery:"-"`
	UpdateDate   time.Time `json:"updateDate" firestore:"updateDate" bigquery:"-"`
}

type Skin struct {
	PrimaryColor   string `json:"primaryColor" firestore:"primaryColor" bigquery:"-"`
	SecondaryColor string `json:"secondaryColor" firestore:"secondaryColor" bigquery:"-"`
	LogoUrl        string `json:"logoUrl" firestore:"logoUrl" bigquery:"-"`
}
