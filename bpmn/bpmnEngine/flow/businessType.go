package flow

import (
	"net/mail"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/fabrick"
)

type CallbackInfo struct {
	base.CallbackInfo
}

func (c *CallbackInfo) GetType() string {
	return "callbackInfo"
}

type PaymentInfoBpmn struct {
	Schedule      string
	PaymentMethod string
	fabrick.FabrickCallback
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

type ClientCallback struct {
	callback_out.CallbackClient
}

func (*ClientCallback) GetType() string {
	return "clientCallback"
}

type Policy struct {
	*models.Policy
}

func (p *Policy) GetType() string {
	return "policy"
}

type Product struct {
	*models.Product
}

func (p *Product) GetType() string {
	return "product"
}

type Network struct {
	*models.NetworkNode
}

func (p *Network) GetType() string {
	return "networkNode"
}

type Addresses struct {
	CcAddress, ToAddress, FromAddress mail.Address
}

func (*Addresses) GetType() string {
	return "addresses"
}
