package inclusive

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var inclusiveRoutes []lib.Route = []lib.Route{
	{
		Route:   "/bankaccount/v1/hype",
		Handler: lib.ResponseLoggerWrapper(BankAccountHypeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/bankaccount/v1/scalapay",
		Handler: lib.ResponseLoggerWrapper(BankAccountScalapayFx),
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

	router := lib.GetRouter("inclusive", inclusiveRoutes)
	router.ServeHTTP(w, r)
}
