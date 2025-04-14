package test

import (
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/accounting"
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
			Desc:      "Contributo per intermediazione",
			Name:      "",
			Code:      "Vita",
			Qty:       1,
			ProductId: 0,
			NetPrice:  0,
			Category:  "Vita",
			Date:      time.Now()}}}

	//"inv.CreateInvoice(false, true)
	accounting.DoInvicePaid(inv, "test/test.pdf")

	return string(""), nil, nil
}
