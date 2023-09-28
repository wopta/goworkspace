package accounting

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Accounting")
	functions.HTTP("Accounting", Accounting)
}

func Accounting(w http.ResponseWriter, r *http.Request) {
	log.Println("Accounting")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{},
	}

	route.Router(w, r)
}
