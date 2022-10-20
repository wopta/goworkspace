package models

type CoverageJson struct {
	Type                       string
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	Price                      float64
	Value                      *CoverageJsonValue
	Offer                      map[string]*CoverageJsonValue
	Slug                       string
	SelfInsurance              string
	IsBase                     bool
	IsYuor                     bool
	IsPremium                  bool
	Base                       *CoverageJsonValue
	Your                       *CoverageJsonValue
	Premium                    *CoverageJsonValue
}
type CoverageJsonValue struct {
	TypeOfSumInsured           string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	SelfInsurance              string
	Tax                        float64
	PremiumNet                 float64
	PremiumTaxAmount           float64
	PremiumGross               float64
}
