package enrich

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var enrichRoutes []lib.Route = []lib.Route{
	{
		Route:       "/vat/munichre/{vat}",
		Fn:          MunichVatFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "enrich.munichre.vat",
	},
	{
		Route:       "/ateco/{ateco}",
		Fn:          AtecoFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "enrich.ateco",
	},

	{
		Route:       "/cities",
		Fn:          CitiesFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "enrich.cities",
	},

	{
		Route:       "/works",
		Fn:          WorksFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "enrich.works",
	},
	{
		Route:       "/naics",
		Fn:          NaicsFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "enrich.naics",
	},
}

func init() {
	log.Println("INIT Enrich")
	functions.HTTP("Enrich", Enrich)
}

func Enrich(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("enrich", enrichRoutes)
	router.ServeHTTP(w, r)
}
