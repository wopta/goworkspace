package callback

/**

workstepFinished : when the workstep was finished
workstepRejected : when the workstep was rejected
workstepDelegated : whe the workstep was delegated
workstepOpened : when the workstep was opened
sendSignNotification : when the sign notification was sent
envelopeExpired : when the envelope was expired
workstepDelegatedSenderActionRequired : when an action from the sender is required because of the delegation
*/
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
			{
				Route:   "/v1/emailVerify",
				Hendler: EmailVerify,
			},
		},
	}
	route.Router(w, r)

}

func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	return "", nil
}
func EmailVerify(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("EmailVerify")
	log.Println("GET params were:", r.URL.Query())

	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	log.Println(token)
	res := lib.WhereFirestore("mail", "email", "==", email)
	objmail, uid := mail.ToListData(res)
	objmail[0].IsValid = true
	lib.SetFirestore("mail", uid[0], objmail)

	return "", nil
}
