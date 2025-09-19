package broker

import (
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"gitlab.dev.wopta.it/goworkspace/broker/renew"
	"gitlab.dev.wopta.it/goworkspace/lib"
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
		Handler: lib.ResponseLoggerWrapper(policyFiscalcodeFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/v1/policy/{uid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	//TODO: why there are two endpoint with same url? Route:   "/policy/reserved/v1/{policyUid}",
	{
		Route:   "/policy/reserved/v1",
		Handler: lib.ResponseLoggerWrapper(requestApprovalFx),
		Method:  http.MethodPost,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgent,
			lib.UserRoleAgency,
		},
	},
	//	{
	//		Route:   "/policy/reserved/v1",
	//		Handler: lib.ResponseLoggerWrapper(RequestApprovalFx),
	//		Method:  http.MethodPost,
	//		Roles: []string{
	//			lib.UserRoleAdmin,
	//			lib.UserRoleManager,
	//			lib.UserRoleAgent,
	//			lib.UserRoleAgency,
	//		},
	//	},
	{
		Route:   "/v1/policy/emit",
		Handler: lib.ResponseLoggerWrapper(emitFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},

	//	{
	//		Route:   "/v1/policy/emit",
	//		Handler: lib.ResponseLoggerWrapper(EmitFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{lib.UserRoleAll},
	//	},
	//	{
	//		Route:   "/draft/policy/body/emit/",
	//		Handler: lib.ResponseLoggerWrapper(DraftEmitWithPolicyFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{lib.UserRoleAll},
	//	},
	{
		Route:   "/v1/policy/lead",
		Handler: lib.ResponseLoggerWrapper(leadFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	//	{
	//		Route:   "/v1/policy/lead",
	//		Handler: lib.ResponseLoggerWrapper(LeadFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{lib.UserRoleAll},
	//	},
	{
		Route:   "/policy/reserved/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(acceptanceFx),
		Method:  http.MethodPut,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	//	{
	//		Route:   "/policy/reserved/v1/{policyUid}",
	//		Handler: lib.ResponseLoggerWrapper(AcceptanceFx),
	//		Method:  http.MethodPut,
	//		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	//	},
	{
		Route:   "/v1/policy/proposal",
		Handler: lib.ResponseLoggerWrapper(proposalFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	//	{
	//		Route:   "/v1/policy/proposal",
	//		Handler: lib.ResponseLoggerWrapper(ProposalFx),
	//		Method:  http.MethodPost,
	//		Roles:   []string{lib.UserRoleAll},
	//	},
	{
		Route:   "/policy/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(updatePolicyFx),
		Method:  http.MethodPatch,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policy/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(deletePolicyFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleManager},
	},
	{
		Route:   "/attachment/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyAttachmentFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/policy/transactions/v1/{policyUid}",
		Handler: lib.ResponseLoggerWrapper(getPolicyTransactionsFx),
		Method:  http.MethodGet,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
	},
	{
		Route:   "/policy/renew/v1/{uid}",
		Handler: lib.ResponseLoggerWrapper(renew.DeleteRenewPolicyFx),
		Method:  http.MethodDelete,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/transaction/v1/{uid}/receipt",
		Handler: lib.ResponseLoggerWrapper(paymentReceiptFx),
		Method:  http.MethodGet,
		Roles:   []string{lib.UserRoleAdmin, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
	},
	{
		Route:   "/policy/v1/init",
		Handler: lib.ResponseLoggerWrapper(initFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
	{
		Route:   "/portfolio/{type}/{version}",
		Method:  http.MethodGet,
		Handler: lib.ResponseLoggerWrapper(getPortfolioFx),
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency,
			lib.UserRoleAgent},
	},
	{
		Route:   "/policy/v1/contract/upload/{uid}",
		Handler: lib.ResponseLoggerWrapper(uploadPolicyContractFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/policy/v1/duplicate/{uid}",
		Handler: lib.ResponseLoggerWrapper(duplicateFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAdmin},
	},
	{
		Route:   "/v1/policy/validation/address",
		Handler: lib.ResponseLoggerWrapper(addressValidationFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {

	router := lib.GetRouter("broker", brokerRoutes)
	router.ServeHTTP(w, r)
}
