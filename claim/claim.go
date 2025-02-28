package claim

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var claimRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1",
		Fn:          PutClaimFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleCustomer},
		Entitlement: "claim.create",
	},
	{
		Route:       "/document/v1/{claimUid}",
		Fn:          GetClaimDocumentFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleCustomer},
		Entitlement: "claim.get.attachment",
	},
}

func init() {
	log.Println("INIT Claim")
	functions.HTTP("Claim", Claim)
}

func Claim(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("claim", claimRoutes)
	router.ServeHTTP(w, r)
}
