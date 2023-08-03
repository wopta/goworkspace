package product

import (
	"github.com/wopta/goworkspace/models"
)

func GetCommissionProducts(data models.Policy, products []models.Product) float64 {
	var commission float64
	for _, prod := range products {
		if prod.Name == data.Name {
			return GetCommissionProduct(data, prod)
		}

	}
	return commission
}

func GetCommissionProduct(data models.Policy, prod models.Product) float64 {
	var (
		amountNet, commissionValue float64
	)

	if data.PaymentSplit == string(models.PaySplitMonthly) {
		amountNet = data.PriceNettMonthly
	} else {
		amountNet = data.PriceNett
	}

	for _, company := range prod.Companies {
		if data.Company == company.Name {
			if company.CommissionSetting.IsFlat {
				return calculateCommission(amountNet, data.IsRenew, company.CommissionSetting.Commissions)
			}
			if company.CommissionSetting.IsByOffer {
				return calculateCommission(amountNet, data.IsRenew, prod.Offers[data.OfferlName].Commissions)
			}

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

func calculateCommission(amount float64, isRenew bool, commissions *models.Commission) float64 {
	if isRenew {
		return amount * commissions.Renew
	}
	return amount * commissions.NewBusiness
}
