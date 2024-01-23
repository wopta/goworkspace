package models

import "github.com/wopta/goworkspace/lib"

type Attachment struct {
	Name        string `firestore:"name,omitempty"        json:"name,omitempty"`
	Link        string `firestore:"link,omitempty"        json:"link,omitempty"`
	Byte        string `firestore:"byte,omitempty"        json:"byte,omitempty"`
	FileName    string `firestore:"fileName,omitempty"    json:"fileName,omitempty"`
	MimeType    string `firestore:"mimeType,omitempty"    json:"mimeType,omitempty"`
	Url         string `firestore:"url,omitempty"         json:"url,omitempty"`
	ContentType string `firestore:"contentType,omitempty" json:"contentType,omitempty"`
	IsPrivate   bool   `firestore:"isPrivate" json:"isPrivate"`
	Section     string `firestore:"section" json:"section"`
	Note        string `firestore:"note,omitempty" json:"note"`
}

func (a *Attachment) Sanitize() {
	a.Name = lib.TrimSpace(a.Name)
	a.FileName = lib.TrimSpace(a.FileName)
	a.Section = lib.TrimSpace(a.Section)
	a.Note = lib.ToUpper(a.Note)
}
