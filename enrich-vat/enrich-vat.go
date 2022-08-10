package enrichVatCode

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	lib "github.com/wopta/goworkspace/lib"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	// err is pre-declared to avoid shadowing client.
}

type publishRequest struct {
	vat string `json:"vat"`
}

//  enrichVatCode
// .
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
func EnrichVat(w http.ResponseWriter, r *http.Request) {

	os.Getenv("munichreBaseUrl")
	os.Getenv("munichreSubscriptionKey")
	os.Getenv("munichreSubscriptionHeader")

	var url = "https://api-devexternal.munichre.com/flowin/dev/api/V1/api/company/vat/01654010345"

	client := lib.ClientCredentials(os.Getenv("munichreClientId"),
		os.Getenv("munichreClientSecret"), os.Getenv("munichreScope"), os.Getenv("munichreTokenEndPoint"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", "59c92bc0095d4b8c803656a207150c32")
	res, err := client.Do(req)
	if err != nil {
		log.Println("errore:")
		log.Println(err)
	}
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()
		fmt.Fprintf(w, string(body))
	}

}
