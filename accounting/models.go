package accounting

import (
	"time"
)

type InvoiceInc struct {
	Name       string
	Desc       string
	VatNumber  string
	TaxCode    string
	Address    string
	PostalCode string
	City       string
	CityCode   string
	Country    string
	Mail       string
	Amount     float32
	Date       time.Time
	PayDate    time.Time
	Items      []Items
}
type Items struct {
	Name       string
	Desc       string
	Code       string
	ProductId  int32
	Qty        int32
	GrossPrice float32
	Category   string
	Date       time.Time
}
