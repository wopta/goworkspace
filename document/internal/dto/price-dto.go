package dto

import (
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type priceDTO struct {
	Gross       numeric
	Net         numeric
	Taxes       numeric
	Consultancy numeric
	Total       numeric
	Split       string
}

func newPriceDTO() *priceDTO {
	return &priceDTO{
		Gross:       newNumeric(),
		Net:         newNumeric(),
		Taxes:       newNumeric(),
		Consultancy: newNumeric(),
		Total:       newNumeric(),
	}
}

func (price *priceDTO) fromPolicy(policy models.Policy) {
	price.Split = getSplit(policy.PaymentSplit)
	price.Gross.ValueFloat = policy.PriceGross
	price.Gross.Text = lib.HumanaizePriceEuro(policy.PriceGross)
	price.Consultancy.ValueFloat = policy.ConsultancyValue.Price
	price.Consultancy.Text = lib.HumanaizePriceEuro(policy.ConsultancyValue.Price)
	price.Total.ValueFloat = policy.ConsultancyValue.Price + policy.PriceGross
	price.Total.Text = lib.HumanaizePriceEuro(policy.ConsultancyValue.Price + policy.PriceGross)
}
