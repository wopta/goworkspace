package consens

import (
	"errors"
	"regexp"
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

	// TODO: this is an oversimplification
	markdown := regexp.MustCompile("[*_~`#]")

	for _, cs := range c.Content {
		text := markdown.ReplaceAll([]byte(cs.Text), []byte(""))
		parts = append(parts, string(text))
	}

	return strings.Join(parts, "\n")
}

type ConsensResp struct {
	Consens []SystemConsens `json:"consens"`
}

type NodeConsensAudit struct {
	Uid             string    `json:"uid" firestore:"uid"`
	Name            string    `json:"name" firestore:"name"`
	RuiCode         string    `json:"ruiCode" firestore:"ruiCode"`
	RuiRegistration time.Time `json:"ruiRegistration" firestore:"ruiRegistration"`
	FiscalCode      string    `json:"fiscalCode" firestore:"fiscalCode"`
	VatCode         string    `json:"vatCode" firestore:"vatCode"`
	Slug            string    `json:"slug" firestore:"slug"`
	Title           string    `json:"title" firestore:"title"`
	Content         string    `json:"content" firestore:"content"`
	Answer          string    `json:"answer" firestore:"answer"`
	GivenAt         time.Time `json:"givenAt" firestore:"givenAt"`
}

func (c *NodeConsensAudit) Save() error {
	c.Uid = lib.NewDoc(lib.NodeConsensAuditCollencion)
	if err := lib.SetFirestoreErr(lib.NodeConsensAuditCollencion, c.Uid, c); err != nil {
		return err
	}

	return lib.InsertRowsBigQuery(lib.WoptaDataset, lib.NodeConsensAuditCollencion, c.BigQueryParse())
}

func (c *NodeConsensAudit) BigQueryParse() NodeConsensAuditBQ {
	return NodeConsensAuditBQ{
		Uid:             c.Uid,
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
	Uid             string                `bigquery:"uid"`
	Name            string                `bigquery:"name"`
	RuiCode         string                `bigquery:"ruiCode"`
	RuiRegistration bigquery.NullDateTime `bigquery:"ruiRegistration"`
	FiscalCode      string                `bigquery:"fiscalCode"`
	VatCode         string                `bigquery:"vatCode"`
	Slug            string                `bigquery:"slug"`
	Title           string                `bigquery:"title"`
	Content         string                `bigquery:"content"`
	Answer          string                `bigquery:"answer"`
	GivenAt         bigquery.NullDateTime `bigquery:"givenAt"`
}
