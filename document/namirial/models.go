package namirial

import (
	"gitlab.dev.wopta.it/goworkspace/document"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type NamirialInput struct {
	FilesFullPath []string
	Policy        models.Policy
	SendEmail     bool
	Origin        string
}

type NamirialOutput struct {
	Url        string
	IdEnvelope string
	FileIds    []string
}

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
	Documents                  []documentDescription      `json:"Documents"`
	Name                       string                     `json:"Name"`
	AddDocumentTimestamp       bool                       `json:"AddDocumentTimestamp"`
	ShareWithTeam              bool                       `json:"ShareWithTeam"`
	LockFormFieldsOnFinish     bool                       `json:"LockFormFieldsOnFinish"`
	Activities                 []document.Activity        `json:"Activities"`
	UnassignedElements         document.Elements          `json:"UnassignedElements"`
	CallbackConfiguration      callbackConfiguration      `json:"CallbackConfiguration"`
	AgentRedirectConfiguration agentRedirectConfiguration `json:"AgentRedirectConfiguration"`
	ReminderConfiguration      reminderConfiguration      `json:"ReminderConfiguration"`
}
type reminderConfiguration struct {
	Enabled                      bool `json:"Enabled"`
	FirstReminderInDays          int  `json:"FirstReminderInDays"`
	ReminderResendIntervalInDays int  `json:"ReminderResendIntervalInDays"`
	BeforeExpirationInDays       int  `json:"BeforeExpirationInDays"`
}
type agentRedirectConfiguration struct {
	Policy             string   `json:"Policy"`
	Allow              bool     `json:"Allow"`
	IframeWhitelisting []string `json:"IframeWhitelisting"`
}

type callbackConfiguration struct {
	CallbackUrl                  string                              `json:"CallbackUrl"`
	StatusUpdateCallbackUrl      string                              `json:"StatusUpdateCallbackUrl"`
	ActivityActionCallbackConfig activityActionCallbackConfiguration `json:"ActivityActionCallbackConfiguration"`
}

type activityActionCallbackConfiguration struct {
	Url                     string                  `json:"Url"`
	ActionCallbackSelection actionCallbackSelection `json:"ActionCallbackSelection"`
}

type actionCallbackSelection struct {
	ConfirmTransactionCode         bool `json:"ConfirmTransactionCode"`
	AgreementAccepted              bool `json:"AgreementAccepted"`
	AgreementRejected              bool `json:"AgreementRejected"`
	PrepareAuthenticationSuccess   bool `json:"PrepareAuthenticationSuccess"`
	AuthenticationFailed           bool `json:"AuthenticationFailed"`
	AuthenticationSuccess          bool `json:"AuthenticationSuccess"`
	AuditTrailRequested            bool `json:"AuditTrailRequested"`
	AuditTrailXmlRequested         bool `json:"AuditTrailXmlRequested"`
	CalledPage                     bool `json:"CalledPage"`
	DocumentDownloaded             bool `json:"DocumentDownloaded"`
	FlattenedDocumentDownloaded    bool `json:"FlattenedDocumentDownloaded"`
	AddedAnnotation                bool `json:"AddedAnnotation"`
	AddedAttachment                bool `json:"AddedAttachment"`
	AppendedDocument               bool `json:"AppendedDocument"`
	FormsFilled                    bool `json:"FormsFilled"`
	ConfirmReading                 bool `json:"ConfirmReading"`
	SendTransactionCode            bool `json:"SendTransactionCode"`
	PrepareSignWorkstepDocument    bool `json:"PrepareSignWorkstepDocument"`
	SignWorkstepDocument           bool `json:"SignWorkstepDocument"`
	UndoAction                     bool `json:"UndoAction"`
	WorkstepCreated                bool `json:"WorkstepCreated"`
	WorkstepFinished               bool `json:"WorkstepFinished"`
	WorkstepRejected               bool `json:"WorkstepRejected"`
	DisablePolicyAndValidityChecks bool `json:"DisablePolicyAndValidityChecks"`
	EnablePolicyAndValidityChecks  bool `json:"EnablePolicyAndValidityChecks"`
	AppendFileToWorkstep           bool `json:"AppendFileToWorkstep"`
	AppendTasksToWorkstep          bool `json:"AppendTasksToWorkstep"`
	SetOptionalDocumentState       bool `json:"SetOptionalDocumentState"`
	PreparePayloadForBatch         bool `json:"PreparePayloadForBatch"`
}

type responseSendDocuments struct {
	EnvelopeId string `json:"EnvelopeId"`
}

type responseGetEvelop struct {
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
