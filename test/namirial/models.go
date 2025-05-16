package namirial

import (
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/models"
)

type StatusNamirial string

const (
	Idle        StatusNamirial = "Idle"
	Upload      StatusNamirial = "Uploaded Files"
	Prepared    StatusNamirial = "Prepared Files"
	Sended      StatusNamirial = "Sended Files"
	GetEnvelope StatusNamirial = "Get Envelope"
)

type dataForDocument struct {
	policy     *models.Policy
	product    *models.Product
	warrant    *models.Warrant
	idDocument string
}

type prepareNamirialDocumentRequest struct {
	FileIds                   []string                 `json:"FileIds"`
	ClearAdvancedDocumentTags bool                     `json:"ClearAdvancedDocumentTags"`
	SigStringConfigurations   []sigStringConfiguration `json:"SigStringConfigurations"`
}

type sigStringConfiguration struct {
	StartPattern         string `json:"StartPattern"`
	EndPattern           string `json:"EndPattern"`
	ClearSigString       bool   `json:"ClearSigString"`
	SearchEntireWordOnly bool   `json:"SearchEntireWordOnly"`
}

type documentDescription struct {
	FileId         string `json:"FileId"`
	DocumentNumber int    `json:"DocumentNumber"`
}

type sendNamirialRequest struct {
	Documents  []documentDescription `json:"Documents"`
	Name       string                `json:"Name"`
	Activities []document.Activity   `json:"Activities"`
}

type responseSendDocuments struct {
	EnvelopeId string `json:"EnvelopeId"`
}

type ResponeGetEvelop struct {
	ViewerLinks []viewerLink `json:"ViewerLinks"`
}

type viewerLink struct {
	ActivityId string `json:"ActivityId"`
	Email      string `json:"Email"`
	ViewerLink string `json:"ViewerLink"`
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
	Documents   []documentDesc `json:"Documents"`
	AuditTrail  auditTrail     `json:"AuditTrail"`
	Disclaimers []string       `json:"Disclaimers"`
}
