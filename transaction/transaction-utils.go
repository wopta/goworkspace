package transaction

import (
	"errors"
	"time"

	"github.com/wopta/goworkspace/models"
)

const (
	transactionStatusReinitialized string = "Reinitialized"
	policyStatusReinitialized      string = "Reinitialized"
)

func ReinitializePaymentInfo(tr *models.Transaction, providerName string) error {
	if tr.IsPay && !tr.IsDelete {
		return errors.New("cannot reinitialize paid transaction")
	}
	tr.ProviderName = providerName
	tr.IsPay = false
	tr.IsDelete = false
	tr.PaymentNote = ""
	tr.PaymentMethod = ""
	tr.PayDate = time.Time{}
	tr.PayUrl = ""
	tr.TransactionDate = time.Time{}
	if !tr.EffectiveDate.IsZero() {
		tr.ScheduleDate = tr.EffectiveDate.Format(time.DateOnly)
		tr.ExpirationDate = tr.EffectiveDate.AddDate(10, 0, 0).Format(time.DateOnly)
	}
	tr.Status = models.TransactionStatusToPay
	tr.StatusHistory = append(tr.StatusHistory, transactionStatusReinitialized, models.TransactionStatusToPay)
	tr.UpdateDate = time.Now().UTC()
	return nil
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
