package partnership

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

var partnershipRoutes []lib.ChiRoute = []lib.ChiRoute{
	{
		Route:   "/v1/life/{partnershipUid}",
		Handler: lib.ResponseLoggerWrapper(LifePartnershipFx),
		Method:  http.MethodGet,
		Roles:   []string{models.UserRoleAll},
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

type PartnershipResponse struct {
	Policy      models.Policy          `json:"policy"`
	Partnership models.PartnershipNode `json:"partnership"`
	Product     models.Product         `json:"product"`
}
