package renew

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var routes []lib.Route = []lib.Route{
	{
		Route:       "/v1/draft",
		Method:      http.MethodPost,
		Fn:          DraftFx,
		Roles:       []string{},
		Entitlement: "renew.draft",
	},
	{
		Route:       "/v1/promote",
		Method:      http.MethodPost,
		Fn:          PromoteFx,
		Roles:       []string{},
		Entitlement: "renew.promote",
	},
	{
		Route:       "/v1/notice/e-commerce",
		Method:      http.MethodPost,
		Fn:          RenewMailFx,
		Roles:       []string{},
		Entitlement: "renew.notice.ecommerce",
	},
}

func init() {
	log.Println("INIT Renew")
	functions.HTTP("Renew", Renew)
}

func Renew(w http.ResponseWriter, r *http.Request) {
	router := models.GetExtendedRouter("renew", routes)
	router.ServeHTTP(w, r)
}
