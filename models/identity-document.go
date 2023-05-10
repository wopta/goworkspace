package models

import "time"

type IdentityDocument struct {
	Code             string    `json:"code" firestore:"code" bigquery:"-"`
	Type             string    `json:"type" firestore:"type" bigquery:"-"`
	Number           string    `json:"number" firestore:"number" bigquery:"-"`
	IssuingAuthority string    `json:"issuingAuthority" firestore:"issuingAuthority" bigquery:"-"`
	PlaceOfIssue     string    `json:"placeOfIssue" firestore:"placeOfIssue" bigquery:"-"`
	DateOfIssue      time.Time `json:"dateOfIssue" firestore:"dateOfIssue" bigquery:"-"`
	ExpiryDate       time.Time `json:"expiryDate" firestore:"expiryDate" bigquery:"-"`
	FrontMedia       *Media    `json:"frontMedia" firestore:"frontMedia" bigquery:"-"`
	BackMedia        *Media    `json:"backMedia,omitempty" firestore:"backMedia" bigquery:"-"`
	LastUpdate       time.Time `json:"lastUpdate,omitempty" firestore:"lastUpdate,omitempty" bigquery:"-"`
}

type Media struct {
	Filename    string `json:"filename" firestore:"filename" bigquery:"-"`
	Link        string `json:"link,omitempty" firestore:"link,omitempty" bigquery:"-"`
	MimeType    string `json:"mimeType,omitempty" firestore:"mimeType,omitempty" bigquery:"-"`
	Base64Bytes string `json:"base64Bytes,omitempty" firestore:"-" bigquery:"-"`
}

func (id *IdentityDocument) IsExpired() bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return id.ExpiryDate.UTC().Before(today)
}
