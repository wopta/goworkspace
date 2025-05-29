package catnat

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

type sellableCatnat func(policy *models.Policy, product *models.Product, isValidationForQuote bool) (*sellable.SellableOutput, error)

func CatnatQuote(policy *models.Policy, product *models.Product, sellable sellableCatnat, catnatClient INetClient) (resp QuoteResponse, err error) {
	outSellable, err := sellable(policy, product, true)
	if err != nil {
		return resp, err
	}
	internal.AddGuaranteesSettingsFromProduct(policy, outSellable.Product)

	var cnReq QuoteRequest
	err = cnReq.FromPolicyForQuote(policy)
	if err != nil {
		log.ErrorF("error building NetInsurance DTO: %s", err.Error())
		return resp, err
	}

	resp, err = catnatClient.Quote(cnReq)
	log.PrintStruct("response quote", resp)
	if err != nil {
		return resp, err
	}
	err = MappingQuoteResponseToPolicy(resp, policy)
	if err != nil {
		return resp, err
	}
	err = MappingQuoteResponseToGuarantee(resp, policy)
	internal.AddConsultacyPrice(policy, product)
	return
}
