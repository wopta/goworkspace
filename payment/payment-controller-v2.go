package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"os"
	"time"
)

func PaymentControllerV2(policy models.Policy, product models.Product, transactions []models.Transaction) (string, []models.Transaction, error) {
	var (
		err                error
		payUrl             string
		paymentMethods     []string
		updatedTransaction []models.Transaction
	)

	log.Printf("init")

	if err = checkPaymentModes(policy); err != nil {
		log.Printf("mismatched payment configuration: %s", err.Error())
		return "", nil, err
	}

	paymentMethods = getPaymentMethodsV2(policy, product)

	switch policy.Payment {
	case models.FabrickPaymentProvider:
		payUrl, updatedTransaction, err = fabrickIntegration(transactions, paymentMethods, policy)
	case models.ManualPaymentProvider:
		payUrl, updatedTransaction, err = remittanceIntegration(transactions)
	default:
		return "", nil, fmt.Errorf("payment provider %s not supported", policy.Payment)
	}

	return payUrl, updatedTransaction, nil
}

func remittanceIntegration(transactions []models.Transaction) (payUrl string, updatedTransaction []models.Transaction, err error) {
	for index, _ := range transactions {
		transactions[index].PaymentMethod = models.PayMethodRemittance
	}
	return "", transactions, nil
}

func fabrickIntegration(transactions []models.Transaction, paymentMethods []string, policy models.Policy) (payUrl string, updatedTransactions []models.Transaction, err error) {
	customerId := uuid.New().String()
	now := time.Now().UTC()

	for index, tr := range transactions {
		isFirstRate := index == 0
		createMandate := (policy.PaymentMode == models.PaymentModeRecurrent) && isFirstRate
		if isFirstRate {
			tr.ScheduleDate = ""
		}

		res := <-createFabrickTransactionV2(&policy, tr, isFirstRate, createMandate, customerId, paymentMethods)
		if res.Payload == nil || res.Payload.PaymentPageURL == nil {
			return "", nil, errors.New("error creating transaction on Fabrick")
		}
		if isFirstRate {
			payUrl = *res.Payload.PaymentPageURL
		}
		tr.ProviderId = *res.Payload.PaymentID
		tr.UserToken = customerId

		/*
			operation that has to be done if transaction has been already paid and canceled.
			Is it correct to do them here?
		*/
		tr.ProviderName = models.FabrickPaymentProvider
		tr.IsPay = false
		tr.IsDelete = false
		tr.PaymentNote = ""
		tr.PaymentMethod = ""
		tr.PayDate = time.Time{}
		tr.TransactionDate = time.Time{}
		tr.Status = models.TransactionStatusToPay
		tr.StatusHistory = append(tr.StatusHistory, models.TransactionStatusToPay)

		tr.UpdateDate = now
		updatedTransactions = append(updatedTransactions, tr)
	}

	return payUrl, updatedTransactions, nil
}

func createFabrickTransactionV2(
	policy *models.Policy,
	transaction models.Transaction,
	firstSchedule, createMandate bool,
	customerId string,
	paymentMethods []string,
) <-chan FabrickPaymentResponse {
	r := make(chan FabrickPaymentResponse)

	go func() {
		defer close(r)

		body := getFabrickRequestBody(policy, firstSchedule, transaction.ScheduleDate, transaction.ExpirationDate,
			customerId, transaction.Amount, "", paymentMethods)
		if body == "" {
			return
		}
		request := getFabrickPaymentRequest(body)
		if request == nil {
			return
		}

		log.Printf("policy '%s' request headers: %s", policy.Uid, request.Header)
		log.Printf("policy '%s' request body: %s", policy.Uid, request.Body)

		if os.Getenv("env") == "local" {
			status := "200"
			local := "local"
			url := "www.dev.wopta.it"
			r <- FabrickPaymentResponse{
				Status: &status,
				Errors: nil,
				Payload: &Payload{
					ExternalID:        &local,
					PaymentID:         &local,
					MerchantID:        &local,
					PaymentPageURL:    &url,
					PaymentPageURLB2B: &url,
					TokenB2B:          &local,
					Coupon:            &local,
				},
			}
		} else {

			res, err := lib.RetryDo(request, 5, 10)
			lib.CheckError(err)

			if res != nil {
				log.Printf("policy '%s' response headers: %s", policy.Uid, res.Header)
				body, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				lib.CheckError(err)
				log.Printf("policy '%s' response body: %s", policy.Uid, string(body))

				var result FabrickPaymentResponse

				if res.StatusCode != 200 {
					log.Printf("exiting with statusCode: %d", res.StatusCode)
					result.Errors = append(result.Errors, res.Status, res.StatusCode)
				} else {
					err = json.Unmarshal([]byte(body), &result)
					lib.CheckError(err)
				}

				r <- result
			}
		}
	}()

	return r
}

func getPaymentMethodsV2(policy models.Policy, product models.Product) []string {
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
