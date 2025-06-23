package quote

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func RemoveOfferRate(policy *models.Policy, availableRates []string) {
	log.AddPrefix("RemoveOfferRate")
	defer log.PopPrefix()
	for offerKey, offerValue := range policy.OffersPrices {
		for rate := range offerValue {
			if !lib.SliceContains(availableRates, rate) {
				log.Printf("removing rate %s", rate)
				delete(policy.OffersPrices[offerKey], rate)
			}
		}
	}
}

func GetAvailableRates(product *models.Product, flow string) []string {
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

func AddConsultacyPrice(policy *models.Policy, product *models.Product) {
	if product.ConsultancyConfig == nil || !product.ConsultancyConfig.IsActive {
		policy.ConsultancyValue.Percentage = 0
		policy.ConsultancyValue.Price = 0
		return
	}

	if !product.ConsultancyConfig.IsConfigurable {
		policy.ConsultancyValue.Percentage = product.ConsultancyConfig.DefaultValue

		policy.ConsultancyValue.Price = lib.RoundFloat(policy.PriceGross*policy.ConsultancyValue.Percentage, 2)
	}
}

func AddGuaranteesSettingsFromProduct(policy *models.Policy, product *models.Product) {
	if len(policy.Assets) == 0 {
		return
	}
	for i, guaranteeReq := range policy.Assets[0].Guarantees {
		if guarantee, ok := product.Companies[0].GuaranteesMap[guaranteeReq.Slug]; ok && guarantee != nil {
			policy.Assets[0].Guarantees[i].IsSelected = product.Companies[0].GuaranteesMap[guaranteeReq.Slug].IsSelected
			policy.Assets[0].Guarantees[i].IsMandatory = product.Companies[0].GuaranteesMap[guaranteeReq.Slug].IsMandatory
			policy.Assets[0].Guarantees[i].IsSellable = product.Companies[0].GuaranteesMap[guaranteeReq.Slug].IsSellable
			policy.Assets[0].Guarantees[i].IsConfigurable = product.Companies[0].GuaranteesMap[guaranteeReq.Slug].IsConfigurable
		}
	}
}
