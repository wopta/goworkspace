package companydata

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Println("companydata")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/pmi/global/emit",
				Handler: PmiGlobalEmit,
				Method:  "GET",
			},
			{
				Route:   "/v1/axa/emit",
				Handler: PmiGlobalEmit,
				Method:  http.MethodPost,
			},
			{
				Route:   "/v1/emit",
				Handler: Emit,
				Method:  http.MethodPost,
			},
		},
	}
	route.Router(w, r)

}
