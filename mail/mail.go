package mail

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var mailRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/send",
		Handler: lib.ResponseLoggerWrapper(sendFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/score",
		Handler: lib.ResponseLoggerWrapper(scoreFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/validate",
		Handler: lib.ResponseLoggerWrapper(validateFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Mail")
	functions.HTTP("Mail", mail)
}

func mail(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("mail", mailRoutes)
	router.ServeHTTP(w, r)
}
