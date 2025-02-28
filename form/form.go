package form

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var formRoutes []lib.Route = []lib.Route{
	{
		Route:       "/axafleet",
		Fn:          AxaFleetTway,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "form.axa.fleet",
	},
	{
		Route:       "/v1/{uid}",
		Fn:          FleetAssistenceInclusiveMovement,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "form.fleet.assistance",
	},
	{
		Route:       "fleet/assistance/v1",
		Fn:          FleetAssistenceInclusiveMovement,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "form.fleet.assistance",
	},
}

func init() {
	log.Println("INIT Form")
	functions.HTTP("Form", Form)
}

func Form(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("form", formRoutes)
	router.ServeHTTP(w, r)
}
