package sellable

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	rulesFilename = "sellable"
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
				Handler: LifeFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/risk/person",
				Handler: PersonaFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}
