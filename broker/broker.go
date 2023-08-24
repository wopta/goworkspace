package broker

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {
	log.Println("Broker")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/policies/fiscalcode/:fiscalcode",
				Handler: PolicyFiscalcode,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/:uid",
				Handler: GetPolicyFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/proposal",
				Handler: Proposal,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/emit",
				Handler: EmitFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policy/v1/:uid",
				Handler: UpdatePolicy,
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policy/v1/:uid",
				Handler: DeletePolicy,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "attachment/v1/:policyUid",
				Handler: GetPolicyAttachmentFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policies/v1",
				Handler: GetPoliciesFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "policy/transactions/v1/:policyUid",
				Handler: GetPolicyTransactions,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/policy/reserved/v1/:policyUid",
				Handler: PutPolicyReservedFx,
				Method:  http.MethodPut,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "/policies/auth/v1",
				Handler: GetPoliciesByAuthFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAgent, models.UserRoleAgency},
			},
		},
	}
	route.Router(w, r)
}
