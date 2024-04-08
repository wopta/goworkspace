package transaction

import (
	"github.com/wopta/goworkspace/models"
	"time"
)

func ReinitializePaymentInfo(tr *models.Transaction) {
	if tr.IsPay && !tr.IsDelete {
		return
	}
	tr.IsPay = false
	tr.IsDelete = false
	tr.PaymentNote = ""
	tr.PaymentMethod = ""
	tr.PayDate = time.Time{}
	tr.TransactionDate = time.Time{}
	tr.Status = models.TransactionStatusToPay
	tr.StatusHistory = append(tr.StatusHistory, "Reinitialized", models.TransactionStatusToPay)
}

func getMonthlyAmountsFlat(policy *models.Policy) (grossAmounts []float64, nettAmounts []float64) {
	numberOfRates := 12
	grossAmounts = make([]float64, numberOfRates)
	nettAmounts = make([]float64, numberOfRates)

	for rateIndex := 0; rateIndex < numberOfRates; rateIndex++ {
		grossAmounts[rateIndex] = policy.PriceGrossMonthly
		nettAmounts[rateIndex] = policy.PriceNettMonthly
	}

	return grossAmounts, nettAmounts
}

func getYearlyAmountsByGuarantee(policy *models.Policy) (grossAmounts []float64, nettAmounts []float64) {
	durationInYears := policy.GetDurationInYears()
	grossAmounts = make([]float64, durationInYears)
	nettAmounts = make([]float64, durationInYears)

	for _, guarantee := range policy.Assets[0].Guarantees {
		for rateIndex := 0; rateIndex < guarantee.Value.Duration.Year; rateIndex++ {
			grossAmounts[rateIndex] += guarantee.Value.PremiumGrossYearly
			nettAmounts[rateIndex] += guarantee.Value.PremiumNetYearly
		}
	}

	return grossAmounts, nettAmounts
}

func getYearlyAmountsByGuaranteeFlat(policy *models.Policy) (grossAmounts []float64, nettAmounts []float64) {
	durationInYears := policy.GetDurationInYears()
	grossAmounts = make([]float64, durationInYears)
	nettAmounts = make([]float64, durationInYears)

	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			for rateIndex := 0; rateIndex < durationInYears; rateIndex++ {
				grossAmounts[rateIndex] += guarantee.PriceGross
				nettAmounts[rateIndex] += guarantee.PriceNett
			}
		}
	}

	return grossAmounts, nettAmounts
}

func getYearlyAmountsFlat(policy *models.Policy) (grossAmounts []float64, nettAmounts []float64) {
	grossAmounts = append(grossAmounts, policy.PriceGross)
	nettAmounts = append(nettAmounts, policy.PriceNett)

	return grossAmounts, nettAmounts
}
