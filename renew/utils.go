package renew

import (
	"fmt"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/product"
	"strings"
)

func getProductsMapByPolicyType(policyType string) map[string]models.Product {
	products := make(map[string]models.Product)

	productsInfo := product.GetAllProductsByChannel(models.MgaChannel)
	for _, pr := range productsInfo {
		prd := product.GetProductV2(pr.Name, pr.Version, models.MgaChannel, nil, nil)
		if strings.EqualFold(prd.PolicyType, policyType) {
			key := fmt.Sprintf("%s-%s", prd.Name, prd.Version)
			products[key] = *prd
		}
	}

	return products
}
