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

func calculateCommissionV2(commissions *models.Commissions, isRenew, isActive bool, amount float64) float64 {
	log.Printf("[calculateCommissionV2] amount: %f", amount)
	var commission float64

	if isRenew {
		if isActive {
			log.Printf("[calculateCommissionV2] commission renew active at %f", commissions.Renew)
			commission = amount * commissions.Renew
		} else {
			log.Printf("[calculateCommissionV2] commission renew passive at %f", commissions.RenewPassive)
			commission = amount * commissions.RenewPassive
		}
	} else {
		if isActive {
			log.Printf("[calculateCommissionV2] commission new business active at %f", commissions.NewBusiness)
			commission = amount * commissions.NewBusiness
		} else {
			log.Printf("[calculateCommissionV2] commission new business passive at %f", commissions.NewBusinessPassive)
			commission = amount * commissions.NewBusinessPassive
		}
	}

	commission = lib.RoundFloat(commission, 2)
	log.Printf("[calculateCommissionV2] calculated commission: %f", commission)

	return commission
}

func GetCommissionByNode(policy *models.Policy, prod *models.Product, isActive bool) float64 {
	log.Println("[GetCommissionByNode]")

	var (
		amountNet, commissionValue float64
	)

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly), string(models.PaySplitSemestral):
		amountNet = policy.PriceNettMonthly
		log.Printf("[GetCommissionByNode] using PriceNettMonthly as amountNet: %g", amountNet)
	default:
		amountNet = policy.PriceNett
		log.Printf("[GetCommissionByNode] using PriceNett as amountNet: %g", amountNet)
	}

	for _, company := range prod.Companies {
		if policy.Company == company.Name {
			if company.CommissionSetting.IsFlat {
				log.Println("[GetCommissionByNode] Flat commission")
				return calculateCommissionV2(company.CommissionSetting.Commissions, policy.IsRenew, isActive, amountNet)
			}

			if company.CommissionSetting.IsByOffer {
				log.Println("[GetCommissionByNode] By offer commission")
				return calculateCommissionV2(prod.Offers[policy.OfferlName].Commissions, policy.IsRenew, isActive, amountNet)
			}

			log.Println("[GetCommissionByNode] By guarantee commission")
			for _, asset := range policy.Assets {
				for _, guarantee := range asset.Guarantees {
					if policy.PaymentSplit == string(models.PaySplitMonthly) {
						amountNet = guarantee.Value.PremiumNetMonthly
					} else {
						amountNet = guarantee.Value.PremiumNetYearly
					}
					commissionValue += calculateCommissionV2(guarantee.Commissions, policy.IsRenew, isActive, amountNet)
				}
			}
		}
	}

	return commissionValue
}
