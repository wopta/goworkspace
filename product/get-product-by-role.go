package product

import (
	"errors"
	"log"

	"github.com/wopta/goworkspace/models"
)

const (
	minAge         = "minAge"
	minReservedAge = "minReservedAge"
)

// DEPRECATED
func GetProductByRole(productName, version, company string, authToken models.AuthToken) (models.Product, error) {
	log.Println("GetProductByRole")
	var (
		responseProduct *models.Product
		err             error
	)

	switch authToken.Role {
	case models.UserRoleAdmin, models.UserRoleManager:
		responseProduct, err = getMgaProduct(productName, version, company)
	case models.UserRoleAll, models.UserRoleCustomer:
		responseProduct, err = getEcommerceProduct(productName, version, company)
	case models.UserRoleAgency, models.UserRoleAgent:
		responseProduct, err = getNetworkNodeProduct(authToken.Type, productName, version, company)
	default:
		responseProduct, err = productNotFound()
	}

	return *responseProduct, err
}

// DEPRECATED
func getMgaProduct(productName, version, company string) (*models.Product, error) {
	log.Println("getMgaProduct")
	return GetProduct(productName, version, models.MgaChannel)
}

// DEPRECATED
func getEcommerceProduct(productName, version, company string) (*models.Product, error) {
	log.Println("getEcommerceProduct")
	ecomProduct, err := GetProduct(productName, version, models.ECommerceChannel)

	if !ecomProduct.IsEcommerceActive {
		return productNotActive()
	}

	return ecomProduct, err
}

// DEPRECATED
func getNetworkNodeProduct(nodeType, productName, version, company string) (*models.Product, error) {
	log.Println("[getNetworkNodeProduct]")

	channel := models.AgentChannel

	if nodeType == models.AgencyNetworkNodeType {
		channel = models.AgencyChannel
	}

	return GetProduct(productName, version, channel)
}

// DEPRECATED
func getProductByName(products []models.Product, productName string) *models.Product {
	log.Println("getProductByName")
	mapProduct := map[string]models.Product{}
	for _, p := range products {
		mapProduct[p.Name] = p
	}
	if p, ok := mapProduct[productName]; ok {
		return &p
	}
	return nil
}

// DEPRECATED
func overrideProduct(baseProduct *models.Product, insertedProduct *models.Product) {
	log.Println("overrideProduct")
	if len(insertedProduct.Steps) > 0 {
		baseProduct.Steps = insertedProduct.Steps
	}

	for _, c := range insertedProduct.Companies {
		for _, c2 := range baseProduct.Companies {
			if c2.Name == c.Name {
				c2.Mandate = c.Mandate
			}
		}
	}
}

func productNotActive() (*models.Product, error) {
	return nil, errors.New("product not active")
}

func productNotFound() (*models.Product, error) {
	return nil, errors.New("product not found")
}
