package product

import (
	"github.com/wopta/goworkspace/models"
)

func GetCommissionProduct(data models.Policy, prod models.Product) float64 {
	var commission float64
	for _, x := range prod.Companies {
		if x.Name == data.Company {
			if data.IsRenew {
				//TODO when pmi migration in done delete shit code check
				if x.Commission == 0 {
					return x.Mandate.CommissionRenew
				} else {
					return x.CommissionRenew
				}

			} else {
				//TODO when pmi migration in done delete shit code check
				if x.Commission == 0 {
					return x.Mandate.Commission
				} else {
					return x.Commission
				}

			}
		}

	}
	return commission
}
func GetCommissionProducts(data models.Policy, products []models.Product) float64 {
	var commission float64
	for _, prod := range products {
		if prod.Name == data.Name {
			return GetCommissionProduct(data, prod)
		}

	}
	return commission
}
