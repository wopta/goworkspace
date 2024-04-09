package companydata

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var companydataRoutes []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/v1/global/transactions",
		Handler: lib.ResponseLoggerWrapper(GlobalTransaction),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/global/pmi/emit",
		Handler: lib.ResponseLoggerWrapper(PmiGlobalEmit),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/global/person/emit",
		Handler: lib.ResponseLoggerWrapper(PersonGlobalEmit),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/axa/life/emit",
		Handler: lib.ResponseLoggerWrapper(LifeAxaEmit),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/sogessur/gap/emit",
		Handler: lib.ResponseLoggerWrapper(GapSogessurEmit),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/axa/life/delete",
		Handler: lib.ResponseLoggerWrapper(LifeAxaDelete),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/emit",
		Handler: lib.ResponseLoggerWrapper(Emit),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/axa/inclusive/bankaccount",
		Handler: lib.ResponseLoggerWrapper(BankAccountAxaInclusive),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/in/life",
		Handler: lib.ResponseLoggerWrapper(LifeInFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("companydata", companydataRoutes)
	router.ServeHTTP(w, r)
}
