package models

type Roles string

type PolicyStatus string

const (
	PolicyStatusInit               = "Init"
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
	PolicyStatusDraftRenew         = "DraftRenew"
	PolicyStatusRenewed            = "Renewed"
	PolicyStatusUnsolved           = "Unsolved"
	PolicyStatusManualSigned       = "ManualSigned"
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
	TransactionStatusRefunded    = "Refunded"
)

type PaySplit string

const (
	PaySplitMonthly           PaySplit = "monthly"
	PaySplitYear              PaySplit = "year"
	PaySplitYearly            PaySplit = "yearly"
	PaySplitSemestral         PaySplit = "semestral"
	PaySplitSingleInstallment PaySplit = "singleInstallment"
	PaySplitTrimestral        PaySplit = "trimestral"
)

// Map how many rate there are in PaySplit
var PaySplitRateMap = map[PaySplit]int{
	PaySplitMonthly:           12,
	PaySplitYear:              1,
	PaySplitYearly:            1,
	PaySplitSemestral:         2,
	PaySplitTrimestral:        4,
	PaySplitSingleInstallment: 1,
}

// Map how many months there  in a PaySplit
var PaySplitMonthsMap = map[PaySplit]int{
	PaySplitMonthly:           1,
	PaySplitYear:              12,
	PaySplitYearly:            12,
	PaySplitSemestral:         6,
	PaySplitTrimestral:        3,
	PaySplitSingleInstallment: 1,
}

type PayType string

const (
	PayTypeCc       PayType = "credit card"
	PayTypeSdd      PayType = "sdd"
	PayTypeTransfer PayType = "transfer"
)

const (
	PaymentModeSingle    string = "single"
	PaymentModeRecurrent string = "recurrent"
)

func GetAllowedShortTermInstallmentModes() []string {
	return []string{PaymentModeRecurrent}
}

func GetAllowedYearlyModes() []string {
	return []string{PaymentModeRecurrent, PaymentModeSingle}
}

func GetAllowedSingleInstallmentModes() []string {
	return []string{PaymentModeSingle}
}

const (
	PartnershipBeProf      string = "beprof"
	PartnershipFacile      string = "facile"
	PartnershipFpinsurance string = "fpinsurance"
	PartnershipELeads      string = "eleads"
	PartnershipSegugio     string = "segugio"
	PartnershipSwitcho     string = "switcho"
)

// DEPRECATED - use lib version instead
const (
	UserRoleAll         string = "all"
	UserRoleCustomer    string = "customer"
	UserRoleAdmin       string = "admin"
	UserRoleManager     string = "manager"
	UserRoleAgent       string = "agent"
	UserRoleAgency      string = "agency"
	UserRoleAreaManager string = "area-manager"
)

// DEPRECATED - use lib version instead
func GetAllRoles() []string {
	return []string{UserRoleAll, UserRoleCustomer, UserRoleAdmin, UserRoleManager, UserRoleAgent, UserRoleAgency}
}

const (
	TimeDateOnly string = "2006-01-02" //  DEPRECATED use time.DateOnly
)

// DEPRECATED - use lib version instead
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
	PoliciesViewCollection       string = "policiesView"
	TransactionsViewCollection   string = "transactionsView"
	NetworkTreeStructureTable    string = "network-tree-structure"
	NetworkNodesView             string = "networkNodesView"
)

const (
	VehiclePrivateUse string = "private"
)

const (
	BeneficiaryLegalAndWillSuccessors string = "legalAndWillSuccessors"
	BeneficiaryChosenBeneficiary      string = "chosenBeneficiary"
	BeneficiaryLegalEntity            string = "legalEntity"
	BeneficiarySelfLegalEntity        string = "selfLegalEntity"
)

const (
	LifeProduct               string = "life"
	PmiProduct                string = "pmi"
	PersonaProduct            string = "persona"
	GapProduct                string = "gap"
	CommercialCombinedProduct string = "commercial-combined"
	CatNatProduct             string = "cat-nat"
)

