package enrichVatCode

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Rules")

	functions.HTTP("EnrichVat", EnrichVat)
}

func EnrichVat(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the main request.
	//lib.EnableCors(&w, r)
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	log.Println("EnrichVat")
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	w.Header().Set("Content-Type", "application/json")
	if len(vat) > 2 {
		if strings.EqualFold(vat[2], "munichre") {
			var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/company/vat/" + vat[1]
			u, err := url.Parse(urlstring)
			if err != nil {
				panic(err)
			}
			log.Println("url parse:", u)
			client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
				os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))

			req, _ := http.NewRequest("GET", urlstring, nil)
			req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MUNICHRESUBSCRIPTIONKEY"))
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
			} else {
				fmt.Fprintf(w, "campo RequestURI banca dati non presente")
			}
		} else {
			fmt.Fprintf(w, "campi RequestURI partita iva e banca dati mancanti")
		}

	}
	log.Println("Header", w.Header())

}
