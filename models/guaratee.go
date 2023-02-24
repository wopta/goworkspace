package models

type Guarantee struct {
	DailyAllowance             string                     `firestore:"dailyAllowance,omitempty" json:"dailyAllowance,omitempty"`
	OrderAsset                 int                        `firestore:"orderAsset,omitempty" json:"orderAsset,omitempty"`
	LegalDefence               string                     `firestore:"legalDefence,omitempty" json:"legalDefence,omitempty"`
	Assistance                 string                     `firestore:"assistance ,omitempty" json:"assistance ,omitempty"`
	Type                       string                     `firestore:"type,omitempty" json:"type,omitempty"`
	Beneficiary                *User                      `firestore:"beneficiaries,omitempty" json:"beneficiaries,omitempty"`
	Beneficiaries              *[]User                    `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"`
	TypeOfSumInsured           string                     `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Description                string                     `firestore:"description,omitempty" json:"description,omitempty"`
	ContractDetail             string                     `firestore:"contractDetail,omitempty" json:"contractDetail,omitempty"`
	CompanyCodec               string                     `firestore:"companyCodec,omitempty" json:"companyCodec,omitempty"`
	CompanyName                string                     `firestore:"companyName,omitempty" json:"companyName,omitempty"`
	Group                      string                     `firestore:"group,omitempty" json:"group,omitempty"`
	Value                      *GuaranteValue             `firestore:"value,omitempty" json:"value,omitempty"`
	Config                     *GuaranteValue             `firestore:"config,omitempty" json:"config,omitempty"`
	ExtraValue                 string                     `firestore:"extraValue,omitempty" json:"extraValue,omitempty"`
	ValueDesc                  string                     `firestore:"valueDesc,omitempty" json:"valueDesc,omitempty"`
	Offer                      *map[string]*GuaranteValue `firestore:"offer,omitempty" json:"offer,omitempty"`
	Slug                       string                     `firestore:"slug,omitempty" json:"slug,omitempty"`
	IsMandatory                bool                       `firestore:"isMandatory ,omitempty" json:"isMandatory ,omitempty"`
	IsExtension                bool                       `firestore:"isExtension,omitempty" json:"isExtension,omitempty"`
	IsBase                     bool                       `firestore:"isBase,omitempty" json:"isBase,omitempty"`
	IsYour                     bool                       `firestore:"isYour,omitempty" json:"isYour,omitempty"`
	IsPremium                  bool                       `firestore:"isPremium,omitempty" json:"isPremium,omitempty"`
	Base                       *GuaranteValue             `firestore:"base,omitempty" json:"base,omitempty"`
	Discount                   float64                    `json:"discount,omitempty" json:"discount,omitempty"`
	Your                       *GuaranteValue             `firestore:"your,omitempty" json:"your,omitempty"`
	Premium                    *GuaranteValue             `firestore:"premium,omitempty" json:"premium,omitempty"`
	Name                       string                     `firestore:"name,omitempty" json:"name,omitempty"`
	SumInsuredLimitOfIndemnity float64                    `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	Deductible                 string                     `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleDesc             string                     `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsurance              string                     `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SelfInsuranceDesc          string                     `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	Tax                        float64                    `json:"tax,omitempty" json:"tax,omitempty"`
	Taxs                       []Tax                      `json:"taxs,omitempty" json:"taxs,omitempty"`
	Price                      float64                    `firestore:"price,omitempty" json:"price,omitempty"`
	PriceNett                  float64                    `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross                 float64                    `firestore:"priceGross,omitempty" json:"priceGross,omitempty"`
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
type GuaranteFieldValue struct {
	Min    float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max    float64   `firestore:"max,omitempty" json:"max,omitempty"`
	Step   float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values []float64 `firestore:"values,omitempty" json:"values,omitempty"`
}
type Tax struct {
	Tax        float64 `firestore:"tax,omitempty" json:"tax,omitempty"`
	Percentage float64 `firestore:"percentage,omitempty" json:"percentage,omitempty"`
}
type Duration struct {
	Year int `firestore:"year,omitempty" json:"year,omitempty"`
}
