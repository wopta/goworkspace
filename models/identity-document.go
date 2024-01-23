package models

import (
	"github.com/wopta/goworkspace/lib"
	"time"
)

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
	SourceFileName string `json:"sourceFileName" firestore:"sourceFileName" bigquery:"-"`
	FileName       string `json:"fileName" firestore:"fileName" bigquery:"-"`
	Link           string `json:"link,omitempty" firestore:"link,omitempty" bigquery:"-"`
	MimeType       string `json:"mimeType,omitempty" firestore:"mimeType,omitempty" bigquery:"-"`
	Base64Bytes    string `json:"base64Bytes,omitempty" firestore:"-" bigquery:"-"`
}

func (id *IdentityDocument) Sanitize() {
	id.Code = lib.TrimSpace(id.Code)
	id.Type = lib.TrimSpace(id.Type)
	id.Number = lib.TrimSpace(id.Number)
	id.IssuingAuthority = lib.TrimSpace(id.IssuingAuthority)
	id.PlaceOfIssue = lib.ToUpper(id.PlaceOfIssue)
	if id.FrontMedia != nil {
		id.FrontMedia.Sanitize()
	}
	if id.BackMedia != nil {
		id.BackMedia.Sanitize()
	}
}

func (id *IdentityDocument) IsExpired() bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return id.ExpiryDate.UTC().Before(today)
}

func (m *Media) Sanitize() {
	m.FileName = lib.TrimSpace(m.FileName)
	m.SourceFileName = lib.TrimSpace(m.SourceFileName)
	m.MimeType = lib.TrimSpace(m.MimeType)
}
