package enrichVatCode

import (
	"fmt"

	"log"
	"net/http"

	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Enrich")
	functions.HTTP("Enrich", Enrich)
}

func Enrich(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the main request.
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Println("EnrichVat")
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	base := 1
	if strings.Contains(vat[0], "enrich") {
		base = 1
	} else {
		base = 0
	}

	log.Println(vat[base])
	switch vat[base] {
	case "vat":
		switch vat[base+1] {
		case "munichre":
			MunichVat(w, vat[base+2])
		default:
			fmt.Fprintf(w, "")
		}

		MunichVat(w, vat[base+1])
	case "/works":

	case "/ateco":

	default:
		fmt.Fprintf(w, "")
	}

}
