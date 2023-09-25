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
				Route:   "/v1/policy/lead",
				Handler: LeadFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/proposal",
				Handler: ProposalFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/policy/reserved/v1",
				Handler: RequestApprovalFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager, models.AgentChannel},
			},
			{
				Route:   "/v1/policy/emit",
				Handler: EmitFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policy/v1/:uid",
				Handler: UpdatePolicyFx,
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
				Handler: GetPolicyTransactionsFx,
				Method:  http.MethodGet,
				Roles: []string{
					models.UserRoleAdmin,
					models.UserRoleManager,
					models.UserRoleAgency,
					models.UserRoleAgent,
				},
			},
			{
				Route:   "/policy/reserved/v1/:policyUid",
				Handler: AcceptanceFx,
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
