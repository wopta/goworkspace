package policy

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/policy/renew"
)

var policyRoutes []lib.Route = []lib.Route{
	{
		Route:   "/fiscalcode/v1/{fiscalcode}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyByFiscalCodeFx), // Broker.PolicyFiscalcode
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(DeletePolicyFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/attachment/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyAttachmentsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/node/v1/{nodeUid}",
		Handler: lib.ResponseLoggerWrapper(GetNodePoliciesFx),
		Method:  http.MethodPost,
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager,
			lib.UserRoleAgent, lib.UserRoleAgency},
	},
	{
		Route:   "/portfolio/v1",
		Handler: lib.ResponseLoggerWrapper(GetPortfolioPoliciesFx),
		Method:  http.MethodPost,
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager,
			lib.UserRoleAgent, lib.UserRoleAgency},
	},
	{
		Route:   "/media/upload/v1",
		Handler: lib.ResponseLoggerWrapper(UploadPolicyMediaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/media/v1",
		Handler: lib.ResponseLoggerWrapper(GetPolicyMediaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgent, lib.UserRoleAgency},
	},
	{
		Route:   "/renewed/v1",
		Handler: lib.ResponseLoggerWrapper(renew.GetRenewedPoliciesFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/renewed/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(renew.GetRenewPolicyByUidFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgent, lib.UserRoleAgency},
	},
}

func init() {
	log.Println("INIT Policy")
	functions.HTTP("Policy", Policy)
}

func Policy(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("policy", policyRoutes)
	router.ServeHTTP(w, r)
}
