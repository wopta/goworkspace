package form

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT Form")
	functions.HTTP("Form", Form)
}

func Form(w http.ResponseWriter, r *http.Request) {
	log.Println("Product")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/axafleet",
				Handler: AxaFleetTway,
				Method:  "GET",
			},
			{
				Route:   "/v1/:uid",
				Handler: GetFx,
				Method:  "GET",
			},
		},
	}
	route.Router(w, r)

}
