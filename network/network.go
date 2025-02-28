package network

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var networkRoutes []lib.Route = []lib.Route{
	{
		Route:       "/import/v1",
		Method:      http.MethodPost,
		Fn:          ImportNodesFx,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "network.networknodes.import",
	},
	{
		Route:  "/subtree/v1/{nodeUid}",
		Method: http.MethodGet,
		Fn:     NodeSubTreeFx,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgent,
			lib.UserRoleAgency,
		},
		Entitlement: "network.get.networknodes.subtree",
	},
}

func init() {
	log.Println("INIT Network")
	functions.HTTP("Network", Network)
}

func Network(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("network", networkRoutes)
	router.ServeHTTP(w, r)
}
