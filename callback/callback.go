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

func Payment(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	return "", nil
}
