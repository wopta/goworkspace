package inclusive

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var inclusiveRoutes []lib.Route = []lib.Route{
	{
		Route:       "/bankaccount/v1/hype",
		Fn:          BankAccountHypeFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "inclusive.hype.bankaccount",
	},
	{
		Route:       "/bankaccount/v1/scalapay",
		Fn:          BankAccountScalapayFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "inclusive.scalapay.bankaccount",
	},
	{
		Route:       "/bankaccount/v1/hype/count",
		Fn:          CountHypeFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "inclusive.hype.bankaccount.count",
	},
	{
		Route:       "/bankaccount/in/v1",
		Fn:          HypeImportMovementbankAccountFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "inclusive.hybe.bankaccount.import",
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
