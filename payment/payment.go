package payment

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
)

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {

	log.Println("Payment")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/fabrick",
				Handler: FabrickPay,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/fabrick/montly",
				Handler: FabrickPayMontly,
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
				Handler: FabrickExpireBill,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/manual/v1/:transactionUid",
				Handler: ManualPaymentFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}
	route.Router(w, r)

}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	return "", nil, nil
}
