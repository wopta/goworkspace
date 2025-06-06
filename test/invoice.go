package test

import (
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/accounting"
)

func createInvoice(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("createInvoice")
	inv := accounting.InvoiceInc{

		Name:       "luca",
		VatNumber:  "brblcu81h03f205q",
		TaxCode:    "brblcu81h03f205q",
		Address:    "via test",
		PostalCode: "15057",
		City:       "Tortona",
		CityCode:   "AL",
		Country:    "Italia",
		Mail:       "luca.barbieri@wopta.it",
		Amount:     0,
		Date:       time.Now(),
		PayDate:    time.Now(),
		Items: []accounting.Items{{
			Desc:       "Contributo per intermediazione",
			Name:       "",
			Code:       "Vita",
			Qty:        1,
			ProductId:  0,
			GrossPrice: 0,
			Category:   "Vita",
			Date:       time.Now()}}}

	//"inv.CreateInvoice(false, true)
	accounting.DoInvoicePaid(inv, "test/test.pdf")

	return string(""), nil, nil
}
