package quote

import (
	"log"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func removeOfferRate(policy *models.Policy, availableRates []string) {
	for offerKey, offerValue := range policy.OffersPrices {
		for rate := range offerValue {
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

func addConsultacyPrice(policy *models.Policy, product *models.Product) {
	if product.ConsultancyConfig == nil || !product.ConsultancyConfig.IsActive {
		policy.ConsultancyValue.Percentage = 0
		policy.ConsultancyValue.Price = 0
		return
	}

	if !product.ConsultancyConfig.IsConfigurable {
		policy.ConsultancyValue.Percentage = product.ConsultancyConfig.DefaultValue
	}
	
	policy.ConsultancyValue.Price = lib.RoundFloat(policy.PriceGross * policy.ConsultancyValue.Percentage, 2)
}
