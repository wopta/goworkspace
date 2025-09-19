package enrich

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

var enrichRoutes []lib.Route = []lib.Route{
	{
		Route:   "/vat/munichre/{vat}",
		Handler: lib.ResponseLoggerWrapper(munichVatFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/ateco/{ateco}",
		Handler: lib.ResponseLoggerWrapper(atecoFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/cat-nat/ateco/{fiscalCode}",
		Handler: lib.ResponseLoggerWrapper(catnatAtecoFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},

	{
		Route:   "/cities",
		Handler: lib.ResponseLoggerWrapper(citiesFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},

	{
		Route:   "/works",
		Handler: lib.ResponseLoggerWrapper(worksFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/naics",
		Handler: lib.ResponseLoggerWrapper(naicsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Enrich")
	functions.HTTP("Enrich", enrich)
}

func enrich(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("enrich", enrichRoutes)
	router.ServeHTTP(w, r)
}
