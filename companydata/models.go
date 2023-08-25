package companydata

type DataReq struct {
	Day string `firestore:"-" json:"day,omitempty" bigquery:"-"`
}
