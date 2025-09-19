package reserved

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

var reservedRoutes []lib.Route = []lib.Route{
	{
		Route:   "/coverage/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(setCoverageReservedFx),
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
