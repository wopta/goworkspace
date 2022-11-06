package broker

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Callback")
	functions.HTTP("Callback", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback")
	lib.EnableCors(&w, r)
	//w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/sign",
				Hendler: Sign,
			},
			{
				Route:   "/v1/payment",
				Hendler: Payment,
			},
		},
	}
	route.Router(w, r)

}

func Sign(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("Sign")
	log.Println("GET params were:", r.URL.Query())

	envelope := r.URL.Query().Get("envelope")
	action := r.URL.Query().Get("action")
	log.Println(action)
	log.Println(envelope)
	if envelope != "" {
		// ... process it, will be the first (only) if multiple were given
		// note: if they pass in like ?param1=&param2= param1 will also be "" :|
	}

	return "", nil
}
func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	return "", nil
}
