package companydata

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var companydataRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/global/transactions",
		Fn:          GlobalTransaction,
		Method:      http.MethodGet,
		Roles:       []string{},
		Entitlement: "companydata.global.transactions",
	},
	{
		Route:       "/v1/global/pmi/emit",
		Fn:          PmiGlobalEmit,
		Method:      http.MethodGet,
		Roles:       []string{},
		Entitlement: "companydata.global.pmi.emit",
	},
	{
		Route:       "/v1/global/person/emit",
		Fn:          PersonGlobalEmit,
		Method:      http.MethodGet,
		Roles:       []string{},
		Entitlement: "companydata.global.persona.emit",
	},
	{
		Route:       "/v1/axa/life/emit",
		Fn:          LifeAxaEmit,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "companydata.axa.life.emit",
	},
	{
		Route:       "/v1/sogessur/gap/emit",
		Fn:          GapSogessurEmit,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "companydata.sogessur.gap.emit",
	},
	{
		Route:       "/v1/axa/life/delete",
		Fn:          LifeAxaDelete,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "companydata.axa.life.delete",
	},
	{
		Route:       "/v1/emit",
		Fn:          ProductTrackOutFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "companydata.emit",
	},
	{
		Route:       "/v1/axa/inclusive/bankaccount",
		Fn:          BankAccountInclusive,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "companydata.axa.inclusive.bankaccount",
	},
	{
		Route:       "/v1/in/life",
		Fn:          LifeInFx,
		Method:      http.MethodPost,
		Roles:       []string{models.UserRoleAdmin},
		Entitlement: "companydata.axa.life.import",
	},
}

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("companydata", companydataRoutes)
	router.ServeHTTP(w, r)
}
