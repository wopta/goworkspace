package models

type Guarantee struct {
	Type                       string
	Beneficiary                User
	TypeOfSumInsured           string
	Description                string
	Value                      *CoverageValue
	Offer                      map[string]*CoverageValue
	Slug                       string
	IsBase                     bool
	IsYour                     bool
	IsPremium                  bool
	Base                       *GuaranteeValue
	Your                       *GuaranteeValue
	Premium                    *GuaranteeValue
	Name                       *string `json:"name,omitempty"`
	Deductible                 *string `json:"deductible,omitempty"`
	SelfInsurance              *string `json:"selfInsurance,omitempty"`
	SumInsuredLimitOfIndemnity *int64  `json:"sumInsuredLimitOfIndemnity,omitempty"`
	Price                      *int64  `json:"price,omitempty"`
	PriceNett                  *int64  `json:"priceNett,omitempty"`
}
type GuaranteeValue struct {
	TypeOfSumInsured           string
	Deductible                 string
	DeductibleType             string
	SumInsuredLimitOfIndemnity float64
	SelfInsurance              string
	Tax                        float64
	PremiumNet                 float64
	PremiumTaxAmount           float64
	PremiumGross               float64
}
