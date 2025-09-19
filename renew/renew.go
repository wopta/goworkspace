package renew

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

var routes []lib.Route = []lib.Route{
	{
		Route:   "/v1/draft",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(draftFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/promote",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(promoteFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/notice/e-commerce",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(renewMailFx),
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Renew")
	functions.HTTP("Renew", renew)
}

func renew(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("renew", routes)
	router.ServeHTTP(w, r)
}
