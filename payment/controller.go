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

	// TODO: fix me
	if policy.Payment == "" || policy.Payment == "fabrik" {
		policy.Payment = models.FabrickPaymentProvider
	}
	paymentMethods = getPaymentMethods(*policy, product)

	log.Printf("generating payment URL")
	switch policy.Payment {
	case models.FabrickPaymentProvider:
		var payRes FabrickPaymentResponse

		switch policy.PaymentSplit {
		case string(models.PaySplitYear), string(models.PaySplitYearly), string(models.PaySplitSingleInstallment):
			log.Printf("fabrick yearly pay")
			payRes = FabrickYearPay(*policy, origin, paymentMethods, mgaProduct)
		case string(models.PaySplitMonthly):
			log.Printf("fabrick monthly pay")
			payRes = FabrickMonthlyPay(*policy, origin, paymentMethods, mgaProduct)
		}
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

	// TODO: remove me once established standard
	if policy.PaymentSplit == string(models.PaySplitYear) {
		policy.PaymentSplit = string(models.PaySplitYearly)
	}

	for _, provider := range product.PaymentProviders {
		if provider.Name == policy.Payment {
			for _, method := range provider.Methods {
				if lib.SliceContains(method.Rates, policy.PaymentSplit) {
					paymentMethods = append(paymentMethods, method.Name)
				}
			}
		}
	}

	log.Printf("[GetPaymentMethods] found %v", paymentMethods)
	return paymentMethods
}
