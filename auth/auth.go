package auth

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var authRoutes []lib.Route = []lib.Route{
	{
		Route:   "/authorize/v1",
		Handler: lib.ResponseLoggerWrapper(AuthorizeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/token/v1",
		Handler: lib.ResponseLoggerWrapper(TokenFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/sso/jwt/{provider}/v1",
		Handler: lib.ResponseLoggerWrapper(DynamicJwtFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleInternal},
	},
	{
		Route:   "/sso/external/v1/{productName}",
		Handler: lib.ResponseLoggerWrapper(GetTokenForExternalIntegrationFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAgent, lib.UserRoleAgency},
	},
}

var origin string

func init() {
	log.Println("INIT Auth")
	functions.HTTP("Auth", Auth)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("auth", authRoutes)
	router.ServeHTTP(w, r)
}
