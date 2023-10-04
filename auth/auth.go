package auth

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

var origin string

func init() {
	log.Println("INIT Auth")
	functions.HTTP("Auth", Auth)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/authorize/v1",
				Handler: AuthorizeFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/token/v1",
				Handler: TokenFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/sso/jwt/aua/v1",
				Handler: JwtFx,
				Method:  http.MethodGet,
				Roles:   []string{"internal"},
			},
		},
	}
	route.Router(w, r)

}
