package models

import (
	"time"

	"github.com/wopta/goworkspace/lib"
)

type ReservedInfo struct {
	MgaApproval       StakeholderApproval `json:"mgaApproval" firestore:"mgaApproval" bigquery:"-"`
	CompanyApproval   StakeholderApproval `json:"companyApproval" firestore:"companyApproval" bigquery:"-"`
	ReservedReasons   []ReservedData      `json:"reservedReasons" firestore:"reservedReasons" bigquery:"-"`
	RequiredDocuments []ReservedData      `json:"requiredDocuments" firestore:"requiredDocuments" bigquery:"-"`
	Attachments       []Attachment        `json:"attachments" firestore:"attachments" bigquery:"-"`

	// DEPRECATED FIELDS

	RequiredExams  []string     `json:"requiredExams,omitempty" firestore:"requiredExams,omitempty" bigquery:"-"`   // DEPRECATED - use RequiredDocuments
	Contacts       []Contact    `json:"contacts,omitempty" firestore:"contacts,omitempty" bigquery:"-"`             // DEPRECATED - fixed in the document
	Documents      []Attachment `json:"documents,omitempty" firestore:"documents,omitempty" bigquery:"-"`           // DEPRECATED - use Attachments
	Reasons        []string     `json:"reasons,omitempty" firestore:"reasons,omitempty" bigquery:"-"`               // DEPRECATED - use ReservedReasons
	AcceptanceNote string       `json:"acceptanceNote,omitempty" firestore:"acceptanceNote,omitempty" bigquery:"-"` // DEPRECATED - use the relative inside MgaApproval/CompanyApproval
	AcceptanceDate time.Time    `json:"acceptanceDate,omitempty" firestore:"acceptanceDate,omitempty" bigquery:"-"` // DEPRECATED - use the relative inside MgaApproval/CompanyApproval
}

func (ri *ReservedInfo) Normalize() {
	ri.AcceptanceNote = lib.ToUpper(ri.AcceptanceNote)
	ri.MgaApproval.AcceptanceNotes = lib.SliceMap(ri.MgaApproval.AcceptanceNotes, func(n string) string {
		return lib.ToUpper(n)
	})
	ri.CompanyApproval.AcceptanceNotes = lib.SliceMap(ri.CompanyApproval.AcceptanceNotes, func(n string) string {
		return lib.ToUpper(n)
	})
}

// TODO add tags
type StakeholderApproval struct {
	Mandatory       bool
	Status          ApprovalStatus
	AcceptanceDate  time.Time
	AcceptanceNotes []string
	UpdateDate      time.Time
}

type ApprovalStatus string

const (
	NeedsApproval   ApprovalStatus = "NeedsApproval"
	WaitingApproval ApprovalStatus = "WaitingApproval"
	Approved        ApprovalStatus = "Approved"
	Rejected        ApprovalStatus = "Rejected"
)

// TODO add tags
type ReservedData struct {
	Id          int
	Name        string
	Description string
}

// DEPRECATED
type Contact struct {
	Title   string `json:"title,omitempty" firestore:"title,omitempty"`
	Type    string `json:"type,omitempty" firestore:"type,omitempty"`
	Address string `json:"address,omitempty" firestore:"address,omitempty"`
	Subject string `json:"subject,omitempty" firestore:"subject,omitempty"`
}
