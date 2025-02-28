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
		Fn:          PolicyFiscalcodeFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.get.policy.fiscalcode",
	},
	{
		Route:       "/v1/policy/{uid}",
		Fn:          GetPolicyFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.get.policy.uid",
	},
	{
		Route:       "/v1/policy/lead",
		Fn:          LeadFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.lead",
	},
	{
		Route:       "/v1/policy/proposal",
		Fn:          ProposalFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.proposal",
	},
	{
		Route:  "/policy/reserved/v1",
		Fn:     RequestApprovalFx,
		Method: http.MethodPost,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgent,
			lib.UserRoleAgency,
		},
		Entitlement: "broker.requestapproval",
	},
	{
		Route:       "/v1/policy/emit",
		Fn:          EmitFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.emit",
	},
	{
		Route:       "/policy/v1/{uid}",
		Fn:          UpdatePolicyFx,
		Method:      http.MethodPatch,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.update.policy",
	},
	{
		Route:       "/policy/v1/{uid}",
		Fn:          DeletePolicyFx,
		Method:      http.MethodDelete,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		Entitlement: "broker.delete.policy",
	},
	{
		Route:       "/attachment/v1/{policyUid}",
		Fn:          GetPolicyAttachmentFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.get.policy.attachment",
	},
	{
		Route:  "/policy/transactions/v1/{policyUid}",
		Fn:     GetPolicyTransactionsFx,
		Method: http.MethodGet,
		Roles: []string{
			lib.UserRoleAdmin,
			lib.UserRoleManager,
			lib.UserRoleAgency,
			lib.UserRoleAgent,
		},
		Entitlement: "broker.get.policy.transactions",
	},
	{
		Route:       "/policy/reserved/v1/{policyUid}",
		Fn:          AcceptanceFx,
		Method:      http.MethodPut,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleManager},
		// TODO: consider splitting different subjects acceptance into different endpoints
		Entitlement: "broker.acceptance",
	},
	{
		Route:       "/policy/renew/v1/{uid}",
		Fn:          renew.DeleteRenewPolicyFx,
		Method:      http.MethodDelete,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "broker.delete.renew",
	},
	{
		Route:       "/transaction/v1/{uid}/receipt",
		Fn:          PaymentReceiptFx,
		Method:      http.MethodGet,
		Roles:       []string{lib.UserRoleAdmin, lib.UserRoleAreaManager, lib.UserRoleAgency, lib.UserRoleAgent},
		Entitlement: "broker.get.transaction.receipt",
	},
	{
		Route:       "/policy/v1/init",
		Fn:          InitFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAll},
		Entitlement: "broker.init",
	},
	{
		Route:  "/portfolio/{type}/{version}",
		Method: http.MethodGet,
		Fn:     GetPortfolioFx,
		Roles: []string{lib.UserRoleAdmin, lib.UserRoleManager, lib.UserRoleAreaManager, lib.UserRoleAgency,
			lib.UserRoleAgent},
		Entitlement: "broker.get.portfolio",
	},
	{
		Route:       "/policy/v1/contract/upload/{uid}",
		Fn:          UploadPolicyContractFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "broker.upload.policy.contract",
	},
	{
		Route:       "/policy/v1/duplicate/{uid}",
		Fn:          DuplicateFx,
		Method:      http.MethodPost,
		Roles:       []string{lib.UserRoleAdmin},
		Entitlement: "broker.duplicate.policy",
	},
}

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := models.GetExtendedRouter("broker", brokerRoutes)
	router.ServeHTTP(w, r)
}
