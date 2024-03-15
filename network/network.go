package network

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Network")
	functions.HTTP("Network", Network)
}

func Network(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.Println("Network")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/import/v1",
				Method:  http.MethodPost,
				Handler: ImportNodesFx,
				Roles:   []string{models.UserRoleAdmin},
			},
			{
				Route:   "/subtree/v1/:nodeUid",
				Method:  http.MethodGet,
				Handler: NodeSubTreeFx,
				Roles: []string{models.UserRoleAdmin, models.UserRoleManager, models.UserRoleAreaManager, models.UserRoleAgent,
					models.UserRoleAgency},
			},
		},
	}
	route.Router(w, r)
}
