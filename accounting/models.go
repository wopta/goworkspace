package accounting

import (
	"time"
)

type InvoiceInc struct {
	Name       string
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
	Name      string
	Code      string
	ProductId int32
	NetPrice  float32
	Category  string
	Date      time.Time
}
