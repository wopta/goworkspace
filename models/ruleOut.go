package models

type RuleOut struct {
	Coverages  map[string]*CoverageOut      `json:"coverages"`
	OfferPrice map[string]map[string]*Price `json:"offerPrice"`
}

type Price struct {
	Net      float64 `json:"net"`
	Tax      float64 `json:"tax"`
	Gross    float64 `json:"gross"`
	Delta    float64 `json:"delta"`
	Discount float64 `json:"discount"`
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
	Tax                        float64                      `json:"tax"`
	Taxes                      []TaxOut                     `json:"taxes"`
	SumInsuredLimitOfIndemnity float64                      `json:"sumInsuredLimitOfIndemnity"`
	Price                      float64                      `json:"price"`
	PriceNett                  float64                      `json:"priceNett"`
	PriceGross                 float64                      `json:"priceGross"`
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
	TypeOfSumInsured           string  `json:"typeOfSumInsured"`
	Deductible                 string  `json:"deductible"`
	DeductibleType             string  `json:"deductibleType"`
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
	SelfInsurance              string  `json:"selfInsurance"`
	Tax                        float64 `json:"tax"`
	Percentage                 float64 `json:"percentage"`
	PremiumNet                 float64 `json:"premiumNet"`
	PremiumTaxAmount           float64 `json:"premiumTaxAmount"`
	PremiumGross               float64 `json:"premiumGross"`
}

type TaxOut struct {
	Tax        float64 `json:"tax"`
	Percentage float64 `json:"percentage"`
}
