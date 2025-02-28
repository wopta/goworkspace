package accounting

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var accountingRoutes []lib.Route = []lib.Route{
	{
		Route:       "/network/transactions/v1/transaction/{transactionUid}",
		Fn:          GetNetworkTransactionsFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "accounting.get.networktransactions",
	},
	{
		Route:       "/network/transactions/v1/{uid}",
		Fn:          PutNetworkTransactionFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "accounting.put.networktransaction",
	},
	{
		Route:       "/network/transactions/v1/transaction/{transactionUid}",
		Fn:          CreateNetworkTransactionFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "accounting.create.networktransaction",
	},
}

func init() {
	log.Println("INIT Accounting")
	functions.HTTP("Accounting", Accounting)
}

func Accounting(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("accounting", accountingRoutes)
	router.ServeHTTP(w, r)
}
