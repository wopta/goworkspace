package inclusive

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Inclusive")
	functions.HTTP("Inclusive", InclusiveFx)
}

func InclusiveFx(w http.ResponseWriter, r *http.Request) {

	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	log.Println(r.RequestURI)
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/bankaccount/v1/hype",
				Handler: BankAccountHypeFx,
				Method:  "POST",
				Roles:   []string{},
			},
			{
				Route:   "/bankaccount/v1/hype/count",
				Handler: CountHypeFx,
				Method:  "POST",
				Roles:   []string{},
			},
			{
				Route:   "/bankaccount/in/v1",
				Handler: HypeImportMovementbankAccountFx,
				Method:  "POST",
				Roles:   []string{},
			},
		},
	}
	route.Router(w, r)

}
