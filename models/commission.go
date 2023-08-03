package models

type CommissionsSetting struct {
	IsFlat        bool        `json:"isFlat" firestore:"isFlat" bigquery:"-"`
	IsByOffer     bool        `json:"isByOffer" firestore:"isByOffer" bigquery:"-"`
	IsByGuarantee bool        `json:"isByGuarantee" firestore:"isByGuarantee" bigquery:"-"`
	Commissions   *Commission `json:"commissions,omitempty" firestore:"commissions,omitempty" bigquery:"-"`
}

type Commission struct {
	NewBusiness float64 `json:"newBusiness" firestore:"newBusiness" bigquery:"-"`
	Renew       float64 `json:"renew" firestore:"renew" bigquery:"-"`
}
