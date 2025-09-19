package payment

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/fabrick"
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
		Handler: lib.ResponseLoggerWrapper(criptoPayFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(deleteTransactionFx),
		Method:  http.MethodDelete,
		Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
	},
	{
		Route:   "/manual/v1/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(manualPaymentFx),
		Method:  http.MethodPost,
		Roles: []string{models.UserRoleAdmin, models.UserRoleManager, models.UserRoleAreaManager,
			models.UserRoleAgency, models.UserRoleAgent},
	},
	{
		Route:   "/manual/v1/renew/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(renewManualPaymentFx),
		Method:  http.MethodPost,
		Roles: []string{models.UserRoleAdmin, models.UserRoleManager, models.UserRoleAreaManager,
			models.UserRoleAgency, models.UserRoleAgent},
	},
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(changePaymentProviderFx),
		Method:  http.MethodPatch,
		Roles:   []string{models.UserRoleAdmin},
	},
	{
		Route:   "/v1/renew",
		Handler: lib.ResponseLoggerWrapper(renewChangePaymentProviderFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAdmin},
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
	functions.HTTP("Payment", payment)
}

func payment(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("payment", paymentRoutes)
	router.ServeHTTP(w, r)
}

func criptoPayFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}
