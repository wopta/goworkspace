package reserved

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var reservedRoutes []lib.Route = []lib.Route{
	{
		Route:       "/coverage/v1/{policyUid}",
		Fn:          SetCoverageReservedFx,
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "reserved.policy.coverage",
	},
}

func init() {
	log.Println("INIT Reserved")
	functions.HTTP("Reserved", Reserved)
}

func Reserved(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("reserved", reservedRoutes)
	router.ServeHTTP(w, r)
}
