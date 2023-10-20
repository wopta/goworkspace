package reserved

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Reserved")
	functions.HTTP("Reserved", Reserved)
}

func Reserved(w http.ResponseWriter, r *http.Request) {
	log.Println("Reserved")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/coverage/v1/:policyUid",
				Handler: SetCoverageReservedFx,
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}
