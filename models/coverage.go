package models

type Coverage struct {
	Type                       string
	Beneficiary                User
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	Price                      float64
	Value                      *CoverageValue
	Offer                      map[string]*CoverageValue
	Slug                       string
	SelfInsurance              string
	IsBase                     bool
	IsYour                     bool
	IsPremium                  bool
	Base                       *CoverageValue
	Your                       *CoverageValue
	Premium                    *CoverageValue
}
type CoverageValue struct {
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
