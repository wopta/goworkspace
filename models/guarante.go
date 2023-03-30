package models

type Guarante struct {
	DailyAllowance             string                    `firestore:"dailyAllowance" json:"dailyAllowance,omitempty"`
	OrderAsset                 int                       `firestore:"orderAsset,omitempty" json:"orderAsset,omitempty"`
	LegalDefence               string                    `firestore:"legalDefence" json:"legalDefence,omitempty"`
	Assistance                 string                    `firestore:"assistance " json:"assistance ,omitempty"`
	Type                       string                    `firestore:"type,omitempty" json:"type,omitempty"`
	Beneficiary                *User                     `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"`
	Beneficiaries              *[]User                   `firestore:"beneficiaries,omitempty" json:"beneficiaries,omitempty"`
	TypeOfSumInsured           string                    `firestore:"typeOfSumInsured" json:"typeOfSumInsured,omitempty"`
	Description                string                    `firestore:"description,omitempty" json:"description,omitempty"`
	ContractDetail             string                    `firestore:"contractDetail,omitempty" json:"contractDetail,omitempty"`
	CompanyCodec               string                    `firestore:"companyCodec,omitempty" json:"companyCodec,omitempty"`
	CompanyName                string                    `firestore:"companyName,omitempty" json:"companyName,omitempty"`
	Group                      string                    `firestore:"group,omitempty" json:"group,omitempty"`
	Value                      *GuaranteValue            `firestore:"value,omitempty" json:"value,omitempty"`
	Config                     *GuaranteValue            `firestore:"config,omitempty" json:"config,omitempty"`
	ExtraValue                 string                    `firestore:"extraValue,omitempty" json:"extraValue,omitempty"`
	ValueDesc                  string                    `firestore:"valueDesc,omitempty" json:"valueDesc,omitempty"`
	Offer                      map[string]*GuaranteValue `firestore:"offer,omitempty" json:"offer,omitempty"`
	Slug                       string                    `firestore:"slug" json:"slug,omitempty"`
	IsMandatory                bool                      `firestore:"isMandatory" json:"isMandatory"`
	IsExtension                bool                      `firestore:"isExtension" json:"isExtension"`
	Discount                   float64                   `json:"discount,omitempty" json:"discount,omitempty"`
	Name                       string                    `firestore:"name,omitempty" json:"name,omitempty"`
	SumInsuredLimitOfIndemnity float64                   `json:"sumInsuredLimitOfIndemnity" json:"sumInsuredLimitOfIndemnity,omitempty"`
	Deductible                 string                    `firestore:"deductible" json:"deductible,omitempty"`
	DeductibleDesc             string                    `firestore:"deductibleDesc" json:"deductibleDesc,omitempty"`
	SelfInsurance              string                    `firestore:"selfInsurance" json:"selfInsurance,omitempty"`
	SelfInsuranceDesc          string                    `firestore:"selfInsuranceDesc" json:"selfInsuranceDesc,omitempty"`
	Tax                        float64                   `json:"tax,omitempty" json:"tax,omitempty"`
	Taxs                       []Tax                     `json:"taxs,omitempty" json:"taxs,omitempty"`
	Price                      float64                   `firestore:"price,omitempty" json:"price,omitempty"`
	PriceNett                  float64                   `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross                 float64                   `firestore:"priceGross,omitempty" json:"priceGross,omitempty"`
	IsSellable                 bool                      `firestore:"isSellable" json:"isSellable"`
	IsConfigurable             bool                      `firestore:"isConfigurable" json:"isConfigurable"`
}
type GuaranteValue struct {
	TypeOfSumInsured           string              `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string              `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleValues           GuaranteFieldValue  `firestore:"deductibleValues,omitempty" json:"deductibleValues,omitempty"`
	DeductibleType             string              `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	SumInsuredLimitOfIndemnity float64             `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SumInsured                 float64             `json:"sumInsured,omitempty" json:"sumInsured,omitempty"`
	LimitOfIndemnity           float64             `json:"limitOfIndemnity,omitempty" json:"limitOfIndemnity,omitempty"`
	SelfInsurance              string              `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SumInsuredValues           GuaranteFieldValue  `firestore:"sumInsuredValues,omitempty" json:"sumInsuredValues,omitempty"`
	DeductibleDesc             string              `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsuranceValues        GuaranteFieldValue  `firestore:"selfInsuranceValues,omitempty" json:"selfInsuranceValues,omitempty"`
	SelfInsuranceDesc          string              `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	Duration                   Duration            `firestore:"duration,omitempty" json:"duration,omitempty"`
	DurationValues             *DurationFieldValue `firestore:"durationValues,omitempty" json:"durationValues,omitempty"`
	Tax                        float64             `firestore:"tax" json:"tax"`
	Percentage                 float64             `firestore:"percentage" json:"percentage"`
	PremiumNetYearly           float64             `firestore:"premiumNetYearly" json:"premiumNetYearly"`
	PremiumTaxAmountYearly     float64             `firestore:"premiumTaxAmountYearly" json:"premiumTaxAmountYearly"`
	PremiumGrossYearly         float64             `firestore:"premiumGrossYearly" json:"premiumGrossYearly"`
	PremiumNetMonthly          float64             `firestore:"premiumNetMonthly" json:"premiumNetMonthly"`
	PremiumTaxAmountMonthly    float64             `firestore:"premiumTaxAmountMonthly" json:"premiumTaxAmountMonthly"`
	PremiumGrossMonthly        float64             `firestore:"premiumGrossMonthly" json:"premiumGrossMonthly"`
	MinimumGrossMonthly        float64             `firestore:"minimumGrossMonthly,omitempty" json:"minimumGrossMonthly,omitempty"`
	MinimumGrossYearly         float64             `firestore:"minimumGrossYearly,omitempty" json:"minimumGrossYearly,omitempty"`
}
type GuaranteFieldValue struct {
	Min    float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max    float64   `firestore:"max,omitempty" json:"max,omitempty"`
	Step   float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values []float64 `firestore:"values,omitempty" json:"values,omitempty"`
}

type DurationFieldValue struct {
	Min  int `firestore:"min,omitempty" json:"min,omitempty"`
	Max  int `firestore:"max,omitempty" json:"max,omitempty"`
	Step int `firestore:"step,omitempty" json:"step,omitempty"`
}

type Tax struct {
	Tax        float64 `firestore:"tax,omitempty" json:"tax,omitempty"`
	Percentage float64 `firestore:"percentage,omitempty" json:"percentage,omitempty"`
}
type Duration struct {
	Year int `firestore:"year,omitempty" json:"year,omitempty"`
}
