package accounting

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

var accountingRoutes []lib.Route = []lib.Route{
	{
		Route:   "/network/transactions/v1/transaction/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(GetNetworkTransactionsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/transactions/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(PutNetworkTransactionFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/transactions/v1/transaction/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(CreateNetworkTransactionFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
}

func init() {
	log.Println("INIT Accounting")
	functions.HTTP("Accounting", Accounting)
}

func Accounting(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("accounting", accountingRoutes)
	router.ServeHTTP(w, r)
}
