package product

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
	//"google.golang.org/api/firebaseappcheck/v1"
)

func init() {
	log.Println("INIT Product")
	functions.HTTP("Product", Product)
}

func Product(w http.ResponseWriter, r *http.Request) {

	log.Println("Product")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1",
				Handler: PutFx,
				Method:  "PUT",
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}
