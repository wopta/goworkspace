package callback

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/callback/fabrick"
	"github.com/wopta/goworkspace/lib"
)

var callbackRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/sign",
		Handler: lib.ResponseLoggerWrapper(SignFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/payment",
		Handler: lib.ResponseLoggerWrapper(PaymentFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route: "/v1/payment/{provider}/first-rate",
		// TODO: create an extra handler wrapper that switches on provider.
		// For now as fabrick is the single provbider it is hardcoded.
		Handler: lib.ResponseLoggerWrapper(fabrick.AnnuityFirstRateFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route: "/v1/payment/{provider}/single-rate",
		// TODO: create an extra handler wrapper that switches on provider.
		// For now as fabrick is the single provbider it is hardcoded.
		Handler: lib.ResponseLoggerWrapper(fabrick.AnnuitySingleRateFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/emailVerify",
		Handler: lib.ResponseLoggerWrapper(EmailVerifyFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Callback")
	functions.HTTP("Callback", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("callback", callbackRoutes)
	router.ServeHTTP(w, r)
}
