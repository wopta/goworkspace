package rules

type Coverage struct {
	DailyAllowance             string
	Name                       string
	LegalDefence               string
	Assistance                 string
	Group                      string
	CompanyCodec               string
	CompanyName                string
	IsExtension                bool
	IsSellable                 bool
	IsYuor                     bool
	Type                       string
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	Tax                        float64
	Taxes                      []Tax
	SumInsuredLimitOfIndemnity float64
	Price                      float64
	PriceNett                  float64
	PriceGross                 float64
	Value                      *CoverageValue
	Offer                      map[string]*CoverageValue
	Slug                       string
	SelfInsurance              string
	SelfInsuranceDesc          string
	Config                     *GuaranteValue
	IsBase                     bool
	IsYour                     bool
	IsPremium                  bool
}

type GuaranteValue struct {
	TypeOfSumInsured           string             `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string             `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleValues           GuaranteFieldValue `firestore:"deductibleValues,omitempty" json:"deductibleValues,omitempty"`
	DeductibleType             string             `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	SumInsuredLimitOfIndemnity float64            `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SumInsured                 float64            `json:"sumInsured,omitempty" json:"sumInsured,omitempty"`
	LimitOfIndemnity           float64            `json:"limitOfIndemnity,omitempty" json:"limitOfIndemnity,omitempty"`
	SelfInsurance              string             `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SumInsuredValues           GuaranteFieldValue `firestore:"sumInsuredValues,omitempty" json:"sumInsuredValues,omitempty"`
	DeductibleDesc             string             `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsuranceValues        GuaranteFieldValue `firestore:"selfInsuranceValues,omitempty" json:"selfInsuranceValues,omitempty"`
	SelfInsuranceDesc          string             `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	Duration                   GuaranteFieldValue `firestore:"duration,omitempty" json:"duration,omitempty"`
}

type CoverageValue struct {
	TypeOfSumInsured           string
	Deductible                 string
	DeductibleType             string
	SumInsuredLimitOfIndemnity float64
	SelfInsurance              string
	Tax                        float64
	Percentage                 float64
	PremiumNet                 float64
	PremiumTaxAmount           float64
	PremiumGross               float64
}

type GuaranteFieldValue struct {
	Min    float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max    float64   `firestore:"max,omitempty" json:"max,omitempty"`
	Step   float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values []float64 `firestore:"values,omitempty" json:"values,omitempty"`
}
type Tax struct {
	Tax        float64
	Percentage float64
}
