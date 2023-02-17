package models

type Guarantee struct {
	Type                       string                     `firestore:"type,omitempty" json:"type,omitempty"`
	Beneficiary                *User                      `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"`
	TypeOfSumInsured           string                     `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Description                string                     `firestore:"description,omitempty" json:"description,omitempty"`
	ContractDetail             string                     `firestore:"contractDetail,omitempty" json:"contractDetail,omitempty"`
	CompanyCodec               string                     `firestore:"companyCodec,omitempty" json:"companyCodec,omitempty"`
	CompanyName                string                     `firestore:"companyName,omitempty" json:"companyName,omitempty"`
	Group                      string                     `firestore:"group,omitempty" json:"group,omitempty"`
	Value                      *GuaranteeValue            `firestore:"value,omitempty" json:"value,omitempty"`
	ExtraValue                 string                     `firestore:"extraValue,omitempty" json:"extraValue,omitempty"`
	ValueDesc                  string                     `firestore:"valueDesc,omitempty" json:"valueDesc,omitempty"`
	Offer                      *map[string]*CoverageValue `firestore:"offer,omitempty" json:"offer,omitempty"`
	Slug                       string                     `firestore:"slug,omitempty" json:"slug,omitempty"`
	IsExtension                bool                       `firestore:"isExtension,omitempty" json:"isExtension,omitempty"`
	IsBase                     bool                       `firestore:"isBase,omitempty" json:"isBase,omitempty"`
	IsYour                     bool                       `firestore:"isYour,omitempty" json:"isYour,omitempty"`
	IsPremium                  bool                       `firestore:"isPremium,omitempty" json:"isPremium,omitempty"`
	Base                       *GuaranteeValue            `firestore:"base,omitempty" json:"base,omitempty"`
	Discount                   float64                    `json:"discount,omitempty" json:"discount,omitempty"`
	Your                       *GuaranteeValue            `firestore:"your,omitempty" json:"your,omitempty"`
	Premium                    *GuaranteeValue            `firestore:"premium,omitempty" json:"premium,omitempty"`
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
type GuaranteeValue struct {
	TypeOfSumInsured           string  `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string  `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleType             string  `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SelfInsurance              string  `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	Tax                        float64 `json:"tax,omitempty" json:"tax,omitempty"`
	Price                      float64 `firestore:"price,omitempty" json:"price,omitempty"`
	PriceNett                  float64 `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross                 float64 `firestore:"priceGross,omitempty" json:"priceGross   ,omitempty"`
}
type Tax struct {
	Tax        float64 `firestore:"tax,omitempty" json:"tax,omitempty"`
	Percentage float64 `firestore:"percentage,omitempty" json:"percentage,omitempty"`
}
