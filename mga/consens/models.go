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
	errNetworkNodeNotFound     = errors.New("network node not found")
	errRuiSectionNotSet        = errors.New("node does not have rui section")
	errStrategyNotFound        = errors.New("strategy not found")
	errInvalidRequest          = errors.New("invalid request body")
	errConsensExpired          = errors.New("consens already expired")
	errInvalidConsentValue     = errors.New("invalid consent value")
	errInvalidConsensToBeGiven = errors.New("invalid consens to be given")
)

const (
	allProducts = "all"
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

func (c SystemConsens) ToOutput() OutputConsens {
	return OutputConsens{
		Slug:    c.Slug,
		Title:   c.Title,
		Content: c.Content,
	}
}

type ConsensResp struct {
	Consens []OutputConsens `json:"consens"`
}

type OutputConsens struct {
	Slug    string           `json:"slug"`
	Title   string           `json:"title"`
	Content []ConsensContent `json:"content"`
}

type NodeConsensAudit struct {
	Uid             string            `json:"uid" firestore:"uid"`
	NetworkNodeUid  string            `json:"networkNodeUid" firestore:"networkNodeUid"`
	Name            string            `json:"name" firestore:"name"`
	RuiCode         string            `json:"ruiCode" firestore:"ruiCode"`
	RuiRegistration time.Time         `json:"ruiRegistration" firestore:"ruiRegistration"`
	FiscalCode      string            `json:"fiscalCode" firestore:"fiscalCode"`
	VatCode         string            `json:"vatCode" firestore:"vatCode"`
	Slug            string            `json:"slug" firestore:"slug"`
	Title           string            `json:"title" firestore:"title"`
	Content         string            `json:"content" firestore:"content"`
	Answers         map[string]string `json:"answers" firestore:"answers"`
	GivenAt         time.Time         `json:"givenAt" firestore:"givenAt"`
}

func (c *NodeConsensAudit) Save() error {
	c.Uid = lib.NewDoc(lib.NodeConsensAuditsCollencion)
	if err := lib.SetFirestoreErr(lib.NodeConsensAuditsCollencion, c.Uid, c); err != nil {
		return err
	}

	return lib.InsertRowsBigQuery(lib.WoptaDataset, lib.NodeConsensAuditsCollencion, c.BigQueryParse())
}

func (c *NodeConsensAudit) BigQueryParse() NodeConsensAuditBQ {
	answers := make([]BigQueryConsens, 0, len(c.Answers))
	for key, value := range c.Answers {
		answers = append(answers, BigQueryConsens{
			Key:   key,
			Value: value,
		})
	}

	return NodeConsensAuditBQ{
		Uid:             c.Uid,
		NetworkNodeUid:  c.NetworkNodeUid,
		Name:            c.Name,
		RuiCode:         c.RuiCode,
		RuiRegistration: lib.GetBigQueryNullDateTime(c.RuiRegistration),
		FiscalCode:      c.FiscalCode,
		VatCode:         c.VatCode,
		Slug:            c.Slug,
		Title:           c.Title,
		Content:         c.Content,
		Answers:         answers,
		GivenAt:         lib.GetBigQueryNullDateTime(c.GivenAt),
	}
}

type BigQueryConsens struct {
	Key   string `bigquery:"key"`
	Value string `bigquery:"value"`
}

type NodeConsensAuditBQ struct {
	Uid             string                `bigquery:"uid"`
	NetworkNodeUid  string                `bigquery:"networkNodeUid"`
	Name            string                `bigquery:"name"`
	RuiCode         string                `bigquery:"ruiCode"`
	RuiRegistration bigquery.NullDateTime `bigquery:"ruiRegistration"`
	FiscalCode      string                `bigquery:"fiscalCode"`
	VatCode         string                `bigquery:"vatCode"`
	Slug            string                `bigquery:"slug"`
	Title           string                `bigquery:"title"`
	Content         string                `bigquery:"content"`
	Answers         []BigQueryConsens     `bigquery:"answers"`
	GivenAt         bigquery.NullDateTime `bigquery:"givenAt"`
}
