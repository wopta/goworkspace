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
		Route:       "/v1/fabrick/recreate",
		Fn:          fabrick.RefreshPayByLinkFx,
		Method:      http.MethodPost,
		Roles:       []string{models.UserRoleAdmin, models.UserRoleManager},
		Entitlement: "payment.fabrick.recreate.link",
	},
	{
		Route:       "/v1/cripto",
		Fn:          CriptoPay,
		Method:      http.MethodPost,
		Roles:       []string{models.UserRoleAll},
		Entitlement: "payment.crypto",
	},
	{
		Route:       "/v1/{uid}",
		Fn:          DeleteTransactionFx,
		Method:      http.MethodDelete,
		Roles:       []string{models.UserRoleAdmin, models.UserRoleManager},
		Entitlement: "payment.delete.transaction",
	},
	{
		Route:  "/manual/v1/{transactionUid}",
		Fn:     manual.ManualPaymentFx,
		Method: http.MethodPost,
		Roles: []string{
			models.UserRoleAdmin,
			models.UserRoleManager,
			models.UserRoleAreaManager,
			models.UserRoleAgency,
			models.UserRoleAgent,
		},
		Entitlement: "payment.pay.manual.transacation",
	},
	{
		Route:  "/manual/v1/renew/{transactionUid}",
		Fn:     manual.RenewManualPaymentFx,
		Method: http.MethodPost,
		Roles: []string{
			models.UserRoleAdmin,
			models.UserRoleManager,
			models.UserRoleAreaManager,
			models.UserRoleAgency,
			models.UserRoleAgent,
		},
		Entitlement: "payment.pay.manual.renew.transacation",
	},
	{
		Route:       "/v1",
		Fn:          ChangePaymentProviderFx,
		Method:      http.MethodPatch,
		Roles:       []string{models.UserRoleAdmin},
		Entitlement: "payment.change.provider",
	},
	{
		Route:       "/v1/renew",
		Fn:          RenewChangePaymentProviderFx,
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "payment.change.provider.renew",
	},
	{
		Route:       "/v1/fabrick/refresh-token",
		Fn:          fabrick.RefreshTokenFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "payment.fabrick.refresh.token",
	},
}

func init() {
	log.Println("INIT Payment")
	functions.HTTP("Payment", Payment)
}

func Payment(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("payment", paymentRoutes)
	router.ServeHTTP(w, r)
}

func CriptoPay(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "", nil, nil
}
