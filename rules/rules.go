package rules

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var rulesRoutes []lib.Route = []lib.Route{
	{
		Route:   "/risk/pmi",
		Handler: lib.ResponseLoggerWrapper(PmiAllriskHandler),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Rules")
	functions.HTTP("Rules", Rules)
}

func Rules(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("rules", rulesRoutes)
	router.ServeHTTP(w, r)
}
