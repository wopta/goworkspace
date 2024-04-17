package enrich

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
)

func MunichVatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[MunichVatFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	var body []byte

	vat := chi.URLParam(r, "vat")
	log.Println(vat)

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
		body = lib.ErrorByte(io.ReadAll(res.Body))
		defer res.Body.Close()
	}

	log.Println("Handler end -------------------------------------------------")

	return string(body), nil, err
}
