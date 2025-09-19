package policy

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/policy/renew"
)

var policyRoutes []lib.Route = []lib.Route{
	{
		Route:   "/fiscalcode/v1/{fiscalcode}",
		Handler: lib.ResponseLoggerWrapper(getPolicyByFiscalCodeFx), // Broker.PolicyFiscalcode
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/annulment/{uid}",
		Handler: lib.ResponseLoggerWrapper(deletePolicyFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/attachment/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyAttachmentsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/media/upload/v1",
		Handler: lib.ResponseLoggerWrapper(uploadPolicyMediaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/media/v1",
		Handler: lib.ResponseLoggerWrapper(getPolicyMediaFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/renew/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(renew.GetRenewPolicyByUidFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/v1/notes/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyNotesFx),
		Method:  http.MethodGet,
		Roles:   []string{},
	},
	{
		Route:   "/v1/note/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(postPolicyNote),
		Method:  http.MethodPost,
		Roles:   []string{},
	},
}

func init() {
	log.Println("INIT Policy")
	functions.HTTP("Policy", policy)
}

func policy(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("policy", policyRoutes)
	router.ServeHTTP(w, r)
}
