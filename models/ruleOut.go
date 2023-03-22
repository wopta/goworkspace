package models

import "github.com/shopspring/decimal"

type RuleOut struct {
	Coverages  map[string]map[string]*CoverageOut `json:"coverages"`
	OfferPrice map[string]map[string]*Price       `json:"offerPrice"`
}

type Price struct {
	Net      decimal.Decimal `json:"net"`
	Tax      decimal.Decimal `json:"tax"`
	Gross    decimal.Decimal `json:"gross"`
	Delta    decimal.Decimal `json:"delta"`
	Discount decimal.Decimal `json:"discount"`
}

type CoverageOut struct {
	DailyAllowance             string                       `json:"dailyAllowance"`
	Name                       string                       `json:"name"`
	LegalDefence               string                       `json:"legalDefence"`
	Assistance                 string                       `json:"assistance"`
	Group                      string                       `json:"group"`
	CompanyCodec               string                       `json:"companyCodec"`
	CompanyName                string                       `json:"companyName"`
	IsExtension                bool                         `json:"isExtension"`
	IsSellable                 bool                         `json:"isSellable"`
	IsYuor                     bool                         `json:"isYuor"`
	Type                       string                       `json:"type"`
	TypeOfSumInsured           string                       `json:"typeOfSumInsured"`
	Description                string                       `json:"description"`
	Deductible                 string                       `json:"deductible"`
	Tax                        decimal.Decimal              `json:"tax"`
	Taxes                      []TaxOut                     `json:"taxes"`
	SumInsuredLimitOfIndemnity decimal.Decimal              `json:"sumInsuredLimitOfIndemnity"`
	Price                      decimal.Decimal              `json:"price"`
	PriceNett                  decimal.Decimal              `json:"priceNett"`
	PriceGross                 decimal.Decimal              `json:"priceGross"`
	Value                      *CoverageValueOut            `json:"value"`
	Offer                      map[string]*CoverageValueOut `json:"offer"`
	Slug                       string                       `json:"slug"`
	SelfInsurance              string                       `json:"selfInsurance"`
	SelfInsuranceDesc          string                       `json:"selfInsuranceDesc"`
	Config                     *GuaranteValue               `json:"config"`
	IsBase                     bool                         `json:"isBase"`
	IsYour                     bool                         `json:"isYour"`
	IsPremium                  bool                         `json:"isPremium"`
}

type CoverageValueOut struct {
	TypeOfSumInsured           string          `json:"typeOfSumInsured"`
	Deductible                 string          `json:"deductible"`
	DeductibleType             string          `json:"deductibleType"`
	SumInsuredLimitOfIndemnity decimal.Decimal `json:"sumInsuredLimitOfIndemnity"`
	SelfInsurance              string          `json:"selfInsurance"`
	Tax                        decimal.Decimal `json:"tax"`
	Percentage                 decimal.Decimal `json:"percentage"`
	PremiumNet                 decimal.Decimal `json:"premiumNet"`
	PremiumTaxAmount           decimal.Decimal `json:"premiumTaxAmount"`
	PremiumGross               decimal.Decimal `json:"premiumGross"`
}

type TaxOut struct {
	Tax        decimal.Decimal `json:"tax"`
	Percentage decimal.Decimal `json:"percentage"`
}
