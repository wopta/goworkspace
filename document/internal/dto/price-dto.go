package dto

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
