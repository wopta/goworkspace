package models

type Roles string

type PolicyStatus string

const (
	PolicyStatusInit        = "Inizialize"
	PolicyStatusInitData    = "InizializeData"
	PolicyStatusInitLead    = "Lead"
	PolicyStatusProspet     = "Prospet"
	PolicyStatusProposal    = "Proposal"
	PolicyStatusContact     = "Contact"
	PolicyStatusToEmit      = "ToEmit"
	PolicyStatusEmited      = "Emited"
	PolicyStatusToSign      = "ToSign"
	PolicyStatusSign        = "Signed"
	PolicyStatusPay         = "Paid"
	PolicyStatusToPay       = "ToPay"
	PolicyStatusToRenew     = "Renew"
	PolicyStatusPS          = "Pay&Sign"
	PolicyStatusCompanyEmit = "CompanyEmited"
	PolicyStatusDeleted     = "Deleted"
)

const (
	TransactionStatusInit        = "Inizialize"
	TransactionStatusToEmit      = "ToEmit"
	TransactionStatusEmited      = "Emited"
	TransactionStatusToPay       = "ToPay"
	TransactionStatusPay         = "Paid"
	TransactionStatusCompanyEmit = "CompanyEmited"
)

type PaySplit string

const (
	PaySplitMonthly PaySplit = "monthly"
	PaySplitYear    PaySplit = "year"
)

type PayType string

const (
	PayTypeCc       PayType = "credit card"
	PayTypeSdd      PayType = "sdd"
	PayTypeTransfer PayType = "transfer"
)

const (
	PartnershipBeProf string = "beprof"
)

const (
	UserRoleAll      string = "all"
	UserRoleCustomer string = "customer"
	UserRoleAdmin    string = "admin"
	UserRoleManager  string = "manager"
	UserRoleAgent    string = "agent"
	UserRoleAgency   string = "agency"
)
