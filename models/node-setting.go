package models

type NodeSetting struct {
	EmitFlow     []Process `json:"emitFlow" firestore:"emitFlow" bigquery:"-"`
	LeadFlow     []Process `json:"leadFlow" firestore:"leadFlow" bigquery:"-`
	ProposalFlow []Process `json:"proposalFlow" firestore:"proposalFlow" bigquery:"-"`
	PayFlow      []Process `json:"payFlow" firestore:"payFlow" bigquery:"-"`
	SignFlow     []Process `json:"signFlow" firestore:"signFlow" bigquery:"-"`
	MailProposal string    `firestore:"mailProposal" json:"mailProposal" bigquery:"-"`
	MailEmitted  string    `firestore:"mailEmitted" json:"mailEmittedx" bigquery:"-"`
	MailPay      string    `firestore:"mailPay" json:"mailPay" bigquery:"-"`
	MailSign     string    `firestore:"mailSign" json:"mailSign" bigquery:"-"`
}

type Process struct {
	Id              int                    `firestore:"id" json:"id" bigquery:"-"`
	LayerId         int                    `firestore:"layer" json:"layer" bigquery:"-"`
	Name            string                 `firestore:"name" json:"name" bigquery:"-"`
	Shape           string                 `firestore:"shape" json:"shape" bigquery:"-"`
	Type            string                 `firestore:"type" json:"type" bigquery:"-"`
	Status          string                 `firestore:"status" json:"status" bigquery:"-"`
	Decision        string                 `firestore:"decision" json:"decision" bigquery:"-"`
	DecisionData    map[string]interface{} `firestore:"decisionData" json:"decisionData" bigquery:"-"`
	Data            interface{}            `firestore:"data" json:"data" bigquery:"-"`
	X               float64                `firestore:"x" json:"x" bigquery:"-"`
	Y               float64                `firestore:"y" json:"y" bigquery:"-"`
	InProcess       []int                  `firestore:"inProcess" json:"inProcess" bigquery:"-"`
	OutProcess      []int                  `firestore:"outProcess" json:"outProcess" bigquery:"-"`
	OutTrueProcess  []int                  `firestore:"outTrueProcess" json:"outTrueProcess" bigquery:"-"`
	OutFalseProcess []int                  `firestore:"outFalseProcess" json:"outFalseProcess" bigquery:"-"`
	IsCompleted     bool                   `firestore:"isCompleted" json:"isCompleted" bigquery:"-"`
	IsFailed        bool                   `firestore:"isFailed" json:"isFailed" bigquery:"-"`
	Error           string                 `firestore:"error" json:"error" bigquery:"-"`
}
