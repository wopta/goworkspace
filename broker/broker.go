package broker

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/broker/renew"
	"github.com/wopta/goworkspace/lib"
)

type BrokerBaseRequest struct {
	PolicyUid    string `json:"policyUid"`
	PaymentSplit string `json:"paymentSplit"`
	Payment      string `json:"payment"`
	PaymentMode  string `json:"paymentMode"`
}

var brokerRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/policies/fiscalcode/{fiscalcode}",
		Handler: lib.ResponseLoggerWrapper(PolicyFiscalcodeFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/policy/{uid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/policy/lead",
		Handler: lib.ResponseLoggerWrapper(LeadFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/policy/proposal",
		Handler: lib.ResponseLoggerWrapper(ProposalFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policy/reserved/v1",
		Handler: lib.ResponseLoggerWrapper(RequestApprovalFx),
		Method:  http.MethodPost,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgent,
			lib.UserRoleAgency,
		},
	},
	{
		Route:   "/v1/policy/emit",
		Handler: lib.ResponseLoggerWrapper(EmitFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policy/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(UpdatePolicyFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policy/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(DeletePolicyFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/attachment/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyAttachmentFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policies/v1",
		Handler: lib.ResponseLoggerWrapper(GetPoliciesFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/policy/transactions/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(GetPolicyTransactionsFx),
		Method:  http.MethodGet,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
	},
	{
		Route:   "/policy/reserved/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(AcceptanceFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/policies/auth/v1",
		Handler: lib.ResponseLoggerWrapper(GetPoliciesByAuthFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAgent, lib.UserRoleAgency},
	},
	{
		Route:   "/policy/renew/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(renew.DeleteRenewPolicyFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/transaction/v1/{uid}/receipt",
		Handler: lib.ResponseLoggerWrapper(PaymentReceiptFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
}

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("broker", brokerRoutes)
	router.ServeHTTP(w, r)
}
