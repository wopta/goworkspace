package enrich

import (
	"github.com/wopta/goworkspace/models"
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
	log.Println("Enrich")
	log.Println(r.RequestURI)
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/vat/munichre/:vat",
				Handler: MunichVat,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/ateco/:ateco",
				Handler: Ateco,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/cities",
				Handler: Cities,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/works",
				Handler: Works,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}
