package partnership

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var partnershipRoutes []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/v1/life/{partnershipUid}",
		Handler: lib.ResponseLoggerWrapper(LifePartnershipFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/auth/{partnershipUid}",
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetChiRouter("partnership", partnershipRoutes)
	router.ServeHTTP(w, r)
}
