package lib

// ROLES

const (
	UserRoleAll         string = "all"
	UserRoleCustomer    string = "customer"
	UserRoleAdmin       string = "admin"
	UserRoleManager     string = "manager"
	UserRoleAgent       string = "agent"
	UserRoleAgency      string = "agency"
	UserRoleAreaManager string = "area-manager"
	UserRoleInternal    string = "internal"
)

func GetAllRoles() []string {
	return []string{UserRoleAll, UserRoleCustomer, UserRoleAdmin, UserRoleManager, UserRoleAgent, UserRoleAgency, UserRoleAreaManager}
}

// CHANNELS

const (
	ECommerceChannel string = "e-commerce"
	MgaChannel       string = "mga"
	NetworkChannel   string = "network"
)

// BIGQUERY DATASETS

const (
	WoptaDataset string = "wopta"
)

// COLLECTIONS

const (
	UserCollection               string = "users"                  // firestore and bigquery
	UsersViewCollection          string = "usersView"              // only for bigquery
	PolicyCollection             string = "policy"                 // firestore and bigquery
	TransactionsCollection       string = "transactions"           // firestore and bigquery
	ClaimsCollection             string = "claims"                 // only for bigquery
	AuditsCollection             string = "audits"                 // only for bigquery
	GuaranteeCollection          string = "guarante"               // only for bigquery
	NetworkNodesCollection       string = "networkNodes"           // firestore and bigquery
	NetworkTransactionCollection string = "networkTransactions"    // only for bigquery
	InvitesCollection            string = "invites"                // only for firestore
	EmergencyNumbersCollection   string = "emergencyNumbers"       // only for firestore
	PoliciesViewCollection       string = "policiesView"           // only for bigquery
	TransactionsViewCollection   string = "transactionsView"       // only for bigquery
	NetworkTreeStructureTable    string = "network-tree-structure" // only for bigquery
	NetworkNodesView             string = "networkNodesView"       // only for bigquery
	MailCollection               string = "mail"                   // only for firestore
	RenewPolicyCollection        string = "renewPolicy"            // firestore and bigquery
	RenewPolicyViewCollection    string = "renewPolicyView"        // bigquery
	RenewTransactionCollection   string = "renewTransactions"      // firestore and bigquery
	MailReportCollection         string = "mailReport"             // only for bigquery
	NodeConsensAuditsCollencion  string = "nodeConsensAudits"      // firestore and bigquery
)
const (
	BaseStorageGoogleUrl = "https://storage.googleapis.com/"
)
