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

type CustomerRole string

const (
	UserRolesCustumer CustomerRole = "custumer"
	UserRolesAgent    CustomerRole = "agent"
	UserRolesManager  CustomerRole = "manager"
)

func t() {

}
