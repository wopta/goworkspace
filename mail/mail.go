package mail

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var mailRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/send",
		Handler: lib.ResponseLoggerWrapper(SendFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/score",
		Handler: lib.ResponseLoggerWrapper(ScoreFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/validate",
		Handler: lib.ResponseLoggerWrapper(ValidateFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/renew-notice",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(RenewNoticeFx),
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Mail")
	functions.HTTP("Mail", Mail)
}

func Mail(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("mail", mailRoutes)
	router.ServeHTTP(w, r)
}
