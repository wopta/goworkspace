package quote

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/wopta/goworkspace/models"
)

type LifeModQuoteTestInput struct {
	Age        int `json:"int"`
	Guarantees []struct {
		Name     string `json:"name"`
		Value    int    `json:"value"`
		Duration int    `json:"duration"`
	} `json:"guarantees"`
}

type QuoteLifeModSuite struct {
	suite.Suite
	baseStub models.Policy
	inputs   []LifeModQuoteTestInput
	outputs  []models.Policy
	quoter   Quoter
}

func (qls *QuoteLifeModSuite) SetupSuite() {
	os.Setenv("env", "local-test")
	var (
		inputs   []LifeModQuoteTestInput
		outputs  []models.Policy
		baseStub models.Policy
	)

	// Load all files and inject into suite
	qls.baseStub = baseStub
	qls.inputs = inputs
	qls.outputs = outputs
	qls.quoter = Quoter{
		policy:     baseStub,
		taxesBytes: []byte("18;0,01024%;0,01024%;0,00123%;0,00126%;0,03360%;0,03444%;0,11824%;0,12120%;0,01286%;0,01286%;0,00154%;0,00158%;0,04679%;0,04796%;0,13976%;0,14326%;0,01593%;0,01593%;0,00191%;0,00196%;0,05364%;0,05498%;0,17136%;0,17565%;0,01966%;0,01966%;0,00236%;0,00242%;0,01888%;0,01888%;0,00227%;0,00232%;\n32;0,02919%;0,02919%;0,00350%;0,00359%;0,06732%;0,06900%;0,23845%;0,24441%;0,03666%;0,03666%;0,00440%;0,00451%;0,09404%;0,09639%;0,31221%;0,32002%;0,04539%;0,04539%;0,00545%;0,00558%;0,10822%;0,11092%;0,39021%;0,39996%;0,05604%;0,05604%;0,00672%;0,00689%;0,06896%;0,06896%;0,00828%;0,00848%;\n54;0,19924%;0,19924%;0,02390%;0,02450%;0,21056%;0,21582%;1,61912%;1,65960%;0,25016%;0,25016%;0,03001%;0,03076%;0,29579%;0,30319%;2,24807%;2,30427%;0,30957%;0,30957%;0,03715%;0,03808%;0,34289%;0,35146%;2,86925%;2,94098%;0,38158%;0,38158%;0,04579%;0,04693%;0,46820%;0,46820%;0,05618%;0,05759%;\n65;0,52059%;0,52059%;0,06241%;0,06397%;0,38223%;0,39178%;0,00000%;0,00000%;0,65345%;0,65345%;0,07832%;0,08028%;0,53861%;0,55208%;0,00000%;0,00000%;0,80761%;0,80761%;0,09691%;0,09934%;0,62657%;0,64224%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;\n74;1,14241%;1,14241%;0,13679%;0,14021%;0,63142%;0,64721%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;0,00000%;"),
	}
}

func (qls *QuoteLifeModSuite) SetupTest() {
	qls.quoter.policy = qls.baseStub
}

func (qls *QuoteLifeModSuite) TearDownSuite() {
	os.Setenv("env", "local")
}

func (qls *QuoteLifeModSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestQuoteLifeEcommerceOffers":
		qls.quoter.channel = models.ECommerceChannel
		qls.quoter.flow = models.ECommerceFlow
	case "TestQuoteLifeProviderOffers":
		qls.quoter.channel = models.NetworkChannel
		qls.quoter.flow = models.ProviderMgaFlow
	case "TestQuoteLifeRemittanceOffers":
		qls.quoter.channel = models.NetworkChannel
		qls.quoter.flow = models.RemittanceMgaFlow
	}
}

func (qls *QuoteLifeModSuite) AfterTest(suiteName, testName string) {}

func getMockSellableLife(age int) func() (*models.Product, error) {
	return func() (*models.Product, error) {
		// inject guarantees by age into product
		return nil, nil
	}
}

func baseLifeModQuoteTest(qls *QuoteLifeModSuite) {
	for index, input := range qls.inputs {
		// inject input in policyStub
		birthDate := time.Now().UTC().AddDate(-input.Age, 0, 0).Format(time.RFC3339)
		qls.quoter.policy.Contractor.BirthDate = birthDate
		// qls.quoter.rulesProduct.Companies[0].GuaranteesMap = input.GuaranteMaps
		qls.quoter.sellable = getMockSellableLife(input.Age)

		LifeMod(qls.quoter)

		qls.Assert().Equal(qls.outputs[index], qls.quoter.policy)
	}
}

func (qls *QuoteLifeModSuite) TestQuoteLifeEcommerceOffers() {
	baseLifeModQuoteTest(qls)
}

func (qls *QuoteLifeModSuite) TestQuoteLifeProviderOffers() {
	baseLifeModQuoteTest(qls)
}

func (qls *QuoteLifeModSuite) TestQuoteLifeRemittanceOffers() {
	baseLifeModQuoteTest(qls)
}

func TestLifeModQuote(t *testing.T) {
	suite.Run(t, new(QuoteLifeModSuite))
}
