package namirial

import (
	"gitlab.dev.wopta.it/goworkspace/models"
)

type NamirialInput struct {
	DocumentsFullPath []string
	Policy            models.Policy
	SendEmail         bool
	Origin            string
}

type NamirialOutput struct {
	Url        string
	IdEnvelope string
	FileIds    []string
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
	Activities                 []activity                 `json:"Activities"`
	UnassignedElements         elements                   `json:"UnassignedElements"`
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

type auditTrail struct {
	FileId    string `json:"FileId"`
	XmlFileId string `json:"XmlFileId"`
}

type prepareResponse struct {
	UnassignedElements elements   `json:"UnassignedElements"`
	Activities         []activity `json:"Activities"`
}

type activity struct {
	Action action `json:"Action"`
}

type action struct {
	Sign sign `json:"Sign"`
}

type sign struct {
	Elements                  elements                  `json:"Elements"`
	RecipientConfiguration    recipientConfiguration    `json:"RecipientConfiguration"`
	FinishActionConfiguration finishActionConfiguration `json:"FinishActionConfiguration"`
}
type finishActionConfiguration struct {
	SignAnyWhereViewer signAnyWhereViewer `json:"SignAnyWhereViewer"`
}

type signAnyWhereViewer struct {
	RedirectUri string `json:"RedirectUri"`
}

type recipientConfiguration struct {
	SendEmails                  bool                        `json:"SendEmails"`
	AllowAccessAfterFinish      bool                        `json:"AllowAccessAfterFinish"`
	AllowDelegation             bool                        `json:"AllowDelegation"`
	ContactInformation          contactInformation          `json:"ContactInformation"`
	PersonalMessage             string                      `json:"PersonalMessage"`
	AuthenticationConfiguration authenticationConfiguration `json:"AuthenticationConfiguration"`
}

type contactInformation struct {
	Email        string `json:"Email"`
	GivenName    string `json:"GivenName"`
	Surname      string `json:"Surname"`
	PhoneNumber  string `json:"PhoneNumber"`
	LanguageCode string `json:"LanguageCode"`
}

type authenticationConfiguration struct {
	SmsOneTimePassword smsOneTimePassword `json:"SmsOneTimePassword"`
	AccessCode         accessCode
}

type accessCode struct {
	Code string
}

type smsOneTimePassword struct {
	PhoneNumber string `json:"PhoneNumber"`
}
type elements struct {
	TextBoxes    []any `json:"TextBoxes"`
	CheckBoxes   []any `json:"CheckBoxes"`
	ComboBoxes   []any `json:"ComboBoxes"`
	RadioButtons []any `json:"RadioButtons"`
	ListBoxes    []any `json:"ListBoxes"`
	Attachments  []any `json:"Attachments"`

	Signatures []signature `json:"Signatures"`
}

type signature struct {
	ElementID             string                `json:"ElementId"`
	Required              bool                  `json:"Required"`
	DocumentNumber        int64                 `json:"DocumentNumber"`
	DisplayName           string                `json:"DisplayName"`
	AllowedSignatureTypes allowedSignatureTypes `json:"AllowedSignatureTypes"`
	FieldDefinition       fieldDefinition       `json:"FieldDefinition"`
	TaskConfiguration     taskConfiguration     `json:"TaskConfiguration"`
}

type allowedSignatureTypes struct {
	ClickToSign      clickToSign `json:"ClickToSign"`
	SignaturePlugins []any       `json:"SignaturePlugins"`
}

type clickToSign struct {
	UseExternalSignatureImage string `json:"UseExternalSignatureImage"`
}

type fieldDefinition struct {
	Position position `json:"Position"`
	Size     size     `json:"Size"`
}

type position struct {
	PageNumber int64   `json:"PageNumber"`
	X          float64 `json:"X"`
	Y          float64 `json:"Y"`
}

type size struct {
	Width  float64 `json:"Width"`
	Height float64 `json:"Height"`
}

type taskConfiguration struct {
	BatchGroup      string          `json:"BatchGroup"`
	OrderDefinition orderDefinition `json:"OrderDefinition"`
}

type orderDefinition struct {
	OrderIndex int64 `json:"OrderIndex"`
}

type namirialFiles struct {
	Documents      []documents     `json:"Documents"`
	AuditTrail     auditTrail      `json:"AuditTrail"`
	LegalDocuments []legalDocument `json:"LegalDocuments"`
}

type documents struct {
	FileID           string       `json:"FileId"`
	FileName         string       `json:"FileName"`
	AuditTrailFileID string       `json:"AuditTrailFileId"`
	Attachments      []attachment `json:"Attachments"`
	PageCount        int64        `json:"PageCount"`
	DocumentNumber   int64        `json:"DocumentNumber"`
}

type attachment struct {
	FileID   string `json:"FileId"`
	FileName string `json:"FileName"`
}

type legalDocument struct {
	FileID     string `json:"FileId"`
	FileName   string `json:"FileName"`
	ActivityID string `json:"ActivityId"`
	Email      string `json:"Email"`
}
