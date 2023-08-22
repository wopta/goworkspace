package models

type ReservedInfo struct {
	RequiredExams []string     `json:"requiredExams,omitempty" firestore:"requiredExams,omitempty" bigquery:"-"`
	Contacts      []Contact    `json:"contacts,omitempty" firestore:"contacts,omitempty" bigquery:"-"`
	Documents     []Attachment `json:"documents,omitempty" firestore:"documents,omitempty" bigquery:"-"`
	Reasons       []string     `json:"reasons,omitempty" firestore:"reasons,omitempty" bigquery:"-"`
	BigReasons    string       `json:"-" firestore:"-" bigquery:"reasons"`
}

type Contact struct {
	Title   string `json:"title,omitempty" firestore:"title,omitempty"`
	Type    string `json:"type,omitempty" firestore:"type,omitempty"`
	Address string `json:"address,omitempty" firestore:"address,omitempty"`
	Object  string `json:"object,omitempty" firestore:"object,omitempty"`
}
