package callback

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Callback")
	functions.HTTP("Callback", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/sign",
				Handler: SignFx,
				Method:  http.MethodGet,
				Roles:   []string{},
			},
			{
				Route:   "/v1/payment",
				Handler: PaymentFx,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v2/payment",
				Handler: PaymentV2Fx,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/emailVerify",
				Handler: EmailVerify,
				Method:  http.MethodGet,
				Roles:   []string{},
			},
		},
	}
	route.Router(w, r)
}
