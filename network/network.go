package network

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var networkRoutes []lib.Route = []lib.Route{
	{
		Route:   "/import/v1",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(ImportNodesFx),
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/subtree/v1/{nodeUid}",
		Method:  http.MethodGet,
		Handler: lib.ResponseLoggerWrapper(NodeSubTreeFx),
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgent,
			lib.UserRoleAgency,
		},
	},
}

func init() {
	log.Println("INIT Network")
	functions.HTTP("Network", Network)
}

func Network(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("network", networkRoutes)
	router.ServeHTTP(w, r)
}
