package quote

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/wopta/goworkspace/models"
)

type LifeQuoteTestInput struct {
	Age                 int `json:"int"`
	Death               int `json:"death"`
	PermanentDisability int `json:"permanent-disability"`
	TemporaryDisability int `json:"temporary-disability"`
	SeriousIllness      int `json:"serious-ill"`
}

type QuoteLifeSuite struct {
	suite.Suite
	baseStub       models.Policy
	currentStub    models.Policy
	channel        string
	flow           string
	hasNetworkNode bool
	networkNode    *models.NetworkNode
	warrant        *models.Warrant
	inputs         []LifeQuoteTestInput
	outputs        []models.Policy
}

var networkNodeStub models.NetworkNode = models.NetworkNode{}

var warrantStub models.Warrant = models.Warrant{}

func (qls *QuoteLifeSuite) SetupSuite() {
	os.Setenv("env", "local-test")
	var (
		inputs   []LifeQuoteTestInput
		outputs  []models.Policy
		baseStub models.Policy
	)

	// Load all files and inject into suite
	qls.baseStub = baseStub
	qls.inputs = inputs
	qls.outputs = outputs
}

func (qls *QuoteLifeSuite) SetupTest() {
	qls.currentStub = qls.baseStub

	if qls.hasNetworkNode {
		qls.networkNode = &networkNodeStub
		qls.warrant = &warrantStub
	}
}

func (qls *QuoteLifeSuite) TearDownSuite() {
	os.Setenv("env", "local")
}

func (qls *QuoteLifeSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestQuoteLifeEcommerceOffers":
		qls.channel = models.ECommerceChannel
		qls.flow = models.ECommerceFlow
		qls.hasNetworkNode = false
	case "TestQuoteLifeProviderOffers":
		qls.channel = models.NetworkChannel
		qls.flow = models.ProviderMgaFlow
		qls.hasNetworkNode = true
	case "TestQuoteLifeRemittanceOffers":
		qls.channel = models.NetworkChannel
		qls.flow = models.RemittanceMgaFlow
		qls.hasNetworkNode = true
	}
}

func (qls *QuoteLifeSuite) AfterTest(suiteName, testName string) {}

func (qls *QuoteLifeSuite) TestQuoteLifeEcommerceOffers() {
	baseLifeQuoteTest(qls)
}

func (qls *QuoteLifeSuite) TestQuoteLifeProviderOffers() {
	baseLifeQuoteTest(qls)
}

func (qls *QuoteLifeSuite) TestQuoteLifeRemittanceOffers() {
	baseLifeQuoteTest(qls)
}

func baseLifeQuoteTest(qls *QuoteLifeSuite) {
	for index, input := range qls.inputs {
		// inject input in policyStub
		birthDate := time.Now().UTC().AddDate(-input.Age, 0, 0).Format(time.RFC3339)
		qls.currentStub.Contractor.BirthDate = birthDate

		Life(qls.currentStub, qls.channel, qls.networkNode, qls.warrant, qls.flow)

		qls.Assert().Equal(qls.outputs[index], qls.currentStub)
	}
}

func TestLifeQuote(t *testing.T) {
	suite.Run(t, new(QuoteLifeSuite))
}
