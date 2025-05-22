package test

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var testRoutes []lib.Route = []lib.Route{
	{
		Route:   "/{operation}",
		Handler: lib.ResponseLoggerWrapper(TestPostFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/{operation}",
		Handler: lib.ResponseLoggerWrapper(TestGetFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/fabrick/{operation}",
		Handler: lib.ResponseLoggerWrapper(TestFabrickFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	}, {
		Route:   "/scalapay/import",
		Handler: lib.ResponseLoggerWrapper(ImportScalapay),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/accounting/createinvoice",
		Handler: lib.ResponseLoggerWrapper(createInvoice),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/log/{severity}/{message}",
		Handler: lib.ResponseLoggerWrapper(logFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("test", testRoutes)
	router.ServeHTTP(w, r)
}
