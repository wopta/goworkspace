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
		Handler: lib.ResponseLoggerWrapper(MunichVatFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/ateco/{ateco}",
		Handler: lib.ResponseLoggerWrapper(AtecoFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/cat-nat/ateco/{fiscalCode}",
		Handler: lib.ResponseLoggerWrapper(CatnatAtecoFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},

	{
		Route:   "/cities",
		Handler: lib.ResponseLoggerWrapper(CitiesFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},

	{
		Route:   "/works",
		Handler: lib.ResponseLoggerWrapper(WorksFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/naics",
		Handler: lib.ResponseLoggerWrapper(NaicsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Enrich")
	functions.HTTP("Enrich", Enrich)
}

func Enrich(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("enrich", enrichRoutes)
	router.ServeHTTP(w, r)
}
