package quote

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var quoteRoutes []lib.Route = []lib.Route{
	{
		Route:   "/pmi/munichre",
		Handler: lib.ResponseLoggerWrapper(PmiMunichFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/incident",
		Handler: lib.ResponseLoggerWrapper(PmiMunichFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/life",
		Handler: lib.ResponseLoggerWrapper(LifeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/person",
		Handler: lib.ResponseLoggerWrapper(PersonaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/gap",
		Handler: lib.ResponseLoggerWrapper(GapFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/excel",
		Handler: lib.ResponseLoggerWrapper(CombinedQbeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Quote")
	functions.HTTP("Quote", Quote)
}

func Quote(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("quote", quoteRoutes)
	router.ServeHTTP(w, r)
}
