package enrichVatCode

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	// err is pre-declared to avoid shadowing client.
}

func QuoteAllrisk(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("QuoteAllrisk")
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	if len(vat) > 1 {
		if strings.EqualFold(vat[1], "munichre") {
			var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/quote/rate/"
			u, err := url.Parse(urlstring)
			if err != nil {
				panic(err)
			}
			log.Println("url parse:", u)
			client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
				os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))
			jsonData, _ := ioutil.ReadAll(r.Body)
			req, _ := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(jsonData))
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
}
