package payment

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment/fabrick"
	"github.com/wopta/goworkspace/payment/manual"
)

var paymentRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/fabrick/recreate",
		Handler: lib.ResponseLoggerWrapper(fabrick.RefreshPayByLinkFx),
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
		Handler: lib.ResponseLoggerWrapper(manual.ManualPaymentFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(ChangePaymentProviderFx),
		Method:  http.MethodPatch,
		Roles:   []string{models.UserRoleAdmin},
	},
	{
		Route:   "/v1/fabrick/refresh-token",
		Handler: lib.ResponseLoggerWrapper(fabrick.RefreshTokenFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("payment", paymentRoutes)
	router.ServeHTTP(w, r)
}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}
