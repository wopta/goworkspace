package partnership

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var partnershipRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/life/{partnershipUid}",
		Fn:          LifePartnershipFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "partnership.init.life",
	},
	{
		Route:       "/v1/product/{partnershipUid}",
		Fn:          GetPartnershipNodeAndProductsFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "partnership.get.nodeproduct",
	},
}

func init() {
	log.Println("INIT Partnership")
	functions.HTTP("Partnership", Partnership)
}

func Partnership(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("partnership", partnershipRoutes)
	router.ServeHTTP(w, r)
}
