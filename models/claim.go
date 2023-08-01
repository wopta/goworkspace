package models

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

func UnmarshalClaim(data []byte) (Claim, error) {
	var r Claim
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Claim) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Claim struct {
	Name              string                `firestore:"name"                        json:"name,omitempty"              bigquery:"-"`
	Date              time.Time             `firestore:"date"                        json:"date,omitempty"              bigquery:"-"`
	BigDate           bigquery.NullDateTime `firestore:"-"                        json:"-"              bigquery:"date"`
	Surname           string                `firestore:"surname"                     json:"surname,omitempty"           bigquery:"-"`
	Mail              string                `firestore:"mail"                        json:"mail,omitempty"              bigquery:"-"`
	PolicyDescription string                `firestore:"policyDescription,omitempty" json:"policyDescription,omitempty" bigquery:"-"`
	PolicyId          string                `firestore:"policyId,omitempty"          json:"policyId,omitempty"          bigquery:"-"`
	PolicyUid         string                `firestore:"policyUid,omitempty"         json:"policyUid,omitempty"         bigquery:"policyUid"`
	PolicyNumber      string                `firestore:"policyNumber,omitempty"      json:"policyNumber,omitempty"      bigquery:"-"`
	CreationDate      time.Time             `firestore:"creationDate,omitempty"      json:"creationDate,omitempty"      bigquery:"-"`
	BigCreationDate   bigquery.NullDateTime `firestore:"-"                           json:"-"                           bigquery:"creationDate"`
	Updated           time.Time             `firestore:"updated,omitempty"           json:"updated,omitempty"           bigquery:"-"`
	BigUpdated        bigquery.NullDateTime `firestore:"-"           json:"-"           bigquery:"updated"`
	Company           string                `firestore:"company,omitempty"           json:"company,omitempty"           bigquery:"-"`
	Policy            string                `firestore:"policy,omitempty"            json:"policy,omitempty"            bigquery:"-"`
	Description       string                `firestore:"description,omitempty"       json:"description,omitempty"       bigquery:"description"`
	IdCompany         string                `firestore:"idCompany,omitempty"         json:"idCompany,omitempty"         bigquery:"-"`
	UserUid           string                `firestore:"userUid,omitempty"           json:"userUid,omitempty"           bigquery:"userUid"`
	ClaimUid          string                `firestore:"claimUid,omitempty"          json:"claimUid,omitempty"          bigquery:"uid"`
	Status            string                `firestore:"status,omitempty"            json:"status,omitempty"            bigquery:"status"`
	StatusHistory     []string              `firestore:"statusHistory,omitempty"     json:"statusHistory,omitempty"     bigquery:"-"`
	BigStatusHistory  string                `firestore:"-"                           json:"-"                           bigquery:"statusHistory"`
	Documents         []Attachment          `firestore:"documents,omitempty"         json:"documents,omitempty"         bigquery:"-"`
	History           []Claim               `firestore:"history,omitempty"           json:"history,omitempty"           bigquery:"-"`
	Data              string                `firestore:"-"                           json:"-"                           bigquery:"data"`
}

func (claim *Claim) BigquerySave(origin string) error {
	claim.BigCreationDate = lib.GetBigQueryNullDateTime(claim.CreationDate)
	claim.BigDate = lib.GetBigQueryNullDateTime(claim.Date)
	claim.BigUpdated = lib.GetBigQueryNullDateTime(claim.Updated)
	claim.BigStatusHistory = strings.Join(claim.StatusHistory, ",")
	data, err := json.Marshal(claim)
	if err != nil {
		return err
	}
	claim.Data = string(data)

	log.Println("claim save big query: " + claim.ClaimUid)
	table := lib.GetDatasetByEnv(origin, ClaimsCollection)

	return lib.InsertRowsBigQuery(WoptaDataset, table, claim)
}

type Attachment struct {
	Name        string `firestore:"name,omitempty"        json:"name,omitempty"`
	Link        string `firestore:"link,omitempty"        json:"link,omitempty"`
	Byte        string `firestore:"byte,omitempty"        json:"byte,omitempty"`
	FileName    string `firestore:"fileName,omitempty"    json:"fileName,omitempty"`
	MimeType    string `firestore:"mimeType,omitempty"    json:"mimeType,omitempty"`
	Url         string `firestore:"url,omitempty"         json:"url,omitempty"`
	ContentType string `firestore:"contentType,omitempty" json:"contentType,omitempty"`
}
