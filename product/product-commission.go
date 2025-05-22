package product

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func calculateCommissionV2(commissions *models.Commissions, isRenew, isActive bool, amount float64) float64 {
	log.AddPrefix("CalculateCommissionV2")
	defer log.PopPrefix()
	log.Printf("amount: %f", amount)
	var commission float64

	if isRenew {
		if isActive {
			log.Printf("commission renew active at %f", commissions.Renew)
			commission = amount * commissions.Renew
		} else {
			log.Printf("commission renew passive at %f", commissions.RenewPassive)
			commission = amount * commissions.RenewPassive
		}
	} else {
		if isActive {
			log.Printf("commission new business active at %f", commissions.NewBusiness)
			commission = amount * commissions.NewBusiness
		} else {
			log.Printf("commission new business passive at %f", commissions.NewBusinessPassive)
			commission = amount * commissions.NewBusinessPassive
		}
	}

	commission = lib.RoundFloat(commission, 2)
	log.Printf("calculated commission: %f", commission)

	return commission
}

func GetCommissionByProduct(policy *models.Policy, prod *models.Product, isActive bool) float64 {
	log.AddPrefix("GetCommissionByProduct")
	defer log.PopPrefix()

	var (
		amountNet, commissionValue float64
	)

	amountNet = getCommissionAmountByPaymentSplit(policy)

	for _, company := range prod.Companies {
		if policy.Company == company.Name {
			if company.CommissionSetting.IsFlat {
				log.Println("Flat commission")
				return calculateCommissionV2(company.CommissionSetting.Commissions, policy.IsRenew, isActive, amountNet)
			}

			if company.CommissionSetting.IsByOffer {
				log.Println("By offer commission")
				return calculateCommissionV2(prod.Offers[policy.OfferlName].Commissions, policy.IsRenew, isActive, amountNet)
			}

			log.Println("By guarantee commission")
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
	log.AddPrefix("GetCommissionAmountByPaymentSplit")
	defer log.PopPrefix()

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly), string(models.PaySplitSemestral):
		log.Printf("using PriceNettMonthly as amountNet: %g", policy.PriceNettMonthly)
		return policy.PriceNettMonthly
	default:
		log.Printf("using PriceNett as amountNet: %g", policy.PriceNett)
		return policy.PriceNett
	}
}
