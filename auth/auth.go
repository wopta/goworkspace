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
	log.Println("INIT AppcheckProxy")
	functions.HTTP("Auth", Auth)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/authorize/v1/",
				Handler: AuthorizeFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/token/v1/",
				Handler: TokenFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/token/v1/sso/jwt/aua",
				Handler: TokenFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}
