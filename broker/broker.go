package broker

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/broker/renew"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type BrokerBaseRequest struct {
	PolicyUid    string `json:"policyUid"`
	PaymentSplit string `json:"paymentSplit"`
	Payment      string `json:"payment"`
	PaymentMode  string `json:"paymentMode"`
}

var brokerRoutes []lib.Route = []lib.Route{
	{
		Route:       "/v1/policies/fiscalcode/{fiscalcode}",
		Handler:     lib.ResponseLoggerWrapper(PolicyFiscalcodeFx),
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerGetFiscalcode,
	},
	{
		Route:       "/v1/policy/{uid}",
		Handler:     lib.ResponseLoggerWrapper(GetPolicyFx),
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerGetPolicy,
	},
	{
		Route:       "/v1/policy/lead",
		Handler:     lib.ResponseLoggerWrapper(LeadFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Middlewares: []lib.RouteMiddleware{models.MiddlewareCheckEntitlement},
		Entitlement: lib.EntitlementBrokerLead,
	},
	{
		Route:       "/v1/policy/proposal",
		Handler:     lib.ResponseLoggerWrapper(ProposalFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerProposal,
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
		Entitlement: lib.EntitlementBrokerRequestApproval,
	},
	{
		Route:       "/v1/policy/emit",
		Handler:     lib.ResponseLoggerWrapper(EmitFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerEmit,
	},
	{
		Route:       "/policy/v1/{uid}",
		Handler:     lib.ResponseLoggerWrapper(UpdatePolicyFx),
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerPolicyUpdate,
	},
	{
		Route:       "/policy/v1/{uid}",
		Handler:     lib.ResponseLoggerWrapper(DeletePolicyFx),
		Method:      http.MethodDelete,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: lib.EntitlementBrokerDeletePolicy,
	},
	{
		Route:       "/attachment/v1/{policyUid}",
		Handler:     lib.ResponseLoggerWrapper(GetPolicyAttachmentFx),
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerPolicyAttachment,
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
		Entitlement: lib.EntitlementBrokerGetTransactions,
	},
	{
		Route:       "/policy/reserved/v1/{policyUid}",
		Handler:     lib.ResponseLoggerWrapper(AcceptanceFx),
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: lib.EntitlementBrokerAcceptance,
	},
	{
		Route:       "/policy/renew/v1/{uid}",
		Handler:     lib.ResponseLoggerWrapper(renew.DeleteRenewPolicyFx),
		Method:      http.MethodDelete,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: lib.EntitlementBrokerDeleteRenew,
	},
	{
		Route:       "/transaction/v1/{uid}/receipt",
		Handler:     lib.ResponseLoggerWrapper(PaymentReceiptFx),
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
		Entitlement: lib.EntitlementBrokerReceipt,
	},
	{
		Route:       "/policy/v1/init",
		Handler:     lib.ResponseLoggerWrapper(InitFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: lib.EntitlementBrokerInit,
	},
	{
		Route:   "/portfolio/{type}/{version}",
		Method:  http.MethodGet,
		Handler: lib.ResponseLoggerWrapper(GetPortfolioFx),
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency,
			lib.UserRoleAgent},
		Entitlement: lib.EntitlementBrokerGetPortfolio,
	},
	{
		Route:       "/policy/v1/contract/upload/{uid}",
		Handler:     lib.ResponseLoggerWrapper(UploadPolicyContractFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: lib.EntitlementBrokerUploadContract,
	},
	{
		Route:       "/policy/v1/duplicate/{uid}",
		Handler:     lib.ResponseLoggerWrapper(DuplicateFx),
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: lib.EntitlementBrokerDuplicate,
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
