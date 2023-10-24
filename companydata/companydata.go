package companydata

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT companydata")
	functions.HTTP("Companydata", Callback)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	log.Println("companydata")
	lib.EnableCors(&w, r)
	// w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/global/transactions",
				Handler: GlobalTransaction,
				Method:  "GET",
				Roles:   []string{},
			},
			{
				Route:   "/v1/global/pmi/emit",
				Handler: PmiGlobalEmit,
				Method:  "GET",
				Roles:   []string{},
			},
			{
				Route:   "/v1/global/person/emit",
				Handler: PersonGlobalEmit,
				Method:  "GET",
				Roles:   []string{},
			},
			{
				Route:   "/v1/axa/life/emit",
				Handler: LifeAxaEmit,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/sogessur/gap/emit",
				Handler: GapSogessurEmit,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/axa/life/delete",
				Handler: LifeAxaDelete,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/emit",
				Handler: Emit,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/axa/inclusive/bankaccount",
				Handler: BankAccountAxaInclusive,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
			{
				Route:   "/v1/in/life",
				Handler: LifeIn,
				Method:  http.MethodPost,
				Roles:   []string{},
			},
		},
	}
	route.Router(w, r)
}
