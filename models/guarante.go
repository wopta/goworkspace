package models

type Guarante struct {
	PolicyUid                  string                    `firestore:"-" json:"-"  bigquery:"policyUid"`
	DailyAllowance             string                    `firestore:"dailyAllowance" json:"dailyAllowance,omitempty"  bigquery:"-"`
	OrderAsset                 int                       `firestore:"orderAsset,omitempty" json:"orderAsset,omitempty"  bigquery:"-"`
	LegalDefence               string                    `firestore:"legalDefence" json:"legalDefence,omitempty"  bigquery:"legalDefence"`
	Assistance                 string                    `firestore:"assistance " json:"assistance,omitempty"  bigquery:"-"`
	Type                       string                    `firestore:"type,omitempty" json:"type,omitempty"  bigquery:"-"`
	Beneficiary                *User                     `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"  bigquery:"-"`
	BeneficiaryReferance       *User                     `firestore:"beneficiaryReferance,omitempty" json:"beneficiaryReferance,omitempty"  bigquery:"-"`
	Beneficiaries              *[]User                   `firestore:"beneficiaries,omitempty" json:"beneficiaries,omitempty"  bigquery:"-"`
	TypeOfSumInsured           string                    `firestore:"typeOfSumInsured" json:"typeOfSumInsured,omitempty"  bigquery:"-"`
	Description                string                    `firestore:"description,omitempty" json:"description,omitempty"  bigquery:"-"`
	ContractDetail             string                    `firestore:"contractDetail,omitempty" json:"contractDetail,omitempty"  bigquery:"-"`
	CompanyCodec               string                    `firestore:"companyCodec,omitempty" json:"companyCodec,omitempty"  bigquery:"-"`
	CompanyName                string                    `firestore:"companyName,omitempty" json:"companyName,omitempty"  bigquery:"companyName"`
	Group                      string                    `firestore:"group,omitempty" json:"group,omitempty"  bigquery:"group"`
	Value                      *GuaranteValue            `firestore:"value,omitempty" json:"value,omitempty"  bigquery:"-"`
	Config                     *GuaranteValue            `firestore:"config,omitempty" json:"config,omitempty"  bigquery:"-"`
	ExtraValue                 string                    `firestore:"extraValue,omitempty" json:"extraValue,omitempty"  bigquery:"-"`
	ValueDesc                  string                    `firestore:"valueDesc,omitempty" json:"valueDesc,omitempty"  bigquery:"-"`
	Offer                      map[string]*GuaranteValue `firestore:"offer,omitempty" json:"offer,omitempty"  bigquery:"-"`
	Slug                       string                    `firestore:"slug" json:"slug,omitempty"  bigquery:"slug"`
	IsMandatory                bool                      `firestore:"isMandatory" json:"isMandatory"  bigquery:"-"`
	IsExtension                bool                      `firestore:"isExtension" json:"isExtension"  bigquery:"-"`
	Discount                   float64                   `firestore:"discount,omitempty" json:"discount,omitempty"  bigquery:"-"`
	Name                       string                    `firestore:"name,omitempty" json:"name,omitempty"  bigquery:"name"`
	SumInsuredLimitOfIndemnity float64                   `firestore:"sumInsuredLimitOfIndemnity" json:"sumInsuredLimitOfIndemnity,omitempty"  bigquery:"sumInsuredLimitOfIndemnity"`
	Deductible                 string                    `firestore:"deductible" json:"deductible,omitempty"  bigquery:"deductible"`
	DeductibleDesc             string                    `firestore:"deductibleDesc" json:"deductibleDesc,omitempty"  bigquery:"-"`
	SelfInsurance              string                    `firestore:"selfInsurance" json:"selfInsurance,omitempty"  bigquery:"selfInsurance"`
	SelfInsuranceDesc          string                    `firestore:"selfInsuranceDesc" json:"selfInsuranceDesc,omitempty"  bigquery:"-"`
	Tax                        float64                   `firestore:"tax,omitempty" json:"tax,omitempty"  bigquery:"tax"`
	Taxs                       []Tax                     `firestore:"taxs,omitempty" json:"taxs,omitempty"  bigquery:"-"`
	Price                      float64                   `firestore:"price,omitempty" json:"price,omitempty"  bigquery:"-"`
	PriceNett                  float64                   `firestore:"priceNett,omitempty" json:"priceNett,omitempty"  bigquery:"priceNett"`
	PriceGross                 float64                   `firestore:"priceGross,omitempty" json:"priceGross,omitempty"  bigquery:"priceGross"`
	IsSellable                 bool                      `firestore:"isSellable" json:"isSellable"  bigquery:"-"`
	IsConfigurable             bool                      `firestore:"isConfigurable" json:"isConfigurable"  bigquery:"-"`
	Subtitle                   string                    `firestore:"subtitle" json:"subtitle"  bigquery:"-"`
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
	PremiumNetYearly           float64             `firestore:"premiumNetYearly,omitempty" json:"premiumNetYearly"`
	PremiumTaxAmountYearly     float64             `firestore:"premiumTaxAmountYearly,omitempty" json:"premiumTaxAmountYearly,omitempty"`
	PremiumGrossYearly         float64             `firestore:"premiumGrossYearly,omitempty" json:"premiumGrossYearly,omitempty"`
	PremiumNetMonthly          float64             `firestore:"premiumNetMonthly,omitempty" json:"premiumNetMonthly,omitempty"`
	PremiumTaxAmountMonthly    float64             `firestore:"premiumTaxAmountMonthly,omitempty" json:"premiumTaxAmountMonthly,omitempty"`
	PremiumGrossMonthly        float64             `firestore:"premiumGrossMonthly,omitempty" json:"premiumGrossMonthly,omitempty"`
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
