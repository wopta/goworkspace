package test

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
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
	},{
		Route:   "/namirial/envelop",
		Handler: lib.ResponseLoggerWrapper(HandlerEnvelop),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Test")
	functions.HTTP("Test", Test)
}

func Test(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("test", testRoutes)
	router.ServeHTTP(w, r)
}
