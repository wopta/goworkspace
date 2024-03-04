package sellable_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
)

type SellableSuite struct {
	suite.Suite
	policy      *models.Policy
	channel     string
	networkNode *models.NetworkNode
	warrant     *models.Warrant
}

func (s *SellableSuite) SetupSuite() {
	os.Setenv("env", "local-test")
}

func (s *SellableSuite) SetupTest() {}

func (s *SellableSuite) TearDownSuite() {}

func (s *SellableSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestGetProductLifeEcommerce":
		s.policy = getPolicyByContractorAge(18)
		s.channel = models.ECommerceChannel
		s.networkNode = nil
		s.warrant = nil
	}
}

func (s *SellableSuite) AfterTest(suiteName, testName string) {}

func TestSellable(t *testing.T) {
	suite.Run(t, new(SellableSuite))
}

func (s *SellableSuite) TestGetProductLifeEcommerce() {
	got, err := sellable.GetProduct(s.policy, s.channel, s.networkNode, s.warrant)
	expected := &models.Product{}

	s.Assert().Equal(got, expected)
	s.Assert().NoError(err, "should not error")
}

// utils
func getPolicyByContractorAge(age int) *models.Policy {
	return &models.Policy{
		Name:           models.LifeProduct,
		ProductVersion: models.ProductV2,
		Contractor: models.Contractor{
			BirthDate: time.Now().UTC().AddDate(-age, 0, 0).Format(time.RFC3339),
		},
	}
}
