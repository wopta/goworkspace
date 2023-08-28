package companydata

type DataReq struct {
	Day    string `firestore:"-" json:"day,omitempty" bigquery:"-"`
	Upload bool   `firestore:"-" json:"upload,omitempty" bigquery:"-"`
}
