package models

type Reserved struct {
	MedicalDocuments  []string     `json:"medicalDocuments"`
	Contacts          []Contact    `json:"contacts"`
	DownloadDocuments []Attachment `json:"downloadDocuments"`
}

type Contact struct {
	ContactType string `json:"contactType"`
	Address     string `json:"address"`
	Object      string `json:"object,omitempty"`
}
