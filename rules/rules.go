package rules

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var rulesRoutes []lib.ChiRoute = []lib.ChiRoute{
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("rules", rulesRoutes)
	router.ServeHTTP(w, r)
}
