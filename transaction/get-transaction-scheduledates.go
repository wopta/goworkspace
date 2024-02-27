package transaction

import (
	"log"
	"time"

	"github.com/wopta/goworkspace/models"
)

func GetTransactionScheduleDates(policy *models.Policy) []time.Time {
	var (
		currentScheduleDate time.Time
		response            []time.Time = make([]time.Time, 0)
		yearDuration        int         = 1
	)

	activeTransactions := GetPolicyActiveTransactions("", policy.Uid)

	if len(activeTransactions) == 0 {
		if policy.PaymentMode == models.PaymentModeRecurrent && policy.PaymentSplit == string(models.PaySplitYearly) {
			yearDuration = policy.GetDurationInYears()
		}

		numberOfRates := policy.GetNumberOfRates() * yearDuration

		for i := 0; i < numberOfRates; i++ {
			if i > 0 {
				switch policy.PaymentSplit {
				case string(models.PaySplitYearly):
					currentScheduleDate = policy.StartDate.AddDate(i, 0, 0)
				case string(models.PaySplitMonthly):
					currentScheduleDate = policy.StartDate.AddDate(0, i, 0)
				default:
					log.Printf("unhandled recurrent payment split: %s", policy.PaymentSplit)
					return nil
				}
			}
			response = append(response, currentScheduleDate)
		}
	} else {
		// TODO: handle when policy already has created transactions
		// isFirstSchedule := true
		// for _, tr := range activeTransactions {
		// 	if tr.IsPay {
		// 		continue
		// 	}
		// 	if isFirstSchedule {
		// 		currentScheduleDate = time.Time{}
		// 		isFirstSchedule = false
		// 	}
		// 	currentScheduleDate, err := time.Parse(models.TimeDateOnly, tr.ScheduleDate)
		// 	if err != nil {
		// 		log.Printf("error parsing schedule date %s: %s", tr.ScheduleDate, err.Error())
		// 		return nil
		// 	}
		// 	response = append(response, currentScheduleDate)
		// }
	}

	return response
}

func GetTransactionsAmounts(policy *models.Policy) (grossAmounts []float64, nettAmounts []float64) {
	switch policy.Name {
	case models.LifeProduct, models.PersonaProduct:
		if policy.PaymentMode == models.PaymentModeSingle {
			return getYearlyAmountsFlat(policy)
		}
		if policy.PaymentSplit == string(models.PaySplitYearly) {
			return getYearlyAmountsByGuarantee(policy)
		}
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			return getMonthlyAmountsFlat(policy)
		}
		log.Printf("not implemented - %s - %s", policy.Name, policy.PaymentSplit)
	case models.PmiProduct:
		if policy.PaymentMode == models.PaymentModeSingle {
			return getYearlyAmountsFlat(policy)
		}
		if policy.PaymentSplit == string(models.PaySplitYearly) {
			return getYearlyAmountsByGuaranteeFlat(policy)
		}
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			return getMonthlyAmountsFlat(policy)
		}
		log.Printf("not implemented - pmi - %s", policy.PaymentSplit)
	case models.GapProduct:
		if policy.PaymentSplit == string(models.PaySplitSingleInstallment) {
			return getYearlyAmountsFlat(policy)
		}
		log.Printf("not implemented - gap - %s", policy.PaymentSplit)
	}

	return grossAmounts, nettAmounts
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
