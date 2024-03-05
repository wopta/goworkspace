package partnership

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Partnership")
	functions.HTTP("Partnership", Partnership)
}

func Partnership(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/life/:partnershipUid",
				Handler: LifePartnershipFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)
}

type PartnershipResponse struct {
	Policy      models.Policy          `json:"policy"`
	Partnership models.PartnershipNode `json:"partnership"`
	Product     models.Product         `json:"product"`
}
