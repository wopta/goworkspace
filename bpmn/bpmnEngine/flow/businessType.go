package flow

import (
	"net/mail"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/models"
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
	ProviderId    string
}

func (*PaymentInfoBpmn) GetType() string {
	return "paymentInfo"
}

type CallbackConfigBpmn struct {
	callback_out.CallbackConfig
}

func (*CallbackConfigBpmn) GetType() string {
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

type Transaction struct {
	*models.Transaction
}

func (*Transaction) GetType() string {
	return "transaction"
}
