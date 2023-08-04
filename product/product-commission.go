package product

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetCommissionProducts(data models.Policy, products []models.Product) float64 {
	log.Println("[GetCommissionProducts]")
	var commission float64
	for _, prod := range products {
		if prod.Name == data.Name {
			log.Printf("[GetCommissionProducts] found product: %s", prod.Name)
			return GetCommissionProduct(data, prod)
		}

	}
	log.Println("[GetCommissionProducts] no product found")
	return commission
}

func GetCommissionProduct(data models.Policy, prod models.Product) float64 {
	log.Println("[GetCommissionProduct]")
	var (
		amountNet, commissionValue float64
	)

	switch data.PaymentSplit {
	case string(models.PaySplitMonthly), string(models.PaySplitSemestral):
		amountNet = data.PriceNettMonthly
		log.Printf("[GetCommissionProduct] using PriceNettMonthly as amountNet: %g", amountNet)
	default:
		amountNet = data.PriceNett
		log.Printf("[GetCommissionProduct] using PriceNett as amountNet: %g", amountNet)
	}

	for _, company := range prod.Companies {
		if data.Company == company.Name {
			if company.CommissionSetting.IsFlat {
				log.Println("[GetCommissionProduct] Flat commission")
				return calculateCommission(amountNet, data.IsRenew, company.CommissionSetting.Commissions)
			}
			if company.CommissionSetting.IsByOffer {
				log.Println("[GetCommissionProduct] By offer commission")
				return calculateCommission(amountNet, data.IsRenew, prod.Offers[data.OfferlName].Commissions)
			}

			log.Println("[GetCommissionProduct] By guarantee commission")
			for _, asset := range data.Assets {
				for _, guarantee := range asset.Guarantees {
					if data.PaymentSplit == string(models.PaySplitMonthly) {
						amountNet = guarantee.Value.PremiumNetMonthly
					} else {
						amountNet = guarantee.Value.PremiumNetYearly
					}
					commissionValue += calculateCommission(amountNet, data.IsRenew, guarantee.Commissions)
				}
			}
		}
	}
	return commissionValue
}

func calculateCommission(amount float64, isRenew bool, commissions *models.Commissions) float64 {
	var commission float64

	if isRenew {
		log.Println("[calculateCommission] commission renew")
		commission = amount * commissions.Renew
	} else {
		log.Println("[calculateCommission] commission new business")
		commission = amount * commissions.NewBusiness
	}

	return lib.RoundFloat(commission, 2)
}
