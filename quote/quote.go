package quote

import (
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Rules")
	functions.HTTP("Quote", Quote)
}

func Quote(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	vat := strings.Split(r.RequestURI, "/")
	log.Println(vat)
	log.Println(len(vat))
	log.Println("QuoteAllrisk")

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/pmi/munichre",
				Handler: PmiMunichFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/incident",
				Handler: PmiMunichFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/life",
				Handler: LifeFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/person",
				Handler: PersonFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}
