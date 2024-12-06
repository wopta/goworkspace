package consens

import (
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/wopta/goworkspace/lib"
)

var (
	errNetworkNodeNotFound = errors.New("network node not found")
	errRuiSectionNotSet    = errors.New("node does not have rui section")
	errStrategyNotFound    = errors.New("strategy not found")
)

const (
	ruiSectionE = "E"
	folderPath  = "consens/network/"
)

type SystemConsens struct {
	Slug        string           `json:"slug"`
	ExpireAt    time.Time        `json:"expireAt"`
	StartAt     time.Time        `json:"startAt"`
	AvailableAt time.Time        `json:"availableAt"`
	Strategy    string           `json:"strategy"`
	Title       string           `json:"title"`
	Content     []ConsensContent `json:"content"`
}

type ConsensContent struct {
	Text       string `json:"text"`
	InputType  string `json:"inputType,omitempty"`
	InputName  string `json:"inputName,omitempty"`
	InputValue string `json:"inputValue,omitempty"`
}

func (c SystemConsens) ToString() string {
	var parts []string

	for _, cs := range c.Content {
		// TODO: regex to remove markdown
		parts = append(parts, cs.Text)
	}

	return strings.Join(parts, " ")
}

type ConsensResp struct {
	Consens []SystemConsens `json:"consens"`
}

type NodeConsensAudit struct {
	Name            string
	RuiCode         string
	RuiRegistration time.Time `json:"ruiRegistration" firestore:"ruiRegistration" bigquery:"-"`
	FiscalCode      string
	VatCode         string
	Slug            string
	Title           string
	Content         string
	Answer          string
	GivenAt         time.Time
}

func (c *NodeConsensAudit) Save() error {
	if err := lib.SetFirestoreErr("", "", c); err != nil {
		return err
	}

	return lib.InsertRowsBigQuery("", "", c.BigQueryParse())
}

func (c *NodeConsensAudit) BigQueryParse() NodeConsensAuditBQ {
	return NodeConsensAuditBQ{
		Name:            c.Name,
		RuiCode:         c.RuiCode,
		RuiRegistration: lib.GetBigQueryNullDateTime(c.RuiRegistration),
		FiscalCode:      c.FiscalCode,
		VatCode:         c.VatCode,
		Slug:            c.Slug,
		Title:           c.Title,
		Content:         c.Content,
		Answer:          c.Answer,
		GivenAt:         lib.GetBigQueryNullDateTime(c.GivenAt),
	}
}

type NodeConsensAuditBQ struct {
	Name            string
	RuiCode         string
	RuiRegistration bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"ruiRegistration"`
	FiscalCode      string
	VatCode         string
	Slug            string
	Title           string
	Content         string
	Answer          string
	GivenAt         bigquery.NullDateTime
}
