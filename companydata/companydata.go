package companydata

import (
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Println("companydata")
	lib.EnableCors(&w, r)
	//w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/global/transactions",
				Handler: GlobalTransaction,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			}, {
				Route:   "/v1/global/pmi/emit",
				Handler: PmiGlobalEmit,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/global/person/emit",
				Handler: PersonGlobalEmit,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/axa/life/emit",
				Handler: LifeAxalEmit,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			}, {
				Route:   "/v1/axa/life/delete",
				Handler: LifeAxaDelete,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/emit",
				Handler: Emit,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}
