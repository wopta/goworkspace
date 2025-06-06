package mail

import (
	m "net/mail"

	"cloud.google.com/go/bigquery"
	"gitlab.dev.wopta.it/goworkspace/models"
)

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
	ContractorName       string
	ContractorFiscalCode string
	NetworkNodeEmail     string
	NetworkNodeName      string
	ProductForm          string
	ProductName          string
	ProductSlug          string
	SignUrl              string
	PayUrl               string
	PaymentMode          string
	InformationSetsUrl   string
	ProposalNumber       int
	ExtraContent         []string
	RenewDate            string
	PriceGross           string
	HasMandate           bool
	PolicyUid            string
}

type MailRequest struct {
	From         string               `json:"from"`
	FromName     string               `json:"fromName"`
	FromAddress  Address              `json:"fromAddress"`
	To           []string             `json:"to"`
	Message      string               `json:"message"`
	Subject      string               `json:"subject"`
	IsHtml       bool                 `json:"isHtml,omitempty"`
	IsAttachment bool                 `json:"isAttachment,omitempty"`
	Attachments  *[]models.Attachment `json:"attachments,omitempty"`
	Cc           string               `json:"cc,omitempty"`
	Bcc          string               `json:"bcc,omitempty"`
	TemplateName string               `json:"templateName,omitempty"`
	Title        string               `json:"title,omitempty"`
	SubTitle     string               `json:"subTitle,omitempty"`
	Content      string               `json:"content,omitempty"`
	Link         string               `json:"link,omitempty"`
	LinkLabel    string               `json:"linkLabel,omitempty"`
	IsLink       bool                 `json:"isLink,omitempty"`
	IsApp        bool                 `json:"isApp,omitempty"`
	Policy       string               `json:"policy,omitempty"`
}

type MailValidate struct {
	Mail      string `firestore:"mail,omitempty" json:"mail,omitempty"`
	IsValid   bool   `firestore:"isValid" json:"isValid"`
	IsValidS  bool   `firestore:"-" json:"isValid "`
	FidoScore int64  `firestore:"fidoScore" json:"fidoScore"`
}

type MailReport struct {
	Policy         string                `bigquery:"policyUid"`
	SenderName     string                `bigquery:"senderName"`
	RecipientEmail string                `bigquery:"recipientEmail"`
	CreationDate   bigquery.NullDateTime `bigquery:"creationDate"`
	MailError      string                `bigquery:"mailError"`
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
	AddressTechnology = Address{
		Name:    "Technology",
		Address: "technology@wopta.it",
	}
	AddressOperations = Address{
		Name:    "Processi",
		Address: "processi@wopta.it",
	}
)
