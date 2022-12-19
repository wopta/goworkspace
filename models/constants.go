package models

type Roles string

const (
	customer Roles = "Customer"
)

type PolicyStatus string

const (
	Init        PolicyStatus = "Inizialize"
	Proposal                 = "Proposal"
	ToEmit                   = "ToEmit"
	Emited                   = "Emited"
	Sign                     = "Sign"
	Pay                      = "Pay"
	PS                       = "Pay&Sign"
	CompanyEmit              = "CompanyEmit"
)

type PaySplit string

const (
	Montly PaySplit = "montly"
	Year   PaySplit = "year"
)

type PayType string

const (
	Cc       PayType = "credit card"
	Sdd      PayType = "sdd"
	Transfer PayType = "transfer"
)

type CustomerRole string

const (
	Custumer CustomerRole = "custumer"
	Agent    CustomerRole = "agent"
	Manager  CustomerRole = "manager"
)

func t() {

}
