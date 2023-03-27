package models

import (
	"encoding/json"
	"time"
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
	Name              string       `firestore:"name" json:"name,omitempty"`
	Date              time.Time    `firestore:"date" json:"date,omitempty"`
	Surname           string       `firestore:"surname" json:"surname,omitempty"`
	Mail              string       `firestore:"mail" json:"mail,omitempty"`
	PolicyDescription string       `firestore:"policyDescription,omitempty" json:"policyDescription,omitempty"`
	PolicyId          string       `firestore:"policyId,omitempty" json:"policyId,omitempty"`
	PolicyUid         string       `firestore:"policyUid,omitempty" json:"policyUid,omitempty"`
	PolicyNumber      string       `firestore:"policyNumber,omitempty" json:"policyNumber,omitempty"`
	CreationDate      time.Time    `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	Updated           time.Time    `firestore:"updated,omitempty" json:"updated,omitempty"`
	Company           string       `firestore:"company,omitempty" json:"company,omitempty"`
	Policy            string       `firestore:"policy,omitempty" json:"policy,omitempty"`
	Description       string       `firestore:"description,omitempty" json:"description,omitempty"`
	IdCompany         string       `firestore:"idCompany,omitempty" json:"idCompany,omitempty"`
	UserUid           string       `firestore:"userUid,omitempty" json:"userUid,omitempty"`
	ClaimUid          string       `firestore:"claimUid,omitempty" json:"claimUid,omitempty"`
	Status            string       `firestore:"status,omitempty" json:"status,omitempty"`
	Documents         []Attachment `firestore:"documents,omitempty" json:"documents,omitempty"`
	History           []Claim      `firestore:"history,omitempty" json:"history,omitempty"`
}

type Attachment struct {
	Name        string `firestore:"name,omitempty" json:"name,omitempty"`
	Link        string `firestore:"link,omitempty" json:"link,omitempty"`
	Byte        string `firestore:"byte,omitempty" json:"byte,omitempty"`
	FileName    string `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	MimeType    string `firestore:"mimeType,omitempty" json:"mimeType,omitempty"`
	Url         string `firestore:"url,omitempty" json:"url,omitempty"`
	ContentType string `firestore:"contentType,omitempty" json:"contentType,omitempty"`
}
