package flow

import (
	"net/http"

	"github.com/wopta/goworkspace/callback"
	"github.com/wopta/goworkspace/models"
)

type CallbackInfo struct {
	Request     *http.Request
	RequestBody []byte
	Response    *http.Response
	Error       error
	Action      string
}

func (c *CallbackInfo) GetType() string {
	return "callbackInfo"
}

type PaymentInfoBpmn struct {
	Schedule      string
	PaymentMethod string
	callback.FabrickCallback
}

func (*PaymentInfoBpmn) GetType() string {
	return "paymentInfo"
}

type CallbackConfig struct {
	Proposal        bool `json:"proposal"`
	RequestApproval bool `json:"requestApproval"`
	Emit            bool `json:"emit"`
	Pay             bool `json:"pay"`
	Sign            bool `json:"sign"`

	//need to integrate inside channel_flow first
	//need to define AcceptanceFx
	Approved bool `json:"approved"`
	Rejected bool `json:"rejected"`
}

func (*CallbackConfig) GetType() string {
	return "callbackConfig"
}

type PolicyDraft struct {
	*models.Policy
}

func (p *PolicyDraft) GetType() string {
	return "policy"
}

type ProductDraft struct {
	*models.Product
}

func (p *ProductDraft) GetType() string {
	return "product"
}

type NetworkDraft struct {
	*models.NetworkNode
}

func (p *NetworkDraft) GetType() string {
	return "networkNode"
}
