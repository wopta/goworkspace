package models

type Coverage struct {
	Type                       string
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	Slug                       string
	SelfInsurance              string
	IsBase                     bool
	IsYuor                     bool
	IsPremium                  bool
	Base                       CoverageValue
	Your                       CoverageValue
	Premium                    CoverageValue
}
type CoverageValue struct {
	TypeOfSumInsured           string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	SelfInsurance              string
}
