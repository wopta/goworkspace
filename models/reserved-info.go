package models

type ReservedInfo struct {
	MedicalDocuments  []string     `json:"medicalDocuments,omitempty" firestore:"medicalDocuments,omitempty" bigquery:"-"`
	Contacts          []Contact    `json:"contacts,omitempty" firestore:"medicalDocuments,omitempty" bigquery:"-"`
	DownloadDocuments []Attachment `json:"downloadDocuments,omitempty" firestore:"downloadDocuments,omitempty" bigquery:"-"`
	Reasons           []string     `json:"reasons,omitempty" firestore:"reasons,omitempty" bigquery:"-"`
	BigReasons        string       `json:"-" firestore:"-" bigquery:"reasons"`
}

type Contact struct {
	ContactType string `json:"contactType"`
	Address     string `json:"address"`
	Object      string `json:"object,omitempty"`
}
