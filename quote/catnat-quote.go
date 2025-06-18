package quote

import (
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/models/catnat"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

type sellableCatnat = func(policy *models.Policy, product *models.Product, isValidationForQuote bool) (*sellable.SellableOutput, error)

type clientQuote = func(dto catnat.QuoteRequest, policy *models.Policy) (response catnat.QuoteResponse, err error)

func CatnatQuote(policy *models.Policy, product *models.Product, sellable sellableCatnat, clientQuote clientQuote) (resp catnat.QuoteResponse, err error) {
	outSellable, err := sellable(policy, product, true)
	if err != nil {
		return resp, err
	}
	internal.AddGuaranteesSettingsFromProduct(policy, outSellable.Product)

	var cnReq catnat.QuoteRequest
	err = cnReq.FromPolicyForQuote(policy)
	if err != nil {
		log.ErrorF("error building NetInsurance DTO: %s", err.Error())
		return resp, err
	}

	resp, err = clientQuote(cnReq, policy)
	log.PrintStruct("response quote", resp)
	if err != nil {
		return resp, err
	}
	internal.AddConsultacyPrice(policy, product)
	return
}
