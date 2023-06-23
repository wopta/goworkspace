package mga

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", Mga)
}

func Mga(w http.ResponseWriter, r *http.Request) {
	log.Println("Mga")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/products/v1",
				Handler: func(w http.ResponseWriter, r *http.Request) (string, interface{}, error) { return "", nil, nil },
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/products/v1",
				Handler: GetProductByRoleFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}

	route.Router(w, r)
}
