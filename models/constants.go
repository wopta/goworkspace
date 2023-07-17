package models

type Roles string

type PolicyStatus string

const (
	PolicyStatusInit            = "Inizialize"
	PolicyStatusInitData        = "InizializeData"
	PolicyStatusInitLead        = "Lead"
	PolicyStatusProspet         = "Prospet"
	PolicyStatusProposal        = "Proposal"
	PolicyStatusContact         = "Contact"
	PolicyStatusToEmit          = "ToEmit"
	PolicyStatusEmited          = "Emited"
	PolicyStatusWaitForApproval = "WaitForApproval"
	PolicyStatusApproved        = "Approved"
	PolicyStatusRejected        = "Rejected"
	PolicyStatusToSign          = "ToSign"
	PolicyStatusSign            = "Signed"
	PolicyStatusPay             = "Paid"
	PolicyStatusToPay           = "ToPay"
	PolicyStatusToRenew         = "Renew"
	PolicyStatusPS              = "Pay&Sign"
	PolicyStatusCompanyEmit     = "CompanyEmited"
	PolicyStatusDeleted         = "Deleted"
)

const (
	TransactionStatusInit        = "Inizialize"
	TransactionStatusToEmit      = "ToEmit"
	TransactionStatusEmited      = "Emited"
	TransactionStatusToPay       = "ToPay"
	TransactionStatusPay         = "Paid"
	TransactionStatusCompanyEmit = "CompanyEmited"
	TransactionStatusDeleted     = "Deleted"
)

type PaySplit string

const (
	PaySplitMonthly      PaySplit = "monthly"
	PaySplitYear         PaySplit = "year"
	PaySingleInstallment PaySplit = "singleInstallment"
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

func GetAllRoles() []string {
	return []string{UserRoleAll, UserRoleCustomer, UserRoleAdmin, UserRoleManager, UserRoleAgent, UserRoleAgency}
}

const (
	TimeDateOnly string = "2006-01-02"
)

const (
	AgentCollection    string = "agents"
	AgencyCollection   string = "agencies"
	UserCollection     string = "users"
	PolicyCollection   string = "policy"
	ProductsCollection string = "products"
)

const (
	VehiclePrivateUse string = "private"
)

const (
	LifeProduct    string = "life"
	PmiProduct     string = "pmi"
	PersonaProduct string = "persona"
	GapProduct     string = "gap"
)