const (
	InternalProductType string = "internal"
	ExternalProductType string = "external"
	FormProductType     string = "form"
)

// DEPRECATED - use lib version instead
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
	AxaCompany          string = "axa"
	GlobalCompany       string = "global"
	SogessurCompany     string = "sogessur"
	QBECompany          string = "qbe"
	NetInsuranceCompany string = "net-insurance"
)

var CompanyMap map[string]string = map[string]string{
	AxaCompany:          "AXA FRANCE VIE S.A.",
	SogessurCompany:     "Sogessur SA",
	GlobalCompany:       "Global Assistance",
	NetInsuranceCompany: "Net Insurance",
}

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

func GetAvailableMethods(role string) []string {
	switch role {
	case UserRoleAdmin, UserRoleAreaManager, UserRoleManager:
		return []string{PayMethodCard, PayMethodTransfer, PayMethodSdd}
	case UserRoleAgency, UserRoleAgent:
		return []string{PayMethodRemittance}
	}
	return []string{}
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
	FlowFileFormat                = "flows/%s.json"
	ContractDocumentFormat        = "%s_Contratto_%s.pdf"
	ProposalDocumentFormat        = "%s_Proposta_%d.pdf"
	RvmInstructionsDocumentFormat = "Scheda_Rapporto_Visita_Medica_Proposta_%d.pdf"
	RvmSurveyDocumentFormat       = "Rapporto_Visita_Medica_Proposta_%d.pdf"
	WarrantFormat                 = "warrants/%s.json"
	WarrantsFolder                = "warrants/"
	ProductsFolder                = "products"
	NetInsuranceDocument          = "%s_Net_Insurance"
)

const (
	ProductV1 = "v1"
	ProductV2 = "v2"
)

const (
	ProposalAttachmentName = "Proposta"
	//to use with namirial, to sign
	PreContractAttachmentName = "Precontrattuale"
	ContractAttachmentName    = "Contratto"
	//
	ContractNonDigitalAttachmentName = "Contratto non digitale"
	RvmInstructionsAttachmentName    = "Scheda Rapporto Visita Medica"
	RvmSurveyAttachmentName          = "Rapporto Visita Medica"
)

const (
	WorksForMgaUid = "__wopta__"
	RuiSectionE    = "E"
	RuiSectionA    = "A"
	RuiSectionB    = "B"
)

// contractorFiscalcode = contractor.fiscalCode
// "physical": $.contractor.type = 'individual' OR $.contractor.type = ""
// "enterprise": (($.contractor.type = 'legalEntity'  and .contractorFiscalcode = "”) or  $.contractor.type = 'enterprise'
// "individualCompany": $.contractor.type = 'legalEntity' and .contractorFiscalcode != ”"
const (
	UserIndividual  = "individual"
	UserEterprise   = "enterprise"
	UserLegalEntity = "legalEntity"
)

const (
	TitolareEffettivo = "titolareEffettivo"
)

const (
	DocumentSectionContracts        = "contract"
	DocumentSectionIdentityDocument = "identity-document"
	DocumentSectionReserved         = "reserved"
	DocumentSectionOther            = "other"
)

func GetProponentRuiSections() []string {
	return []string{RuiSectionE}
}

func GetIssuerRuiSections() []string {
	return []string{RuiSectionA, RuiSectionB}
}

const (
	PolicyTypeMultiYear         = "multiYear"
	PolicyTypeYearly            = "yearly"
	PolicyTypeSingleInstallment = "singleInstallment"
)

const (
	QuoteTypeFixed    = "fixed"
	QuoteTypeVariable = "variable"
)

const (
	AssetTypeEnterprise = "enterprise"
	AssetTypeBuilding   = "building"
)

const (
	FAQ = "https://www.wopta.it/it/faq"
)
