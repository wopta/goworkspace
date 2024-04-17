package claim

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var claimRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1",
		Handler: lib.ResponseLoggerWrapper(PutClaimFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleCustomer},
	},
	{
		Route:   "/document/v1/{claimUid}",
		Handler: lib.ResponseLoggerWrapper(GetClaimDocumentFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleCustomer},
	},
}

func init() {
	log.Println("INIT Claim")
	functions.HTTP("Claim", Claim)
}

func Claim(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("claim", claimRoutes)
	router.ServeHTTP(w, r)
}
