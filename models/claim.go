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
	Name              string       `firestore:"name"                        json:"name,omitempty"`
	Date              time.Time    `firestore:"date"                        json:"date,omitempty"`
	Surname           string       `firestore:"surname"                     json:"surname,omitempty"`
	Mail              string       `firestore:"mail"                        json:"mail,omitempty"`
	PolicyDescription string       `firestore:"policyDescription,omitempty" json:"policyDescription,omitempty"`
	PolicyId          string       `firestore:"policyId,omitempty"          json:"policyId,omitempty"`
	PolicyUid         string       `firestore:"policyUid,omitempty"         json:"policyUid,omitempty"`
	PolicyNumber      string       `firestore:"policyNumber,omitempty"      json:"policyNumber,omitempty"`
	CreationDate      time.Time    `firestore:"creationDate,omitempty"      json:"creationDate,omitempty"`
	Updated           time.Time    `firestore:"updated,omitempty"           json:"updated,omitempty"`
	Company           string       `firestore:"company,omitempty"           json:"company,omitempty"`
	Policy            string       `firestore:"policy,omitempty"            json:"policy,omitempty"`
	Description       string       `firestore:"description,omitempty"       json:"description,omitempty"`
	IdCompany         string       `firestore:"idCompany,omitempty"         json:"idCompany,omitempty"`
	UserUid           string       `firestore:"userUid,omitempty"           json:"userUid,omitempty"`
	ClaimUid          string       `firestore:"claimUid,omitempty"          json:"claimUid,omitempty"`
	Status            string       `firestore:"status,omitempty"            json:"status,omitempty"`
	StatusHistory     []string     `firestore:"statusHistory,omitempty"     json:"statusHistory,omitempty"`
	Documents         []Attachment `firestore:"documents,omitempty"         json:"documents,omitempty"`
	History           []Claim      `firestore:"history,omitempty"           json:"history,omitempty"`
}

type ClaimBigquery struct {
	Uid           string                `bigquery:"uid"`
	PolicyUid     string                `bigquery:"policyUid"`
	UserUid       string                `bigquery:"userUid"`
	Status        string                `bigquery:"status"`
	StatusHistory string                `bigquery:"statusHistory"`
	Description   string                `bigquery:"description"`
	Data          string                `bigquery:"data"`
	CreationDate  bigquery.NullDateTime `bigquery:"creationDate"`
}

func (claim Claim) toBigquery() (ClaimBigquery, error) {
	claimData, err := json.Marshal(claim)
	if err != nil {
		return ClaimBigquery{}, err
	}
	return ClaimBigquery{
		Uid:           claim.ClaimUid,
		PolicyUid:     claim.PolicyId,
		UserUid:       claim.UserUid,
		Status:        claim.Status,
		StatusHistory: strings.Join(claim.StatusHistory, ","),
		Description:   claim.Description,
		Data:          string(claimData),
		CreationDate:  lib.GetBigQueryNullDateTime(claim.CreationDate),
	}, nil
}

func (claim Claim) BigquerySave(origin string) error {
	table := lib.GetDatasetByEnv(origin, "claim")
	claimBigquery, err := claim.toBigquery()
	if err != nil {
		return err
	}

	log.Println("claim save big query: " + claim.ClaimUid)

	return lib.InsertRowsBigQuery("wopta", table, claimBigquery)
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
