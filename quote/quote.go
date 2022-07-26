package quote

import (
	"bytes"
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
	functions.HTTP("Quote", Quote)
}

func Quote(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	log.Println("QuoteAllrisk")
	base := "/quote"
	if strings.Contains(r.RequestURI, "/quote") {
		base = "/quote"
	} else {
		base = ""
	}
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)

	switch os := r.RequestURI; os {
	case base + "/pmi/munichre":
		PmiMunich(w, r)
	case base + "/incident":

	default:
		fmt.Fprintf(w, "select right field")
	}

}
func PmiMunich(w http.ResponseWriter, r *http.Request) {
	var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/quote/rate/"
	u, err := url.Parse(urlstring)
	lib.CheckError(err)
	log.Println("url parse:", u)
	client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
		os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))
	jsonData, _ := ioutil.ReadAll(r.Body)
	req, _ := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(jsonData))
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MUNICHRESUBSCRIPTIONKEY"))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	lib.CheckError(err)
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		fmt.Fprintf(w, string(body))
	}

}
