package models

type Asset struct {
	Name                 string      `firestore:"name,omitempty" json:"name,omitempty"`
	Type                 string      `firestore:"type,omitempty" json:"type,omitempty"`
	Uuid                 string      `firestore:"uuid,omitempty" json:"uuid,omitempty"`
	Building             *Building   `firestore:"building,omitempty" json:"building,omitempty"`
	Person               *User       `firestore:"person,omitempty" json:"person,omitempty"`
	Enterprise           *Enterprise `firestore:"enterprise,omitempty" json:"enterprise,omitempty"`
	IsContractor         bool        `firestore:"isContractor,omitempty" json:"isContractor,omitempty"`
	Guarantees           []Guarante  `firestore:"guarantees,omitempty" json:"guarantees,omitempty"`
	Vehicle              *Vehicle    `firestore:"vehicle,omitempty" json:"vehicle,omitempty"`
	TypeUseOfAsset       *string     `firestore:"typeUseOfAsset,omitempty" json:"typeUseOfAsset,omitempty"`
	RangeNumbesrOfFloors *string     `firestore:"rangeNumbersOfFloors,omitempty" json:"rangeNumbersOfFloors,omitempty"`
	GroundFloor          *string     `firestore:"theGroundFloor,omitempty" json:"theGroundFloor,omitempty"`
}

type Building struct {
	Name             string   `firestore:"name,omitempty" json:"name,omitempty"`
	Type             string   `firestore:"type,omitempty" json:"type,omitempty"`
	StreetNumber     string   `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty"`
	CityCode         string   `firestore:"cityCode,omitempty" json:"cityCode,omitempty"`
	PostalCode       string   `firestore:"postalCode" json:"postalCode,omitempty"`
	City             string   `firestore:"city" json:"city,omitempty"`
	Locality         string   `firestore:"locality" json:"locality,omitempty"`
	Location         Location `firestore:"location" json:"location,omitempty"`
	Address          string   `firestore:"address,omitempty"         json:"address,omitempty"         bigquery:"-"`
	BuildingType     string   `firestore:"buildingType,omitempty" json:"buildingType,omitempty"`
	BuildingMaterial string   `firestore:"buildingMaterial,omitempty" json:"buildingMaterial,omitempty"`
	BuildingYear     string   `firestore:"buildingYear,omitempty" json:"buildingYear,omitempty"`
	Employer         int64    `firestore:"employer,omitempty" json:"employer,omitempty"`
	HasAlarm         bool     `firestore:"hasAlarm" json:"hasAlarm"`
	Floor            string   `firestore:"floor,omitempty" json:"floor,omitempty"`
	Ateco            string   `firestore:"ateco,omitempty" json:"ateco,omitempty"`
	AtecoDesc        string   `firestore:"atecoDesc,omitempty" json:"atecoDesc,omitempty"`
	AtecoMacro       string   `firestore:"atecoMacro,omitempty" json:"atecoMacro,omitempty"`
	AtecoSub         string   `firestore:"atecoSub,omitempty" json:"atecoSub,omitempty"`
	Costruction      string   `firestore:"costruction,omitempty" json:"costruction,omitempty"`
	IsHolder         bool     `firestore:"isHolder,omitempty" json:"isHolder,omitempty"`
	NaicsDetail      string   `firestore:"naicsDetail,omitempty" json:"naicsDetail,omitempty"`
	NaicsCategory    string   `firestore:"naicsCategory,omitempty" json:"naicsCategory,omitempty"`
	Naics            string   `firestore:"naics,omitempty" json:"naics,omitempty"`
	IsNaicsSellable  bool     `firestore:"isNaicsSellable,omitempty" json:"isNaicsSellable,omitempty"`
	HasSandwichPanel bool     `firestore:"hasSandwichPanel" json:"hasSandwichPanel"`
	HasSprinkler     bool     `firestore:"hasSprinkler" json:"hasSprinkler"`
	BuildingAddress  *Address `firestore:"buildingAddress,omitempty" json:"buildingAddress,omitempty"`
	UseType          string   `firestore:"useType,omitempty" json:"useType,omitempty"`
	LowestFloor      string   `firestore:"lowestFloor,omitempty" json:"lowestFloor,omitempty"`
}

