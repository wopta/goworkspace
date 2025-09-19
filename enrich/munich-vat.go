package enrich

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

func munichVatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.AddPrefix("MunichVatFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	var body []byte

	vat := chi.URLParam(r, "vat")
	log.Println(vat)

	var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/company/vat/" + vat
	u, err := url.Parse(urlstring)
	lib.CheckError(err)
	log.Printf("url parse: %v", u)
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
