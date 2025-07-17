package companydata

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

var companydataRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/global/transactions",
		Handler: lib.ResponseLoggerWrapper(GlobalTransaction),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/generate/track/{operation}/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(GenerateTrackFx),
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
		Handler: lib.ResponseLoggerWrapper(ProductTrackOutFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/axa/inclusive/bankaccount",
		Handler: lib.ResponseLoggerWrapper(BankAccountInclusive),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/in/life",
		Handler: lib.ResponseLoggerWrapper(LifeInFx),
		Method:  http.MethodPost,
		Roles:   []string{models.UserRoleAdmin},
	},
}

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("companydata", companydataRoutes)
	router.ServeHTTP(w, r)
}
