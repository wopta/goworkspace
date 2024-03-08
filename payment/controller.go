package payment

import (
	"fmt"
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func PaymentController(origin string, policy *models.Policy, product, mgaProduct *models.Product) (string, error) {
	var (
		payUrl         string
		paymentMethods []string
	)

	log.Printf("init")

	if err := checkPaymentConfiguration(policy); err != nil {
		log.Printf("mismatched payment configuration: %s", err.Error())
		return "", err
	}

	paymentMethods = getPaymentMethods(*policy, product)

	switch policy.Payment {
	case models.FabrickPaymentProvider:
		payRes := fabrickPayment(policy, origin, paymentMethods, mgaProduct)

		if payRes.Payload == nil || payRes.Payload.PaymentPageURL == nil {
			log.Println("fabrick error payload or paymentUrl empty")
			return "", fmt.Errorf("fabrick error: %v", payRes.Errors)
		}
		payUrl = *payRes.Payload.PaymentPageURL
	default:
		return "", fmt.Errorf("payment provider %s not supported", policy.Payment)
	}

	log.Printf("payUrl: %s", payUrl)

	return payUrl, nil
}

func getPaymentMethods(policy models.Policy, product *models.Product) []string {
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

func checkPaymentConfiguration(policy *models.Policy) error {
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
