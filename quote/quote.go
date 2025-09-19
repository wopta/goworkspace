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
		Handler: lib.ResponseLoggerWrapper(pmiMunichFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/incident",
		Handler: lib.ResponseLoggerWrapper(pmiMunichFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/life",
		Handler: lib.ResponseLoggerWrapper(lifeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/person",
		Handler: lib.ResponseLoggerWrapper(personaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/gap",
		Handler: lib.ResponseLoggerWrapper(gapFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/commercial-combined",
		Handler: lib.ResponseLoggerWrapper(combinedQbeFx),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
	{
		Route:   "/v1/generate/document",
		Handler: lib.ResponseLoggerWrapper(generateDocumentFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/cat-nat",
		Handler: lib.ResponseLoggerWrapper(catNatFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Quote")
	functions.HTTP("Quote", quote)
}

func quote(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("quote", quoteRoutes)
	router.ServeHTTP(w, r)
}
