package accounting

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Accounting")
	functions.HTTP("Accounting", Accounting)
}

func Accounting(w http.ResponseWriter, r *http.Request) {
	log.Println("Accounting")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/network/transactions/v1/transaction/:transactionUid",
				Handler: GetNetworkTransactions,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/network/transactions/v1/:uid",
				Handler: PutNetworkTransaction,
				Method:  http.MethodPut,
				// Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
				Roles: []string{},
			},
			// POST
		},
	}

	route.Router(w, r)
}
