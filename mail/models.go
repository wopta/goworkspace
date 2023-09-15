package mail

import m "net/mail"

type Data struct {
	Title     string
	SubTitle  string
	Content   string
	Link      string
	LinkLabel string
	IsLink    bool
	IsApp     bool
}
type BodyData struct {
	ContractorName     string
	ContractorSurname  string
	AgentName          string
	AgentSurname       string
	AgentMail          string
	AgencyName         string
	AgencyMail         string
	ProductForm        string
	ProductName        string
	InformationSetsUrl string
	ProposalNumber     int
	ExtraContent       []string
}
type Attachment struct {
	Name        string `firestore:"name,omitempty" json:"name,omitempty"`
	Link        string `firestore:"link,omitempty" json:"link,omitempty"`
	Byte        string `firestore:"byte,omitempty" json:"byte,omitempty"`
	FileName    string `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	MimeType    string `firestore:"mimeType,omitempty" json:"mimeType,omitempty"`
	Url         string `firestore:"url,omitempty" json:"url,omitempty"`
	ContentType string `firestore:"contentType,omitempty" json:"contentType,omitempty"`
}
type MailRequest struct {
	From         string        `json:"from"`
	FromName     string        `json:"fromName"`
	FromAddress  Address       `json:"fromAddress"`
	To           []string      `json:"to"`
	Message      string        `json:"message"`
	Subject      string        `json:"subject"`
	IsHtml       bool          `json:"isHtml,omitempty"`
	IsAttachment bool          `json:"isAttachment,omitempty"`
	Attachments  *[]Attachment `json:"attachments,omitempty"`
	Cc           string        `json:"cc,omitempty"`
	TemplateName string        `json:"templateName,omitempty"`
	Title        string        `json:"title,omitempty"`
	SubTitle     string        `json:"subTitle,omitempty"`
	Content      string        `json:"content,omitempty"`
	Link         string        `json:"link,omitempty"`
	LinkLabel    string        `json:"linkLabel,omitempty"`
	IsLink       bool          `json:"isLink,omitempty"`
	IsApp        bool          `json:"isApp,omitempty"`
}
type MailValidate struct {
	Mail      string `firestore:"mail,omitempty" json:"mail,omitempty"`
	IsValid   bool   `firestore:"isValid" json:"isValid"`
	IsValidS  bool   `firestore:"-" json:"isValid "`
	FidoScore int64  `firestore:"fidoScore" json:"fidoScore"`
}

type Address = m.Address

var (
	AddressAnna = Address{
		Name:    "Anna di Wopta Assicurazioni",
		Address: "anna@wopta.it",
	}
	AddressAssunzione = Address{
		Name:    "Assunzione",
		Address: "assunzione@wopta.it",
	}
)
