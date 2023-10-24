package models

type Roles string

type PolicyStatus string

const (
	PolicyStatusInit               = "Inizialize"
	PolicyStatusInitData           = "InizializeData"
	PolicyStatusInitLead           = "Lead"
	PolicyStatusPartnershipLead    = "PartnershipLead"
	PolicyStatusProspet            = "Prospet"
	PolicyStatusProposal           = "Proposal"
	PolicyStatusContact            = "Contact"
	PolicyStatusToEmit             = "ToEmit"
	PolicyStatusEmited             = "Emited"
	PolicyStatusNeedsApproval      = "NeedsApproval"
	PolicyStatusWaitForApproval    = "WaitForApproval"
	PolicyStatusWaitForApprovalMga = "WaitForApprovalMga"
	PolicyStatusApproved           = "Approved"
	PolicyStatusRejected           = "Rejected"
	PolicyStatusToSign             = "ToSign"
	PolicyStatusSign               = "Signed"
	PolicyStatusPay                = "Paid"
	PolicyStatusToPay              = "ToPay"
	PolicyStatusToRenew            = "Renew"
	PolicyStatusPS                 = "Pay&Sign"
	PolicyStatusCompanyEmit        = "CompanyEmited"
	PolicyStatusDeleted            = "Deleted"
)

func GetWaitForApprovalStatusList() []string {
	return []string{PolicyStatusWaitForApproval, PolicyStatusWaitForApprovalMga}
}

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
	AgentCollection              string = "agents"
	AgencyCollection             string = "agencies"
	UserCollection               string = "users"
	PolicyCollection             string = "policy"
	ProductsCollection           string = "products"
	TransactionsCollection       string = "transactions"
	ClaimsCollection             string = "claims" //only for bigquery
	AuditsCollection             string = "audits"
	GuaranteeCollection          string = "guarante"
	NetworkNodesCollection       string = "networkNodes"
	NetworkTransactionCollection string = "networkTransactions" //only for bigquery
	InvitesCollection            string = "invites"
	EmergencyNumbersCollection   string = "emergencyNumbers"
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
	AgentChannel     string = "agent"  //DEPRECATED: remove this constant once product versioning is completed
	AgencyChannel    string = "agency" //DEPRECATED: remove this constant once product versioning is completed
	MgaChannel       string = "mga"
	NetworkChannel   string = "network"
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
	ManualPaymentProvider  string = "manual"
)

const (
	PayMethodCard       = "creditcard"
	PayMethodTransfer   = "transfer"
	PayMethodSdd        = "sdd"
	PayMethodRemittance = "remittance"
)

func GetAllPaymentMethods() []string {
	return []string{PayMethodCard, PayMethodTransfer, PayMethodSdd}
}

const (
	AgentNetworkNodeType       string = "agent"
	AgencyNetworkNodeType      string = "agency"
	BrokerNetworkNodeType      string = "broker"
	AreaManagerNetworkNodeType string = "area-manager"
	PartnershipNetworkNodeType string = "partnership"
)

const (
	ECommerceFlow     = "e-commerce"
	MgaFlow           = "mga"
	RemittanceMgaFlow = "remittance_mga"
	ProviderMgaFlow   = "provider_mga"
)

const (
	FlowFileFormat = "flows/%s.json"
)
