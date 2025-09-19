package callback

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var callbackRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/net/incasso",
		Handler: lib.ResponseLoggerWrapper(incassoNetFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/sign",
		Handler: lib.ResponseLoggerWrapper(signFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/payment/{provider}/{rate}",
		Handler: lib.ResponseLoggerWrapper(paymentFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/emailVerify",
		Handler: lib.ResponseLoggerWrapper(emailVerifyFx),
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
