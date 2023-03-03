package callback

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
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
				Handler: Sign,
				Method:  "GET",
			},
			{
				Route:   "/v1/payment",
				Handler: Payment,
				Method:  http.MethodPost,
			},
			{
				Route:   "/v1/emailVerify",
				Handler: EmailVerify,
				Method:  "GET",
			},
		},
	}
	route.Router(w, r)

}

func EmailVerify(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("EmailVerify")
	log.Println("GET params were:", r.URL.Query())

	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	log.Println(token)
	res := lib.WhereFirestore("mail", "email", "==", email)
	objmail, uid := mail.ToListData(res)
	objmail[0].IsValid = true
	lib.SetFirestore("mail", uid[0], objmail)

	return "", nil, nil
}
