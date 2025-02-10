package dto

import (
	"github.com/wopta/goworkspace/lib"
)

type priceDTO struct {
	Gross     float64
	GrossText string
	Net       float64
	NetText   string
	Taxes     float64
	TaxesText string
}

func newPriceDTO() *priceDTO {
	return &priceDTO{
		Gross:     0,
		GrossText: lib.HumanaizePriceEuro(0),
		Net:       0,
		NetText:   lib.HumanaizePriceEuro(0),
		Taxes:     0,
		TaxesText: lib.HumanaizePriceEuro(0),
	}
}
