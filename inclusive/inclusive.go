package inclusive

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var inclusiveRoutes []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/bankaccount/v1/hype",
		Handler: lib.ResponseLoggerWrapper(BankAccountHypeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/bankaccount/v1/hype/count",
		Handler: lib.ResponseLoggerWrapper(CountHypeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/bankaccount/in/v1",
		Handler: lib.ResponseLoggerWrapper(HypeImportMovementbankAccountFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Inclusive")
	functions.HTTP("Inclusive", InclusiveFx)
}

func InclusiveFx(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("inclusive", inclusiveRoutes)
	router.ServeHTTP(w, r)
}
