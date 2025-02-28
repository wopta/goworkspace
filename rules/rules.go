package rules

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var rulesRoutes []lib.Route = []lib.Route{
	{
		Route:       "/risk/pmi",
		Fn:          PmiAllriskHandler,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "rules.risk.pmi",
	},
}

func init() {
	log.Println("INIT Rules")
	functions.HTTP("Rules", Rules)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("rules", rulesRoutes)
	router.ServeHTTP(w, r)
}
