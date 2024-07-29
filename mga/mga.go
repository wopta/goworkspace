package mga

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
)

var mgaRoutes []lib.Route = []lib.Route{
	{
		Route:   "/products/v1",
		Handler: lib.ResponseLoggerWrapper(GetProductsListByChannelFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/products/v1",
		Handler: lib.ResponseLoggerWrapper(GetProductByChannelFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/network/node/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetNetworkNodeByUidFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/network/node/v1",
		Handler: lib.ResponseLoggerWrapper(CreateNetworkNodeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/node/v1",
		Handler: lib.ResponseLoggerWrapper(UpdateNetworkNodeFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/nodes/v1",
		Handler: lib.ResponseLoggerWrapper(GetAllNetworkNodesFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/node/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(DeleteNetworkNodeFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/create",
		Handler: lib.ResponseLoggerWrapper(CreateNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/consume",
		Handler: lib.ResponseLoggerWrapper(ConsumeNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/warrants/v1",
		Handler: lib.ResponseLoggerWrapper(GetWarrantsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/warrant/v1",
		Handler: lib.ResponseLoggerWrapper(CreateWarrantFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/policy/v1",
		Handler: lib.ResponseLoggerWrapper(ModifyPolicyFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/quoter/life/v1",
		Handler: lib.ResponseLoggerWrapper(GetQuoterFileFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
}

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", Mga)
}

func Mga(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("mga", mgaRoutes)
	router.ServeHTTP(w, r)
}
