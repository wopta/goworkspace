package product

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

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

func GetCommissionByProduct(policy *models.Policy, prod *models.Product, isActive bool) float64 {
	log.Println("[GetCommissionByProduct]")

	var (
		amountNet, commissionValue float64
	)

	amountNet = getCommissionAmountByPaymentSplit(policy)

	for _, company := range prod.Companies {
		if policy.Company == company.Name {
			if company.CommissionSetting.IsFlat {
				log.Println("[GetCommissionByProduct] Flat commission")
				return calculateCommissionV2(company.CommissionSetting.Commissions, policy.IsRenew, isActive, amountNet)
			}

			if company.CommissionSetting.IsByOffer {
				log.Println("[GetCommissionByProduct] By offer commission")
				return calculateCommissionV2(prod.Offers[policy.OfferlName].Commissions, policy.IsRenew, isActive, amountNet)
			}

			log.Println("[GetCommissionByProduct] By guarantee commission")
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

func getCommissionAmountByPaymentSplit(policy *models.Policy) float64 {
	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly), string(models.PaySplitSemestral):
		log.Printf("[getCommissionAmountByPaymentSplit] using PriceNettMonthly as amountNet: %g", policy.PriceNettMonthly)
		return policy.PriceNettMonthly
	default:
		log.Printf("[getCommissionAmountByPaymentSplit] using PriceNett as amountNet: %g", policy.PriceNett)
		return policy.PriceNett
	}
}
