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
	log.Println(r.RequestURI)
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	base := 1
	if strings.Contains(r.RequestURI, "enrich") {
		base = 2
	} else {
		base = 1
	}
	log.Println("base ", base)
	log.Println(vat[base])
	switch vat[base] {
	case "vat":
		switch vat[base+1] {
		case "munichre":
			MunichVat(w, vat[base+2])
		default:
			fmt.Fprintf(w, "missing service in path es. munichre")
		}

	case "works":
		Works(w)
	case "ateco":

	default:
		fmt.Fprintf(w, "missing scope in path es. vat, works, ateco")
	}

}
