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
		Route:   "/v1/annulment/{uid}",
		Handler: lib.ResponseLoggerWrapper(DeletePolicyFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/attachment/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyAttachmentsFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
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
		Handler: lib.ResponseLoggerWrapper(getPolicyNotes),
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
	functions.HTTP("Policy", Policy)
}

func Policy(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("policy", policyRoutes)
	router.ServeHTTP(w, r)
}
