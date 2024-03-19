package quote

import (
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
)

func removeOfferRate(policy *models.Policy, availableRates []string) {
	for offerKey, offerValue := range policy.OffersPrices {
		for rate, _ := range offerValue {
			if !lib.SliceContains(availableRates, rate) {
				log.Printf("[removeOfferRate] removing rate %s", rate)
				delete(policy.OffersPrices[offerKey], rate)
			}
		}
	}
}

func getAvailableRates(product *models.Product, flow string) []string {
	availableRates := make([]string, 0)
	for _, paymentProvider := range product.PaymentProviders {
		if lib.SliceContains(paymentProvider.Flows, flow) {
			for _, config := range paymentProvider.Configs {
				if !lib.SliceContains(availableRates, config.Rate) {
					availableRates = append(availableRates, config.Rate)
				}
			}
		}
	}
	return availableRates
}
