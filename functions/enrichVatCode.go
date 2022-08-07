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

func EnrichVatCode(w http.ResponseWriter, r *http.Request) {

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
