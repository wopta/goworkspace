package mga

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mga/consens"
	"github.com/wopta/goworkspace/models"
)

var mgaRoutes []lib.Route = []lib.Route{
	{
		Route:       "/products/v1",
		Fn:          GetProductsListByChannelFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mga.get.products",
	},
	{
		Route:       "/products/v1",
		Fn:          GetProductByChannelFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mga.get.product",
	},
	{
		Route:  "/network/node/v1/{uid}",
		Fn:     GetNetworkNodeByUidFx,
		Method: http.MethodGet,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
		Entitlement: "mga.get.networknode",
	},
	{
		Route:       "/network/node/v1",
		Fn:          CreateNetworkNodeFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.create.networknode",
	},
	{
		Route:       "/network/node/v1",
		Fn:          UpdateNetworkNodeFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.update.networknode",
	},
	{
		Route:       "/network/nodes/v1",
		Fn:          GetAllNetworkNodesFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.get.networknodes",
	},
	{
		Route:       "/network/node/v1/{uid}",
		Fn:          DeleteNetworkNodeFx,
		Method:      http.MethodDelete,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.delete.networknode",
	},
	{
		Route:       "/network/invite/v1/create",
		Fn:          CreateNetworkNodeInviteFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.create.networknode.invite",
	},
	{
		Route:       "/network/invite/v1/consume",
		Fn:          ConsumeNetworkNodeInviteFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "mga.consume.networknode.invite",
	},
	{
		Route:  "/network/consens/v1",
		Fn:     consens.GetUndeclaredConsensFx,
		Method: http.MethodGet,
		Roles: []string{
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
		Entitlement: "mga.get.consens.undeclared",
	},
	{
		Route:  "/network/consens/v1",
		Fn:     consens.AcceptanceFx,
		Method: http.MethodPost,
		Roles: []string{
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
		Entitlement: "mga.give.consent",
	},
	{
		Route:  "/warrants/v1",
		Fn:     GetWarrantsFx,
		Method: http.MethodGet,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAreaManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
		Entitlement: "mga.get.warrants",
	},
	{
		Route:       "/warrant/v1",
		Fn:          CreateWarrantFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "mga.create.warrant",
	},
	{
		Route:       "/policy/v1",
		Fn:          ModifyPolicyFx,
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "mga.modify.policy",
	},
	{
		Route:       "/quoter/life/v1",
		Fn:          GetQuoterFileFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
		Entitlement: "mga.get.quoter.life",
	},
}

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", Mga)
}

func Mga(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("mga", mgaRoutes)
	router.ServeHTTP(w, r)
}
