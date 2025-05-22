package inclusive

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
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

	router := lib.GetRouter("inclusive", inclusiveRoutes)
	router.ServeHTTP(w, r)
}
