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

var sellableRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/sales/life",
		Fn:          LifeFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "sellable.sales.life",
	},
	{
		Route:       "/v1/risk/person",
		Fn:          PersonaFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "sellable.risk.person",
	},
	{
		Route:       "/v1/commercial-combined",
		Fn:          CommercialCombinedFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "sellable.commercialcombined",
	},
}

func init() {
	log.Println("INIT Sellable")
	functions.HTTP("Sellable", Sellable)
}

func Sellable(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("sellable", sellableRoutes)
	router.ServeHTTP(w, r)
}
