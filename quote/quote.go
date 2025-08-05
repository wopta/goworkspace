package quote

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
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
		Route:   "/v1/commercial-combined",
		Handler: lib.ResponseLoggerWrapper(CombinedQbeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/generate/document",
		Handler: lib.ResponseLoggerWrapper(GenerateDocumentFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/cat-nat",
		Handler: lib.ResponseLoggerWrapper(CatNatFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Quote")
	functions.HTTP("Quote", Quote)
}

func Quote(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("quote", quoteRoutes)
	router.ServeHTTP(w, r)
}
