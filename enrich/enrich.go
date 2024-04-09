package enrich

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var enrichRoutes []lib.ChiRoute = []lib.ChiRoute{
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
}

func init() {
	log.Println("INIT Enrich")
	functions.HTTP("Enrich", Enrich)
}

func Enrich(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("enrich", enrichRoutes)
	router.ServeHTTP(w, r)
}
