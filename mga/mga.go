package mga

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mga/consens"
)

var mgaRoutes []lib.Route = []lib.Route{
	{
		Route:   "/products/v1",
		Handler: lib.ResponseLoggerWrapper(getProductsListByChannelFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/products/v1",
		Handler: lib.ResponseLoggerWrapper(getProductByChannelFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/network/node/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(getNetworkNodeByUidFx),
		Method:  http.MethodGet,
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency,
			lib.UserRoleAgent},
	},
	{
		Route:   "/network/node/v1",
		Handler: lib.ResponseLoggerWrapper(createNetworkNodeFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/node/v1",
		Handler: lib.ResponseLoggerWrapper(updateNetworkNodeFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/nodes/v1",
		Handler: lib.ResponseLoggerWrapper(getAllNetworkNodesFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/node/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(deleteNetworkNodeFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/create",
		Handler: lib.ResponseLoggerWrapper(createNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/network/invite/v1/consume",
		Handler: lib.ResponseLoggerWrapper(consumeNetworkNodeInviteFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/network/consens/v1",
		Handler: lib.ResponseLoggerWrapper(consens.GetUndeclaredConsensFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/network/consens/v1",
		Handler: lib.ResponseLoggerWrapper(consens.AcceptanceFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/warrants/v1",
		Handler: lib.ResponseLoggerWrapper(getWarrantsFx),
		Method:  http.MethodGet,
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency,
			lib.UserRoleAgent},
	},
	{
		Route:   "/warrant/v1",
		Handler: lib.ResponseLoggerWrapper(createWarrantFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/policy/v1",
		Handler: lib.ResponseLoggerWrapper(modifyPolicyFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/quoter/life/v1",
		Handler: lib.ResponseLoggerWrapper(getQuoterFileFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/refund/policy/{policyUid}/{transactionUid}",
		Handler: lib.ResponseLoggerWrapper(refundPolicyFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin},
	},
}

func init() {
	log.Println("INIT Mga")
	functions.HTTP("Mga", mga)
}

func mga(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("mga", mgaRoutes)
	router.ServeHTTP(w, r)
}
