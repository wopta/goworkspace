package auth

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var authRoutes []lib.Route = []lib.Route{
	{
		Route:       "/authorize/v1",
		Fn:          AuthorizeFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "auth.authorize",
	},
	{
		Route:       "/token/v1",
		Fn:          TokenFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "auth.token",
	},
	{
		Route:       "/sso/jwt/{provider}/v1",
		Fn:          DynamicJwtFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleInternal},
		Entitlement: "auth.sso.jwt",
	},
	{
		Route:       "/sso/external/v1/{productName}",
		Fn:          GetTokenForExternalIntegrationFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAgent, lib.UserRoleAgency},
		Entitlement: "auth.sso.external.product",
	},
}

var origin string

func init() {
	log.Println("INIT Auth")
	functions.HTTP("Auth", Auth)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("auth", authRoutes)
	router.ServeHTTP(w, r)
}
