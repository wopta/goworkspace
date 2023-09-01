package models

type ReservedInfo struct {
	RequiredExams []string     `json:"requiredExams,omitempty" firestore:"requiredExams,omitempty" bigquery:"-"`
	Contacts      []Contact    `json:"contacts,omitempty" firestore:"contacts,omitempty" bigquery:"-"`
	Documents     []Attachment `json:"documents,omitempty" firestore:"documents,omitempty" bigquery:"-"`
	Reasons       []string     `json:"reasons,omitempty" firestore:"reasons,omitempty" bigquery:"-"`
}

type Contact struct {
	Title   string `json:"title,omitempty" firestore:"title,omitempty"`
	Type    string `json:"type,omitempty" firestore:"type,omitempty"`
	Address string `json:"address,omitempty" firestore:"address,omitempty"`
	Subject string `json:"subject,omitempty" firestore:"subject,omitempty"`
}