type Enterprise struct {
	Name                      string   `firestore:"name,omitempty" json:"name,omitempty"`
	Address                   string   `firestore:"address,omitempty"         json:"address,omitempty"         bigquery:"-"`
	StreetNumber              string   `firestore:"streetNumber,omitempty" json:"streetNumber,omitempty"`
	Location                  Location `firestore:"location" json:"location,omitempty"`
	Type                      string   `firestore:"type,omitempty" json:"type,omitempty"`
	PostalCode                string   `firestore:"postalCode,omitempty" json:"postalCode,omitempty"`
	City                      string   `firestore:"city,omitempty" json:"city,omitempty"`
	CityCode                  string   `firestore:"cityCode,omitempty" json:"cityCode,omitempty"`
	Locality                  string   `firestore:"locality" json:"locality,omitempty"`
	VatCode                   string   `firestore:"vatCode,omitempty" json:"vatCode,omitempty"`
	FiscalCode                string   `firestore:"fiscalCode" json:"fiscalCode,omitempty"`
	Ateco                     string   `firestore:"ateco,omitempty" json:"ateco,omitempty"`
	AtecoDesc                 string   `firestore:"atecoDesc,omitempty" json:"atecoDesc,omitempty"`
	AtecoMacro                string   `firestore:"atecoMacro,omitempty" json:"atecoMacro,omitempty"`
	AtecoSub                  string   `firestore:"atecoSub,omitempty" json:"atecoSub,omitempty"`
	Class                     string   `firestore:"class" json:"class,omitempty"`
	Sector                    string   `firestore:"sector,omitempty" json:"sector,omitempty"`
	Revenue                   float64  `firestore:"revenue,omitempty" json:"revenue,omitempty"`
	Employer                  int64    `firestore:"employer,omitempty" json:"employer,omitempty"`
	WorkEmployersRemuneration float64  `firestore:"workEmployersRemuneration,omitempty" json:"workEmployersRemuneration,omitempty"`
	TotalBilled               float64  `firestore:"totalBilled,omitempty" json:"totalBilled,omitempty"`
	NorthAmericanMarket       float64  `firestore:"northAmericanMarket,omitempty" json:"northAmericanMarket,omitempty"`
	YearOfEstablishment       int      `firestore:"yearOfEstablishment,omitempty" json:"yearOfEstablishment,omitempty"`
	EnterpriseAddress         *Address `firestore:"enterpriseAddress,omitempty" json:"enterpriseAddress,omitempty"`
}

type Ateco struct {
	Name       string `firestore:"name,omitempty" json:"name,omitempty"`
	Type       string `firestore:"type,omitempty" json:"type,omitempty"`
	VatCode    string `firestore:"vatCode,omitempty" json:"vatCode,omitempty"`
	Ateco      string `firestore:"ateco,omitempty" json:"ateco,omitempty"`
	AtecoDesc  string `firestore:"atecoDesc,omitempty" json:"atecoDesc,omitempty"`
	AtecoMacro string `firestore:"atecoMacro,omitempty" json:"atecoMacro,omitempty"`
	AtecoSub   string `firestore:"atecoSub,omitempty" json:"atecoSub,omitempty"`
	Class      string `firestore:"class" json:"class,omitempty"`
	Sector     string `firestore:"sector,omitempty" json:"sector,omitempty"`
}

type Location struct {
	Lat float64 `firestore:"lat,omitempty" json:"lat,omitempty"`
	Lng float64 `firestore:"lng,omitempty" json:"lng,omitempty"`
}

func (a *Asset) Normalize() {
	if a.Person != nil {
		a.Person.Normalize()
	}
	if a.Vehicle != nil {
		a.Vehicle.Normalize()
	}
}
