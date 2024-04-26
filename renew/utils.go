package renew

import (
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"strings"
)

func getProductsByPolicyType(policyType string) []models.Product {
	products := make([]models.Product, 0)

	productsInfo := product.GetAllProductsByChannel(models.MgaChannel)
	for _, pr := range productsInfo {
		prd := product.GetProductV2(pr.Name, pr.Version, models.MgaChannel, nil, nil)
		if strings.EqualFold(prd.PolicyType, policyType) {
			products = append(products, *prd)
		}
	}

	return products
}
