package common

// TODO
// improve me making me dynamic
// the traverse of the workspace should be able to figure it out if a certain
// folder host a module or function or both
const (
	ACCOUNTING   = "accounting"
	BPMN         = "bpmn"
	CALLBACK_OUT = "callback_out"
	DOCUMENT     = "document"
	LIB          = "lib"
	MAIL         = "mail"
	MODELS       = "models"
	NETWORK      = "network"
	PAYMENT      = "payment"
	POLICY       = "policy"
	PRODUCT      = "product"
	QUESTION     = "question"
	QUOTE        = "quote"
	RESERVED     = "reserved"
	SELLABLE     = "sellable"
	TRANSACTION  = "transaction"
	USER         = "user"
	AUTH         = "auth"
	BROKER       = "broker"
	CALLBACK     = "callback"
	CLAIM        = "claim"
	COMPANY_DATA = "companydata"
	ENRICH       = "enrich"
	FORM         = "form"
	MGA          = "mga"
	PARTNERSHIP  = "partnership"
	RULES        = "rules"
	RENEW        = "renew"
)

var MODULES []string = []string{
	ACCOUNTING,
	BPMN,
	CALLBACK_OUT,
	DOCUMENT,
	LIB,
	MAIL,
	MODELS,
	NETWORK,
	PAYMENT,
	POLICY,
	PRODUCT,
	QUESTION,
	QUOTE,
	RESERVED,
	SELLABLE,
	TRANSACTION,
	USER,
}

var FUNCTIONS []string = []string{
	AUTH,
	BROKER,
	CALLBACK,
	CLAIM,
	COMPANY_DATA,
	ENRICH,
	FORM,
	MGA,
	PARTNERSHIP,
	RULES,
	RENEW,
}
