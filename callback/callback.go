package callback

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/callback/fabrick"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var callbackRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/sign",
		Fn:          SignFx,
		Method:      http.MethodGet,
		Roles:       []string{},
		Entitlement: "callback.sign",
	},
	{
		Route:       "/v1/payment",
		Fn:          PaymentFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "callback.payment",
	},
	{
		Route: "/v1/payment/{provider}/first-rate",
		// TODO: create an extra handler wrapper that switches on provider.
		// For now as fabrick is the single provider it is hardcoded.
		Fn:          fabrick.AnnuityFirstRateFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "callback.payment.firstrate",
	},
	{
		Route: "/v1/payment/{provider}/single-rate",
		// TODO: create an extra handler wrapper that switches on provider.
		// For now as fabrick is the single provider it is hardcoded.
		Fn:          fabrick.AnnuitySingleRateFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "callback.payment.singlerate",
	},
	{
		Route:       "/v1/emailVerify",
		Fn:          EmailVerifyFx,
		Method:      http.MethodGet,
		Roles:       []string{},
		Entitlement: "callback.email.verify",
	},
}

func init() {
	log.Println("INIT Callback")
	functions.HTTP("Callback", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("callback", callbackRoutes)
	router.ServeHTTP(w, r)
}
