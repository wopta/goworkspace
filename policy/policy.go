package policy

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/policy/renew"
)

var policyRoutes []lib.Route = []lib.Route{
	{
		Route:       "/fiscalcode/v1/{fiscalcode}",
		Fn:          GetPolicyByFiscalCodeFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "policy.get.policy.fiscalcode",
	},
	{
		Route:       "/v1/{uid}",
		Fn:          GetPolicyFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "policy.get.policy",
	},
	{
		Route:       "/v1/annulment/{uid}",
		Fn:          DeletePolicyFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "policy.delete.policy",
	},
	{
		Route:       "/attachment/v1/{uid}",
		Fn:          GetPolicyAttachmentsFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "policy.get.policy.attachments",
	},
	{
		Route:       "/media/upload/v1",
		Fn:          UploadPolicyMediaFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "policy.upload.policy.media",
	},
	{
		Route:       "/media/v1",
		Fn:          GetPolicyMediaFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
		Entitlement: "policy.get.policy.media",
	},
	{
		Route:       "/renew/v1/{uid}",
		Fn:          renew.GetRenewPolicyByUidFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAgency, lib.UserRoleAgent},
		Entitlement: "policy.get.policy.renew",
	},
}

func init() {
	log.Println("INIT Policy")
	functions.HTTP("Policy", Policy)
}

func Policy(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("policy", policyRoutes)
	router.ServeHTTP(w, r)
}
