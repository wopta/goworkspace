package enrich

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	lib "github.com/wopta/goworkspace/lib"
)

func MunichVat(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	// Set CORS headers for the main request.
	log.Println("Munich Enrich Vat")
	log.Println(r.Header.Get("vat"))
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	var body []byte
	w.Header().Set("Content-Type", "application/json")
	vat := r.Header.Get("vat")
	var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/company/vat/" + vat
	u, err := url.Parse(urlstring)
	lib.CheckError(err)
	log.Println("url parse:", u)
	client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
		os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))
	req, _ := http.NewRequest("GET", urlstring, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MUNICHRESUBSCRIPTIONKEY"))
	res, err := client.Do(req)
	lib.CheckError(err)
	if res != nil {
		body, err = ioutil.ReadAll(res.Body)
		lib.CheckError(err)
		res.Body.Close()

	}
	log.Println("Header", w.Header())
	return string(body), nil, err
}
