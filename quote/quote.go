package quote

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var quoteRoutes []lib.Route = []lib.Route{
	{
		Route:       "/pmi/munichre",
		Fn:          PmiMunichFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "quote.pmi.munichre",
	},
	{
		Route:       "/incident",
		Fn:          PmiMunichFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "quote.pmi.incident",
	},
	{
		Route:       "/v1/life",
		Fn:          LifeFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "quote.life",
	},
	{
		Route:       "/v1/person",
		Fn:          PersonaFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "quote.person",
	},
	{
		Route:       "/v1/gap",
		Fn:          GapFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "quote.gap",
	},
	{
		Route:       "/v1/combined",
		Fn:          CombinedQbeFx,
		Method:      http.MethodPost,
		Roles:       []string{},
		Entitlement: "quote.commercialcombined",
	},
}

func init() {
	log.Println("INIT Quote")
	functions.HTTP("Quote", Quote)
}

func Quote(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("quote", quoteRoutes)
	router.ServeHTTP(w, r)
}
