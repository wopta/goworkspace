package models

type Agent struct {
	User
	ManagerUid string    `json:"managerUid,omitempty" firestore:"managerUid,omitempty" bigquery:"-"`
	AgencyUid  string    `json:"agencyUid" firestore:"agencyUid" bigquery:"-"`
	Agents     []string  `json:"agents" firestore:"agents" bigquery:"-"`
	Portfolio  []string  `json:"portfolio" firestore:"portfolio" bigquery:"-"` // will contain users UIDs
	IsActive   bool      `json:"isActive" firestore:"isActive" bigquery:"-"`
	Products   []Product `json:"products" firestore:"products" bigquery:"-"`
	Policies   []string  `json:"policies" firestore:"policies" bigquery:"-"` // will contain policies UIDs
	RuiCode    string    `json:"ruiCode" firestore:"ruiCode" bigquery:"-"`
}
