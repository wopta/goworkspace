package mail

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var mailRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/send",
		Fn:          SendFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mail.send",
	},
	{
		Route:       "/v1/score",
		Fn:          ScoreFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mail.score",
	},
	{
		Route:       "/v1/validate",
		Fn:          ValidateFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mail.validate",
	},
}

func init() {
	log.Println("INIT Mail")
	functions.HTTP("Mail", Mail)
}

func Mail(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("mail", mailRoutes)
	router.ServeHTTP(w, r)
}
