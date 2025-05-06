package renew

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
)

var routes []lib.Route = []lib.Route{
	{
		Route:   "/v1/draft",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(DraftFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/promote",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(PromoteFx),
		Roles:   []string{},
	},
	{
		Route:   "/v1/notice/e-commerce",
		Method:  http.MethodPost,
		Handler: lib.ResponseLoggerWrapper(RenewMailFx),
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Renew")
	functions.HTTP("Renew", Renew)
}

func Renew(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("renew", routes)
	router.ServeHTTP(w, r)
}
