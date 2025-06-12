package callback

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var callbackRoutes []lib.Route = []lib.Route{
	//	{
	//		Route:   "/v1/sign",
	//		Handler: lib.ResponseLoggerWrapper(SignFx),
	//		Method:  http.MethodGet,
	//		Roles:   []string{},
	//	},
	{
		Route:   "/v1/sign",
		Handler: lib.ResponseLoggerWrapper(DraftSignFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	//TODO: could be this eliminated?
	//	{
	//		Route:   "/v1/payment",
	//		Handler: lib.ResponseLoggerWrapper(PaymentFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{},
	//	},
	{
		//	Route: "/v1/payment/{provider}/first-rate",
		//Route: "/v1/payment/{provider}/single-rate",
		Route: "/v1/payment/{provider}/{rate}",
		// TODO: create an extra handler wrapper that switches on provider.
		// For now as fabrick is the single provider it is hardcoded.
		Handler: lib.ResponseLoggerWrapper(payment),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	//	{
	//		Route: "/v1/payment/{provider}/single-rate",
	//		// TODO: create an extra handler wrapper that switches on provider.
	//		// For now as fabrick is the single provider it is hardcoded.
	//		Handler: lib.ResponseLoggerWrapper(fabrick.AnnuitySingleRateFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{},
	//	},
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

	router := lib.GetRouter("callback", callbackRoutes)
	router.ServeHTTP(w, r)
}
