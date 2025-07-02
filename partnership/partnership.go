package partnership

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
)

var partnershipRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/life/{partnershipUid}",
		Handler: lib.ResponseLoggerWrapper(LifePartnershipFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/life/{partnershipUid}",
		Handler: lib.ResponseLoggerWrapper(NewLifePartnershipFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/product/{partnershipUid}",
		Handler: lib.ResponseLoggerWrapper(GetPartnershipNodeAndProductsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Partnership")
	functions.HTTP("Partnership", Partnership)
}

func Partnership(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("partnership", partnershipRoutes)
	router.ServeHTTP(w, r)
}
