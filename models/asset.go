package models

type Asset struct {
	Name         string      `firestore:"name,omitempty" json:"name,omitempty"`
	Address      string      `firestore:"address,omitempty" json:"address,omitempty"`
	Type         string      `firestore:"type,omitempty" json:"type,omitempty"`
	Building     Building    `firestore:"building,omitempty" json:"building,omitempty"`
	Person       User        `firestore:"person,omitempty" json:"person,omitempty"`
	Enterprise   Enterprise  `firestore:"enterprise,omitempty" json:"enterprise,omitempty"`
	IsContractor bool        `firestore:"isContractor,omitempty" json:"isContractor,omitempty"`
	Guarantees   []Guarantee `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
}
type Building struct {
	Name             string `firestore:"name,omitempty" json:"name,omitempty"`
	Address          string `firestore:"address,omitempty" json:"address,omitempty"`
	Type             string `firestore:"type,omitempty" json:"type,omitempty"`
	SquareMeters     int64  `firestore:"squareMeters,omitempty" json:"squareMeters,omitempty"`
	PostalCode       string `firestore:"postalCode" json:"postalCode,omitempty"`
	City             string `firestore:"city" json:"city,omitempty"`
	BuildingType     string `firestore:"buildingType,omitempty" json:"buildingType,omitempty"`
	BuildingMaterial string `firestore:"buildingMaterial,omitempty" json:"buildingMaterial,omitempty"`
	BuildingYear     string `firestore:"buildingYear,omitempty" json:"buildingYear,omitempty"`
	Employer         int64  `firestore:"employer,omitempty" json:"employer,omitempty"`
	IsAllarm         bool   `firestore:"isAllarm,omitempty" json:"isAllarm,omitempty"`
	Floor            string `firestore:"floor,omitempty" json:"floor,omitempty"`
	Costruction      string `firestore:"costruction,omitempty" json:"costruction,omitempty"`
	IsHolder         bool   `firestore:"isHolder,omitempty" json:"isHolder,omitempty"`
}
type Enterprise struct {
	Name       string `firestore:"name,omitempty" json:"name,omitempty"`
	Address    string `firestore:"address,omitempty" json:"address,omitempty"`
	Type       string `firestore:"type,omitempty" json:"type,omitempty"`
	PostalCode string `firestore:"postalCode,omitempty" json:"postalCode,omitempty"`
	City       string `firestore:"city,omitempty" json:"city,omitempty"`
	VatCode    string `firestore:"vatCode,omitempty" json:"vatCode,omitempty"`
	Ateco      string `firestore:"ateco,omitempty" json:"ateco,omitempty"`
	AtecoDesc  string `firestore:"atecoDesc,omitempty" json:"atecoDesc,omitempty"`
	AtecoMacro string `firestore:"atecoMacro,omitempty" json:"atecoMacro,omitempty"`
	AtecoSub   string `firestore:"atecoSub,omitempty" json:"atecoSub,omitempty"`
	Class      string `firestore:"class" json:"class,omitempty"`
	Sector     string `firestore:"sector,omitempty" json:"sector,omitempty"`
	Revenue    string `firestore:"revenue,omitempty" json:"revenue,omitempty"`
	Employer   int64  `firestore:"employer,omitempty" json:"employer,omitempty"`
}
