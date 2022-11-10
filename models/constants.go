package models

type Roles string

const (
	Customer = "Customer"
)

type PolicyStatus string

const (
	Init        = "Inizialize"
	Proposal    = "Proposal"
	Emit        = "Emit"
	Sign        = "Sign"
	Pay         = "Pay"
	PS          = "Pay&Sign"
	CompanyEmit = "CompanyEmit"
)
