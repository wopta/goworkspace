package sellable

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"log"
	"net/http"
)

func init() {
	log.Println("INIT Sellable")

	functions.HTTP("Sellable", Sellable)
}

func Sellable(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/sales/life",
				Handler: LifeHandler,
				Method:  http.MethodPost,
			},
			{
				Route:   "/v1/risk/person",
				Handler: PersonHandler,
				Method:  http.MethodPost,
			},
		},
	}
	route.Router(w, r)

}
