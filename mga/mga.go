package mga

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", Mga)
}

func Mga(w http.ResponseWriter, r *http.Request) {
	log.Println("Mga")
	lib.EnableCors(&w, r)

	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/products/v1",
				Handler: GetProductsListByChannelFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/products/v1",
				Handler: GetProductByChannelFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/products/channel/v1/:channel",
				Handler: GetActiveProductsByChannelFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/network/node/v1/:uid",
				Handler: GetNetworkNodeByUidFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/network/node/v1",
				Handler: CreateNetworkNodeFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/network/node/v1/:uid",
				Handler: DeleteNetworkNodeFx,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/network/invite/v1/create",
				Handler: CreateNetworkNodeInviteFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/network/invite/v1/consume",
				Handler: ConsumeNetworkNodeInviteFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route: "/warrants/v1",
				Handler: GetWarrantsFx,
				Method: http.MethodGet,
				Roles: []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route: "/warrant/v1",
				Handler: CreateWarrantFx,
				Method: http.MethodPut,
				Roles: []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}

	route.Router(w, r)
}
