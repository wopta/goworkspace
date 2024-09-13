package companydata

type DataReq struct {
	Day    string `firestore:"-" json:"day,omitempty" bigquery:"-"`
	Upload bool   `firestore:"-" json:"upload,omitempty" bigquery:"-"`
	Name   string `firestore:"name,omitempty" json:"name,omitempty"`
	Event  string `firestore:"event,omitempty" json:"event,omitempty"`
}
type Track struct {
	Columns     []Column    `firestore:"columns,omitempty" json:"columns"`
	Name        string      `firestore:"name,omitempty" json:"name,omitempty"`
	Frequency   string      `firestore:"frequency,omitempty" json:"frequency,omitempty"`
	Type        string      `firestore:"type,omitempty" json:"type"`
	UploadType  string      `firestore:"uploadType,omitempty" json:"uploadType,omitempty"`
	Emit        Event       `firestore:"emit,omitempty" json:"Emit,omitempty"`
	FileName    string      `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	FileNameFx  string      `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Payment     Event       `firestore:"payment,omitempty" json:"payment,omitempty"`
	Delete      Event       `firestore:"delete,omitempty" json:"delete,omitempty"`
	CsvConfig   CsvConfig   `firestore:"csvConfig,omitempty" json:"csvConfig,omitempty"`
	ExcelConfig ExcelConfig `firestore:"excelConfig,omitempty" json:"excelConfig,omitempty"`
	FtpConfig   SftpConfig  `firestore:"ftpConfig,omitempty" json:"ftpConfig,omitempty"`
	IsAssetFlat bool        `firestore:"isAssetFlat" json:"isAssetFlat"`
}

type Column struct {
	Value     string            `firestore:"value,omitempty" json:"value"`
	Name      string            `firestore:"name,omitempty" json:"name,omitempty"`
	Type      string            `firestore:"type,omitempty" json:"type"`
	AssetType string            `firestore:"assetType,omitempty" json:"assetType"`
	Format    string            `firestore:"format,omitempty" json:"format,omitempty"`
	MapFx     string            `firestore:"mapFx,omitempty" json:"mapFx,omitempty"`
	MapStatic map[string]string `firestore:"mapStatic,omitempty" json:"mapStatic,omitempty"`
	Frame     string            `firestore:"frame,omitempty" json:"frame,omitempty"`
}
type CsvConfig struct {
	FileNameFx string `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Extension  string `firestore:"extension,omitempty" json:"extension,omitempty"`
	Separator  string `firestore:"separator,omitempty" json:"separator,omitempty"`
}
type ExcelConfig struct {
	FileNameFx string `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Extension  string `firestore:"extension,omitempty" json:"extension,omitempty"`
	SheetName  string `firestore:"sheetName,omitempty" json:"sheetName,omitempty"`
}
type Query struct {
	Field      string      `firestore:"field,omitempty" json:"field,omitempty"`
	Operator   string      `firestore:"operator,omitempty" json:"operator,omitempty"`
	QueryValue interface{} `firestore:"queryValue,omitempty" json:"queryValue,omitempty"`
}
type Event struct {
	Event []Column `firestore:"event,omitempty" json:"event,omitempty"`
	Query []Query  `firestore:"query,omitempty" json:"query,omitempty"`
}
type SftpConfig struct {
	Username     string   `firestore:"username,omitempty" json:"username,omitempty"`
	Password     string   `firestore:"password,omitempty" json:"password,omitempty"`
	PrivateKey   string   `firestore:"privateKey,omitempty" json:"privateKey,omitempty"`
	Server       string   `firestore:"server,omitempty" json:"server,omitempty"`
	KeyExchanges []string `firestore:"keyExchanges,omitempty" json:"keyExchanges,omitempty"`
	KeyPsw       string   `firestore:"keyPsw,omitempty" json:"keyPsw,omitempty"`
	Timeout      int      `firestore:"timeout,omitempty" json:"timeout,omitempty"`
}
