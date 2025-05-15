package catnat

import (
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/quote/internal"
	"github.com/wopta/goworkspace/sellable"
)

type sellableCatnat func(policy *models.Policy, product *models.Product, isValidationForQuote bool) (*sellable.SellableOutput, error)

func CatnatQuote(policy *models.Policy, product *models.Product, sellable sellableCatnat, catnatClient INetClient) (resp QuoteResponse, err error) {
	outSellable, err := sellable(policy, product, true)
	if err != nil {
		return resp, err
	}
	internal.AddGuaranteesSettingsFromProduct(policy, outSellable.Product)

	var cnReq QuoteRequest
	err = cnReq.FromPolicy(policy, false)
	if err != nil {
		log.ErrorF("error building NetInsurance DTO: %s", err.Error())
		return resp, err
	}

	resp, err = catnatClient.Quote(cnReq)
	if err != nil {
		return resp, err
	}
	resp.ToPolicy(policy)
	internal.AddConsultacyPrice(policy, product)
	return
}
