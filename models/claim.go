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
	Status            string       `firestore:"status,omitempty" json:"status,omitempty"`
	Document          []Attachment `firestore:"attachment,omitempty" json:"attachment,omitempty"`
	History           []Claim      `firestore:"history,omitempty" json:"history,omitempty"`
}

type Attachment struct {
	Name *string `json:"name,omitempty" json:"name,omitempty"`
	Link *string `json:"link,omitempty" json:"link,omitempty"`
	Byte *string `json:"byte,omitempty" json:"byte,omitempty"`
}
