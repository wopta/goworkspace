package payment

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func Controller(policy models.Policy, product models.Product, transactions []models.Transaction, scheduleFirstRate bool, customerId string) (string, []models.Transaction, error) {
	var (
		err            error
		paymentMethods []string
	)

	log.Printf("init")

	if len(transactions) == 0 {
		log.Printf("%02d is an invalid number of transactions", len(transactions))
		return "", nil, errors.New("no valid transactions")
	}

	if err = checkPaymentModes(policy); err != nil {
		log.Printf("mismatched payment configuration: %s", err.Error())
		return "", nil, err
	}

	paymentMethods = getPaymentMethods(policy, product)

	switch policy.Payment {
	case models.FabrickPaymentProvider:
		return fabrickIntegration(transactions, paymentMethods, policy, scheduleFirstRate, customerId)
	case models.ManualPaymentProvider:
		return remittanceIntegration(transactions)
	default:
		return "", nil, fmt.Errorf("payment provider %s not supported", policy.Payment)
	}
}

func fabrickIntegration(transactions []models.Transaction, paymentMethods []string, policy models.Policy, scheduleFirstRate bool, customerId string) (payUrl string, updatedTransactions []models.Transaction, err error) {
	hasMandate := false
	now := time.Now().UTC()

	if hasMandate, err = fabrickHasMandate(customerId); err != nil {
		log.Printf("error checking mandate: %s", err.Error())
	}

	if customerId == "" {
		customerId = uuid.New().String()
	}

	for index, tr := range transactions {
		isFirstOfBatch := index == 0
		isFirstRateOfAnnuity := policy.StartDate.Month() == tr.EffectiveDate.Month()

		createMandate := shouldCreateMandate(policy, tr, isFirstOfBatch, isFirstRateOfAnnuity, hasMandate)

		tr.ProviderName = models.FabrickPaymentProvider

		scheduleDate, err := time.Parse(time.DateOnly, tr.ScheduleDate)
		if err != nil {
			log.Printf("error parsing scheduleDate: %s", err.Error())
			return "", nil, err
		}
		if scheduleFirstRate && scheduleDate.Before(now) {
			/*
				sets schedule date to today + 1 in order to avoid corner case in which fabrick is not able to
				execute transaction when recreated at the end of the day
			*/
			tr.ScheduleDate = now.AddDate(0, 0, 1).Format(time.DateOnly)
		}

		res := <-createFabrickTransaction(&policy, tr, createMandate, scheduleFirstRate, isFirstRateOfAnnuity, customerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}
		if isFirstOfBatch && (!hasMandate || createMandate) {
			payUrl = *res.Payload.PaymentPageURL
		}
		log.Printf("transaction %02d payUrl: %s", index+1, *res.Payload.PaymentPageURL)

		tr.ProviderId = *res.Payload.PaymentID
		tr.PayUrl = *res.Payload.PaymentPageURL
		tr.UserToken = customerId
		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	return payUrl, updatedTransactions, nil
}

func remittanceIntegration(transactions []models.Transaction) (payUrl string, updatedTransaction []models.Transaction, err error) {
	updatedTransaction = make([]models.Transaction, 0)

	for index, tr := range transactions {
		now := time.Now().UTC()
		if index == 0 && tr.Annuity == 0 {
			tr.IsPay = true
			tr.Status = models.TransactionStatusPay
			tr.StatusHistory = append(tr.StatusHistory, models.TransactionStatusPay)
			tr.PayDate = now
			tr.TransactionDate = now
			tr.PaymentMethod = models.PayMethodRemittance
		}
		tr.ProviderId = ""
		tr.UserToken = ""
		tr.UpdateDate = now
		updatedTransaction = append(updatedTransaction, tr)
	}
	return "", updatedTransaction, nil
}

func getPaymentMethods(policy models.Policy, product models.Product) []string {
	var paymentMethods = make([]string, 0)

	log.Printf("[GetPaymentMethods] loading available payment methods for %s payment provider", policy.Payment)

	for _, provider := range product.PaymentProviders {
		if provider.Name == policy.Payment {
			for _, config := range provider.Configs {
				if config.Mode == policy.PaymentMode && config.Rate == policy.PaymentSplit {
					paymentMethods = append(paymentMethods, config.Methods...)
				}
			}
		}
	}

	log.Printf("[GetPaymentMethods] found %v", paymentMethods)
	return paymentMethods
}

func checkPaymentModes(policy models.Policy) error {
	var allowedModes []string

	switch policy.PaymentSplit {
	case string(models.PaySplitMonthly):
		allowedModes = models.GetAllowedMonthlyModes()
	case string(models.PaySplitYearly):
		allowedModes = models.GetAllowedYearlyModes()
	case string(models.PaySplitSingleInstallment):
		allowedModes = models.GetAllowedSingleInstallmentModes()
	}

	if !lib.SliceContains(allowedModes, policy.PaymentMode) {
		return fmt.Errorf("mode '%s' is incompatible with split '%s'", policy.PaymentMode, policy.PaymentSplit)
	}

	return nil
}

func shouldCreateMandate(p models.Policy, tr models.Transaction, isFirstTransaction, isFirstRateOfAnnuity, hasMandate bool) bool {
	isFirstRateOfPolicy := p.StartDate.Truncate(time.Hour * 24).Equal(tr.EffectiveDate.Truncate(time.Hour * 24))

	return p.PaymentMode == models.PaymentModeRecurrent && isFirstTransaction &&
		(isFirstRateOfPolicy || (isFirstRateOfAnnuity && !hasMandate))
}
