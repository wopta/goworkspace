package payment

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var paymentRouts []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/v1/fabrick/recreate",
		Handler: lib.ResponseLoggerWrapper(FabrickRefreshPayByLinkFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/v1/cripto",
		Handler: lib.ResponseLoggerWrapper(CriptoPay),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(DeleteTransactionFx),
		Method:  http.MethodDelete,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/manual/v1/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(ManualPaymentFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(ChangePaymentProviderFx),
		Method:  http.MethodPatch,
		Roles:   []string{models.UserRoleAdmin},
	},
}

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("payment", paymentRouts)
	router.ServeHTTP(w, r)
}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}
