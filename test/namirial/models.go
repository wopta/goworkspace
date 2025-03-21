package namirial

import "github.com/wopta/goworkspace/document"

type prepareNamirialDocumentRequest struct {
	FileIds                     []string                    `json:"FileIds"`
	ClearAdvancedDocumentTags   bool                     `json:"ClearAdvancedDocumentTags"`
	SigStringConfigurations     []sigStringConfiguration `json:"SigStringConfigurations"`
}

type sigStringConfiguration struct {
	StartPattern         string `json:"StartPattern"`
	EndPattern           string `json:"EndPattern"`
	ClearSigString       bool   `json:"ClearSigString"`
	SearchEntireWordOnly bool   `json:"SearchEntireWordOnly"`
}

type documentDescription struct {
	FileId string
	DocumentNumber int
}

type sendNamirialRequest struct {
	Documents [] documentDescription
	Name string
	Activities []document.Activity
}

type responseSendDocuments struct {
	EnvelopeId string
}

type ResponeGetEvelop struct{
	ViewerLinks []viewerLink
}

type viewerLink struct {
	ActivityId  string `json:"ActivityId"`
	Email       string `json:"Email"`
	ViewerLink  string `json:"ViewerLink"`
}

type viewerLinksResponse struct {
	ViewerLinks []viewerLink `json:"ViewerLinks"`
}

type documentDesc struct {
	FileId         string   `json:"FileId"`
	FileName       string   `json:"FileName"`
	Attachments    []string `json:"Attachments"`
	PageCount      int      `json:"PageCount"`
	DocumentNumber int      `json:"DocumentNumber"`
}

type auditTrail struct {
	FileId    string `json:"FileId"`
	XmlFileId string `json:"XmlFileId"`
}

type FilesIdsResponse struct {
	Documents   []documentDesc  `json:"Documents"`
	AuditTrail  auditTrail  `json:"AuditTrail"`
	Disclaimers []string    `json:"Disclaimers"`
}

