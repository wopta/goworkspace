package main

// GOOGLE_CLOUD_PROJECT is a user-set environment variable.
/*{
	"BaseUrl":"https://api-devexternal.munichre.com/flowin/dev/api/V1",
	"ApimKey":"59c92bc0095d4b8c803656a207150c32",
	"TokenEndPoint":"https://login.microsoftonline.com/9f2c9c2d-da50-4f33-8dfb-a780f38b50dd/oauth2/v2.0/token",
	"Scope":"46e8daaf-f894-464a-942a-e06852ed4526/.default",
	"ClientId":"194d46f8-0779-4e17-a96d-62c7bdd81901",
	"ClientSecret":"nrDLtgtLiVhvaChj1sU7JiCBUZbztRXw2ROMBYxZ",
	"GrantType":"client_credentials",
	"UWRole":"Agent",
	"SubProductId_PMIW":"35"
}*/
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	enr "github.com/wopta/goworkspace/enrich-vat"
	q "github.com/wopta/goworkspace/quote-allrisk"
	rules "github.com/wopta/goworkspace/rules"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	os.Setenv("munichreBaseUrl", "https://api-devexternal.munichre.com/flowin/dev/api/V1")
	os.Setenv("munichreTokenEndPoint", "https://login.microsoftonline.com/9f2c9c2d-da50-4f33-8dfb-a780f38b50dd/oauth2/v2.0/token")
	os.Setenv("munichreScope", "46e8daaf-f894-464a-942a-e06852ed4526/.default")
	os.Setenv("munichreClientId", "194d46f8-0779-4e17-a96d-62c7bdd81901")
	os.Setenv("munichreClientSecret", "nrDLtgtLiVhvaChj1sU7JiCBUZbztRXw2ROMBYxZ")
	os.Setenv("munichreSubscriptionKey", "59c92bc0095d4b8c803656a207150c32")
	os.Setenv("munichreSubscriptionHeader", "Ocp-Apim-Subscription-Key")

	enrich_vat := r.PathPrefix("/enrich-vat").Subrouter()
	rules_sub := r.PathPrefix("/rules").Subrouter()
	quote := r.PathPrefix("/quote").Subrouter()

	enrich_vat.HandleFunc("/{key}/{key}", enr.EnrichVat).Methods("GET")
	rules_sub.HandleFunc("/{key}", rules.Rules).Methods("POST")
	quote.HandleFunc("/quote", q.QuoteAllrisk).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
