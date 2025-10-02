package companydata

import "time"

// TODO: from to , create a query to use them
type DataReq struct {
	Day    string   `firestore:"-" json:"day,omitempty" bigquery:"-"`
	From   string   `firestore:"-" json:"from,omitempty" bigquery:"-"`
	To     string   `firestore:"-" json:"to,omitempty" bigquery:"-"`
	Upload bool     `firestore:"-" json:"upload,omitempty" bigquery:"-"`
	Name   string   `firestore:"name,omitempty" json:"name,omitempty"`
	Event  []string `firestore:"event,omitempty" json:"event,omitempty"`
}
type Track struct {
	Name        string      `firestore:"name,omitempty" json:"name,omitempty"`
	Frequency   string      `firestore:"frequency,omitempty" json:"frequency,omitempty"`
	Type        string      `firestore:"type,omitempty" json:"type"`
	HasHeader   bool        `firestore:"hasHeader" json:"hasHeader"`
	UploadType  string      `firestore:"uploadType,omitempty" json:"uploadType,omitempty"`
	Emit        Event       `firestore:"emit,omitempty" json:"emit,omitempty"`
	FileName    string      `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	FileNameFx  string      `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Payment     Event       `firestore:"payment,omitempty" json:"payment,omitempty"`
	Delete      Event       `firestore:"delete,omitempty" json:"delete,omitempty"`
	CsvConfig   CsvConfig   `firestore:"csvConfig,omitempty" json:"csvConfig,omitempty"`
	ExcelConfig ExcelConfig `firestore:"excelConfig,omitempty" json:"excelConfig,omitempty"`
	FtpConfig   SftpConfig  `firestore:"ftpConfig,omitempty" json:"ftpConfig,omitempty"`
	IsAssetFlat bool        `firestore:"isAssetFlat" json:"isAssetFlat"`
	SendMail    bool        `firestore:"sendMail" json:"sendMail"`
	MailConfig  MailConfig  `firestore:"mailConfig" json:"mailConfig"`
	now         time.Time
	from        time.Time
	to          time.Time
}

type Column struct {
	Values    []string          `firestore:"values,omitempty" json:"values"`
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
	HasHeader  bool   `firestore:"hasHeader" json:"hasHeader"`
}
type ExcelConfig struct {
	FileNameFx string `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
	Extension  string `firestore:"extension,omitempty" json:"extension,omitempty"`
	SheetName  string `firestore:"sheetName,omitempty" json:"sheetName,omitempty"`
	HasHeader  bool   `firestore:"hasHeader" json:"hasHeader"`
}
type Database struct {
	Dataset    string  `firestore:"dataset,omitempty" json:"dataset,omitempty"`
	Name       string  `firestore:"name,omitempty" json:"name,omitempty"`
	BigQuery   string  `firestore:"bigQuery,omitempty" json:"bigQuery,omitempty"`
	Query      []Query `firestore:"query,omitempty" json:"query,omitempty"`
	RelDataset string  `firestore:"relDataset,omitempty" json:"relDataset,omitempty"`
	RelField   string  `firestore:"relField,omitempty" json:"relField,omitempty"`
}
type Query struct {
	Field      string      `firestore:"field,omitempty" json:"field,omitempty"`
	Operator   string      `firestore:"operator,omitempty" json:"operator,omitempty"`
	QueryValue interface{} `firestore:"queryValue,omitempty" json:"queryValue,omitempty"`
}
type Event struct {
	Event      []Column `firestore:"event,omitempty" json:"event,omitempty"`
	Database   Database `firestore:"database,omitempty" json:"database,omitempty"`
	FileName   string   `firestore:"fileName,omitempty" json:"fileName,omitempty"`
	FileNameFx string   `firestore:"fileNameFx,omitempty" json:"fileNameFx,omitempty"`
}
type SftpConfig struct {
	Username     string   `firestore:"username,omitempty" json:"username,omitempty"`
	Password     string   `firestore:"password,omitempty" json:"password,omitempty"`
	PrivateKey   string   `firestore:"privateKey,omitempty" json:"privateKey,omitempty"`
	Server       string   `firestore:"server,omitempty" json:"server,omitempty"`
	KeyExchanges []string `firestore:"keyExchanges,omitempty" json:"keyExchanges,omitempty"`
	KeyPsw       string   `firestore:"keyPsw,omitempty" json:"keyPsw,omitempty"`
	Timeout      int      `firestore:"timeout,omitempty" json:"timeout,omitempty"`
	Path         string   `firestore:"path,omitempty" json:"path,omitempty"`
}
type MailConfig struct {
	From     string `json:"from"`
	FromName string `json:"fromName"`

	To           []string `json:"to"`
	Message      string   `json:"message"`
	Subject      string   `json:"subject"`
	IsHtml       bool     `json:"isHtml,omitempty"`
	IsAttachment bool     `json:"isAttachment,omitempty"`

	Cc           string `json:"cc,omitempty"`
	Bcc          string `json:"bcc,omitempty"`
	TemplateName string `json:"templateName,omitempty"`
	Title        string `json:"title,omitempty"`
	SubTitle     string `json:"subTitle,omitempty"`
	Content      string `json:"content,omitempty"`
	Link         string `json:"link,omitempty"`
	LinkLabel    string `json:"linkLabel,omitempty"`
	IsLink       bool   `json:"isLink,omitempty"`
	IsApp        bool   `json:"isApp,omitempty"`
}
