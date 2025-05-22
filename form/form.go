package form

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var formRoutes []lib.Route = []lib.Route{
	{
		Route:   "/axafleet",
		Handler: lib.ResponseLoggerWrapper(AxaFleetTway),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(FleetAssistenceInclusiveMovement),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/fleet/assistance/v1",
		Handler: lib.ResponseLoggerWrapper(FleetAssistenceInclusiveMovement),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Form")
	functions.HTTP("Form", Form)
}

func Form(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("form", formRoutes)
	router.ServeHTTP(w, r)
}
