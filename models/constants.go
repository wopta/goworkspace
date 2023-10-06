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
	PolicyStatusNeedsApproval   = "NeedsApproval"
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
	PaySplitMonthly           PaySplit = "monthly"
	PaySplitYear              PaySplit = "year"
	PaySplitYearly            PaySplit = "yearly"
	PaySplitSemestral         PaySplit = "semestral"
	PaySplitSingleInstallment PaySplit = "singleInstallment"
)

type PayType string

const (
	PayTypeCc       PayType = "credit card"
	PayTypeSdd      PayType = "sdd"
	PayTypeTransfer PayType = "transfer"
)

const (
	PartnershipBeProf string = "beprof"
	PartnershipFacile string = "facile"
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
	AgentCollection        string = "agents"
	AgencyCollection       string = "agencies"
	UserCollection         string = "users"
	PolicyCollection       string = "policy"
	ProductsCollection     string = "products"
	TransactionsCollection string = "transactions"
	ClaimsCollection       string = "claims" //only for bigquery
	AuditsCollection       string = "audits"
	GuaranteeCollection    string = "guarante"
	NetworkNodesCollection string = "networkNodes"
)

const (
	VehiclePrivateUse string = "private"
)

const (
	BeneficiaryLegalAndWillSuccessors string = "legalAndWillSuccessors"
	BeneficiaryChosenBeneficiary      string = "chosenBeneficiary"
)

const (
	LifeProduct    string = "life"
	PmiProduct     string = "pmi"
	PersonaProduct string = "persona"
	GapProduct     string = "gap"
)

const (
	ECommerceChannel string = "e-commerce"
	AgentChannel     string = "agent"
	AgencyChannel    string = "agency"
	MgaChannel       string = "mga"
)

const (
	WoptaDataset string = "wopta"
)

const (
	AxaCompany      string = "axa"
	GlobalCompany   string = "global"
	SogessurCompany string = "sogessur"
)

const (
	FabrickPaymentProvider string = "fabrick"
	ManualPaymentProvider  string = "payment"
)

const (
	AgentProducerType       string = "agent"
	AgencyProducerType      string = "agency"
	PartnershipProducerType string = "partnership"
)
