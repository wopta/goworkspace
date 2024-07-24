package companydata

type DataReq struct {
	Day    string `firestore:"-" json:"day,omitempty" bigquery:"-"`
	Upload bool   `firestore:"-" json:"upload,omitempty" bigquery:"-"`
}
type Track struct {
	Columns   []Column  `firestore:"columns,omitempty" json:"columns"`
	Name      string    `firestore:"name,omitempty" json:"name,omitempty"`
	Type      string    `firestore:"type,omitempty" json:"type"`
	Format    string    `firestore:"format,omitempty" json:"format,omitempty"`
	Emit      []Column  `firestore:"emit,omitempty" json:"Emit,omitempty"`
	CsvConfig CsvConfig `firestore:"csvConfig,omitempty" json:"csvConfig,omitempty"`
}

type Column struct {
	Value     string `firestore:"value,omitempty" json:"value"`
	Name      string `firestore:"name,omitempty" json:"name,omitempty"`
	Type      string `firestore:"type,omitempty" json:"type"`
	AssetType string `firestore:"assetType,omitempty" json:"assetType"`
	Format    string `firestore:"format,omitempty" json:"format,omitempty"`
	MapFx     string `firestore:"mapFx,omitempty" json:"mapFx,omitempty"`
	Frame     string `firestore:"frame,omitempty" json:"frame,omitempty"`
}
type CsvConfig struct {
	FileNameFx  string `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Extension string `firestore:"extension,omitempty" json:"extension,omitempty"`
	Separator string `firestore:"separator,omitempty" json:"separator,omitempty"`
}
