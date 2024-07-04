package renew

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var routes []lib.Route = []lib.Route{
	{
		Route:   "/v1/draft",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(DraftFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/promote",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(PromoteFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/policies",
		Method:  http.MethodGet,
		Handler: lib.ResponseLoggerWrapper(GetRenewPoliciesFx),
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgent, lib.UserRoleAgency},
	},
}

func init() {
	log.Println("INIT Renew")
	functions.HTTP("Renew", Renew)
}

func Renew(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("renew", routes)
	router.ServeHTTP(w, r)
}
