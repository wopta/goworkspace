package product

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
)

// DEPRECATED
func GetAllProductsByChannel(channel string) []models.Product {
	const (
		filePath = "products/mga/"
	)

	products := make([]models.Product, 0)

	rawProducts := lib.GetFolderContentByEnv(filePath)
	log.Printf("[GetAllProductsByChannel] found %d products for channel %s", len(rawProducts), models.MgaChannel)
	for _, rawProduct := range rawProducts {
		var product *models.Product
		var isActive bool
		err := json.Unmarshal(rawProduct, &product)
		lib.CheckError(err)
		switch channel {
		case models.MgaChannel:
			products = append(products, *product)
			continue
		case models.AgencyChannel:
			isActive = product.IsAgencyActive
		case models.AgentChannel:
			isActive = product.IsAgentActive
		case models.ECommerceChannel:
			isActive = product.IsEcommerceActive
		}

		log.Printf("[GetAllProductsByChannel] product %s version %s isActive %v", product.Name, product.Version, isActive)

		if isActive {
			product, err = GetProduct(product.Name, product.Version, channel)
			lib.CheckError(err)
			log.Printf("[GetAllProductsByChannel] found product %s version %s", product.Name, product.Version)
			products = append(products, *product)
		}
	}

	return products
}
