package payment

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	log.Println("Payment")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/fabrick/recreate",
				Handler: FabrickRefreshPayByLinkFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/v1/fabrick",
				Handler: FabrickPayFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMonthlyFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/cripto",
				Handler: CriptoPay,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/:uid",
				Handler: DeleteTransactionFx,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/manual/v1/:transactionUid",
				Handler: ManualPaymentFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/v1",
				Handler: ChangePaymentProviderFx,
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}
