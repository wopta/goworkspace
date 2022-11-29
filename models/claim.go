package models

import (
	"encoding/json"
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
	Surname           string       `firestore:"surname" json:"surname,omitempty"`
	Mail              string       `firestore:"mail" json:"mail,omitempty"`
	PolicyDescription string       `firestore:"policyDescription,omitempty" json:"policyDescription,omitempty"`
	PolicyId          string       `firestore:"policyId,omitempty" json:"policyId,omitempty"`
	PolicyUid         string       `firestore:"policyUid,omitempty" json:"policyUid,omitempty"`
	PolicyNumber      string       `firestore:"policyNumber,omitempty" json:"policyNumber,omitempty"`
	CreationDate      string       `firestore:"creationDate,omitempty" json:"creationDate,omitempty"`
	Updated           string       `firestore:"updated,omitempty" json:"updated,omitempty"`
	Company           string       `firestore:"company,omitempty" json:"company,omitempty"`
	Policy            string       `firestore:"policy,omitempty" json:"policy,omitempty"`
	Description       string       `firestore:"description,omitempty" json:"description,omitempty"`
	IdCompany         string       `firestore:"idCompany,omitempty" json:"idCompany,omitempty"`
	Uid               string       `firestore:"uid,omitempty" json:"uid,omitempty"`
	ClaimUid          string       `firestore:"claimUid,omitempty" json:"claimUid,omitempty"`
	Status            string       `firestore:"status,omitempty" json:"status,omitempty"`
	Documents         []Attachment `firestore:"documents,omitempty" json:"documents,omitempty"`
	History           []Claim      `firestore:"history,omitempty" json:"history,omitempty"`
}

type Attachment struct {
	Name     string `json:"name,omitempty" json:"name,omitempty"`
	Link     string `json:"link,omitempty" json:"link,omitempty"`
	Byte     string `json:"byte,omitempty" json:"byte,omitempty"`
	FileName string `json:"fileName,omitempty" json:"fileName,omitempty"`
	MimeType string `json:"mimeType,omitempty" json:"mimeType,omitempty"`
	Url      string `json:"url,omitempty" json:"url,omitempty"`
}
