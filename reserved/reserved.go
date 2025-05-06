package reserved

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
)

var reservedRoutes []lib.Route = []lib.Route{
	{
		Route:   "/coverage/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(SetCoverageReservedFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Reserved")
	functions.HTTP("Reserved", Reserved)
}

func Reserved(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("reserved", reservedRoutes)
	router.ServeHTTP(w, r)
}
