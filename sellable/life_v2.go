package sellable

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	prd "github.com/wopta/goworkspace/product"
)

type SellableProduct struct {
	authToken models.AuthToken
}

type SellableI interface {
	GetProduct(name, version string) *models.Product
	// GetRulesFile() []byte
	Exec(*models.Product, []byte) *models.Product
}

func (s *SellableProduct) GetProduct(name, version string) *models.Product {
	var warrant *models.Warrant
	channel := s.authToken.GetChannelByRoleV2()
	networkNode := network.GetNetworkNodeByUid(s.authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	return prd.GetProductV2(name, version, channel, networkNode, warrant)
}

// func (s *SellableProduct) GetRulesFile(name, version string) []byte {
// 	return lib.GetRulesFileV2(name, version, rulesFilename)
// }

func (s *SellableProduct) Exec(prd *models.Product, in []byte) *models.Product {
	rulesBytes := lib.GetRulesFileV2(prd.Name, prd.Version, rulesFilename)
	fx := new(models.Fx)
	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesBytes, prd, in, nil)

	return ruleOutput.(*models.Product)
}

func LifeMod(policy models.Policy, s SellableI) (product *models.Product, err error) {
	in, err := getInputData(policy)
	if err != nil {
		return nil, err
	}

	product = s.GetProduct(policy.Name, policy.ProductVersion)
	product = s.Exec(product, in)

	return
}
