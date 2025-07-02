package models

import (
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
)

type Guarante struct {
	Status                     string                     `firestore:"status,omitempty" json:"status,omitempty" bigquery:"status"`
	PolicyUid                  string                     `firestore:"-" json:"-"  bigquery:"policyUid"`
	DailyAllowance             string                     `firestore:"dailyAllowance" json:"dailyAllowance,omitempty"  bigquery:"-"`
	OrderAsset                 int                        `firestore:"orderAsset,omitempty" json:"orderAsset,omitempty"  bigquery:"-"`
	DependsOr                  []string                   `firestore:"dependsOr,omitempty" json:"dependsOr,omitempty"  bigquery:"-"`
	Depends                    []string                   `firestore:"depends,omitempty" json:"depends,omitempty"  bigquery:"-"`
	Childrens                  []string                   `firestore:"childrens ,omitempty" json:"childrens,omitempty"  bigquery:"-"`
	SelectedExclude            []string                   `firestore:"selectedExclude,omitempty" json:"selectedExclude,omitempty"  bigquery:"-"`
	LegalDefence               string                     `firestore:"legalDefence" json:"legalDefence,omitempty"  bigquery:"legalDefence"`
	Assistance                 string                     `firestore:"assistance" json:"assistance,omitempty"  bigquery:"-"`
	Type                       string                     `firestore:"type,omitempty" json:"type,omitempty"  bigquery:"type"`
	Beneficiary                *User                      `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"  bigquery:"-"`
	BeneficiaryReference       *User                      `firestore:"beneficiaryReference,omitempty" json:"beneficiaryReference,omitempty"  bigquery:"-"`
	Beneficiaries              *[]Beneficiary             `firestore:"beneficiaries,omitempty" json:"beneficiaries,omitempty"  bigquery:"-"`
	BeneficiaryOptions         map[string]string          `json:"beneficiaryOptions,omitempty" firestore:"beneficiaryOptions,omitempty" bigquery:"-"`
	TypeOfSumInsured           string                     `firestore:"typeOfSumInsured" json:"typeOfSumInsured,omitempty"  bigquery:"typeOfSumInsured"`
	Description                string                     `firestore:"description,omitempty" json:"description,omitempty"  bigquery:"-"`
	ContractDetail             string                     `firestore:"contractDetail,omitempty" json:"contractDetail,omitempty"  bigquery:"-"`
	CompanyCodec               string                     `firestore:"companyCodec,omitempty" json:"companyCodec,omitempty"  bigquery:"-"`
	CompanyName                string                     `firestore:"companyName,omitempty" json:"companyName,omitempty"  bigquery:"companyName"`
	Group                      string                     `firestore:"group,omitempty" json:"group,omitempty"  bigquery:"group"`
	Section                    string                     `firestore:"section,omitempty" json:"section,omitempty"  bigquery:"-"`
	Value                      *GuaranteValue             `firestore:"value,omitempty" json:"value,omitempty"  bigquery:"-"`
	Config                     *GuaranteConfig            `firestore:"config,omitempty" json:"config,omitempty"  bigquery:"-"`
	ExtraValue                 string                     `firestore:"extraValue,omitempty" json:"extraValue,omitempty"  bigquery:"-"`
	ValueDesc                  string                     `firestore:"valueDesc,omitempty" json:"valueDesc,omitempty"  bigquery:"-"`
	Offer                      map[string]*GuaranteValue  `firestore:"offer,omitempty" json:"offer,omitempty"  bigquery:"-"`
	Slug                       string                     `firestore:"slug" json:"slug,omitempty"  bigquery:"slug"`
	IsMandatory                bool                       `firestore:"isMandatory" json:"isMandatory"  bigquery:"-"`
	IsExtension                bool                       `firestore:"isExtension" json:"isExtension"  bigquery:"-"`
	Discount                   float64                    `firestore:"discount" json:"discount"  bigquery:"-"`
	Name                       string                     `firestore:"name,omitempty" json:"name,omitempty"  bigquery:"name"`
	SumInsuredLimitOfIndemnity float64                    `firestore:"sumInsuredLimitOfIndemnity" json:"sumInsuredLimitOfIndemnity,omitempty"  bigquery:"sumInsuredLimitOfIndemnity"`
	Deductible                 string                     `firestore:"deductible" json:"deductible,omitempty"  bigquery:"deductible"`
	DeductibleDesc             string                     `firestore:"deductibleDesc" json:"deductibleDesc,omitempty"  bigquery:"-"`
	SelfInsurance              string                     `firestore:"selfInsurance" json:"selfInsurance,omitempty"  bigquery:"selfInsurance"`
	SelfInsuranceDesc          string                     `firestore:"selfInsuranceDesc" json:"selfInsuranceDesc,omitempty"  bigquery:"-"`
	Tax                        float64                    `firestore:"tax,omitempty" json:"tax,omitempty"  bigquery:"tax"`
	Taxs                       []Tax                      `firestore:"taxs,omitempty" json:"taxs,omitempty"  bigquery:"-"`
	Price                      float64                    `firestore:"price,omitempty" json:"price,omitempty"  bigquery:"-"`
	PriceNett                  float64                    `firestore:"priceNett,omitempty" json:"priceNett,omitempty"  bigquery:"priceNett"`
	PriceGross                 float64                    `firestore:"priceGross,omitempty" json:"priceGross,omitempty"  bigquery:"priceGross"`
	IsSellable                 bool                       `firestore:"isSellable" json:"isSellable"  bigquery:"-"`
	IsConfigurable             bool                       `firestore:"isConfigurable" json:"isConfigurable"  bigquery:"-"`
	Subtitle                   string                     `firestore:"subtitle" json:"subtitle"  bigquery:"-"`
	Reserved                   bool                       `json:"reserved" firestore:"reserved" bigquery:"-"`
	ReservedConfig             map[string]*ReservedLimits `json:"reservedConfig,omitempty" firestore:"reservedConfig,omitempty" bigquery:"-"`
	Order                      int64                      `json:"order,omitempty" firestore:"order,omitempty" bigquery:"-"`
	Commissions                *Commissions               `json:"commissions,omitempty" firestore:"commissions,omitempty" bigquery:"-"`
	IsSelected                 bool                       `json:"isSelected" firestore:"isSelected" bigquery:"-"`
	IsDeleted                  bool                       `json:"isDeleted" firestore:"isDeleted" bigquery:"-"`
	SynchronizeSlug            string                     `json:"synchronizeSlug" firestore:"synchronizeSlug" bigquery:"-"`
}

type GuaranteValue struct {
	StartDate                  *time.Time          `firestore:"startDate,omitempty" json:"startDate,omitempty" bigquery:"-"`
	StartDateString            string              `firestore:"startDateString,omitempty" json:"startDateString,omitempty" bigquery:"-"`
	TypeOfSumInsured           string              `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string              `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	RetroactiveDate            *time.Time          `firestore:"retroactiveDate ,omitempty" json:"retroactiveDate,omitempty" bigquery:"-"`
	Discount                   float64             `firestore:"discount" json:"discount"  bigquery:"-"`
	RetroactiveUsaCanDate      *time.Time          `firestore:"retroactiveUsaCanDate,omitempty" json:"retroactiveUsaCanDate,omitempty" bigquery:"-"`
	DeductibleValues           GuaranteFieldValue  `firestore:"deductibleValues,omitempty" json:"deductibleValues,omitempty"`
	DeductibleType             string              `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	DeductibleUnit             string              `firestore:"deductibleUnit,omitempty" json:"deductibleUnit,omitempty"`
	SumInsuredLimitOfIndemnity float64             `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SumInsured                 float64             `json:"sumInsured,omitempty" json:"sumInsured,omitempty"`
	LimitOfIndemnity           float64             `json:"limitOfIndemnity,omitempty" json:"limitOfIndemnity,omitempty"`
	SelfInsurance              string              `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SumInsuredValues           GuaranteFieldValue  `firestore:"sumInsuredValues,omitempty" json:"sumInsuredValues,omitempty"`
	DeductibleDesc             string              `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsuranceValues        GuaranteFieldValue  `firestore:"selfInsuranceValues,omitempty" json:"selfInsuranceValues,omitempty"`
	SelfInsuranceDesc          string              `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	Duration                   *Duration           `firestore:"duration,omitempty" json:"duration,omitempty"`
	DurationValues             *DurationFieldValue `firestore:"durationValues,omitempty" json:"durationValues,omitempty"`
	Tax                        float64             `firestore:"tax" json:"tax"`
	Percentage                 float64             `firestore:"percentage" json:"percentage"`
	PremiumNetYearly           float64             `firestore:"premiumNetYearly,omitempty" json:"premiumNetYearly"`
	PremiumTaxAmountYearly     float64             `firestore:"premiumTaxAmountYearly,omitempty" json:"premiumTaxAmountYearly"`
	PremiumGrossYearly         float64             `firestore:"premiumGrossYearly,omitempty" json:"premiumGrossYearly"`
	PremiumNetMonthly          float64             `firestore:"premiumNetMonthly,omitempty" json:"premiumNetMonthly,omitempty"`
	PremiumTaxAmountMonthly    float64             `firestore:"premiumTaxAmountMonthly,omitempty" json:"premiumTaxAmountMonthly"`
	PremiumGrossMonthly        float64             `firestore:"premiumGrossMonthly,omitempty" json:"premiumGrossMonthly"`
	MinimumGrossMonthly        float64             `firestore:"minimumGrossMonthly,omitempty" json:"minimumGrossMonthly,omitempty"`
	MinimumGrossYearly         float64             `firestore:"minimumGrossYearly,omitempty" json:"minimumGrossYearly,omitempty"`
}
type GuaranteConfig struct {
	TypeOfSumInsured                       string               `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                             string               `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleValues                       GuaranteFieldValue   `firestore:"deductibleValues,omitempty" json:"deductibleValues,omitempty"`
	DeductibleType                         string               `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	DeductibleUnit                         string               `firestore:"deductibleUnit,omitempty" json:"deductibleUnit,omitempty"`
	SumInsuredLimitOfIndemnity             float64              `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	LimitOfIndemnity                       float64              `json:"limitOfIndemnity,omitempty" json:"limitOfIndemnity,omitempty"`
	SelfInsurance                          string               `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SumInsuredValues                       GuaranteFieldValue   `firestore:"sumInsuredValues,omitempty" json:"sumInsuredValues,omitempty"`
	DeductibleDesc                         string               `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsuranceValues                    GuaranteFieldValue   `firestore:"selfInsuranceValues,omitempty" json:"selfInsuranceValues,omitempty"`
	SelfInsuranceDesc                      string               `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	DiscountConfig                         *DiscountConfig      `firestore:"discountConfig,omitempty" json:"discountConfig,omitempty"  bigquery:"-"`
	Duration                               *Duration            `firestore:"duration,omitempty" json:"duration,omitempty"`
	DurationValues                         *DurationFieldValue  `firestore:"durationValues,omitempty" json:"durationValues,omitempty"`
	Tax                                    float64              `firestore:"tax" json:"tax"`
	Percentage                             float64              `firestore:"percentage" json:"percentage"`
	MinimumGrossMonthly                    float64              `firestore:"minimumGrossMonthly,omitempty" json:"minimumGrossMonthly,omitempty"`
	MinimumGrossYearly                     float64              `firestore:"minimumGrossYearly,omitempty" json:"minimumGrossYearly,omitempty"`
	DurationSlider                         *GuaranteFieldConfig `firestore:"durationSlider,omitempty" json:"durationSlider,omitempty"`
	SumInsuredLimitOfIndemnityDropdown     *GuaranteFieldConfig `firestore:"sumInsuredLimitOfIndemnityDropdown,omitempty" json:"sumInsuredLimitOfIndemnityDropdown,omitempty"`
	StartDatePicker                        *GuaranteFieldConfig `firestore:"startDatePicker,omitempty" json:"startDatePicker,omitempty"`
	SumInsuredLimitOfIndemnityTextField    *GuaranteFieldConfig `firestore:"sumInsuredLimitOfIndemnityTextField,omitempty" json:"sumInsuredLimitOfIndemnityTextField,omitempty"`
	LimitOfIndemnityTextField              *GuaranteFieldConfig `firestore:"limitOfIndemnityTextField,omitempty" json:"limitOfIndemnityTextField,omitempty"`
	SumInsuredTextField                    *GuaranteFieldConfig `firestore:"sumInsuredTextField,omitempty" json:"sumInsuredTextField,omitempty"`
	DurationStartDateStringSliderTextField *GuaranteFieldConfig `firestore:"durationStartDateStringSliderTextField,omitempty" json:"durationStartDateStringSliderTextField,omitempty"`
	RetroactiveDatePicker                  *GuaranteFieldConfig `firestore:"retroactiveDatePicker,omitempty" json:"retroactiveDatePicker,omitempty"`
	RetroactiveUsaCanDatePicker            *GuaranteFieldConfig `firestore:"retroactiveUsaCanDatePicker,omitempty" json:"retroactiveUsaCanDatePicker,omitempty"`
}
type GuaranteFieldValue struct {
	Min    float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max    float64   `firestore:"max,omitempty" json:"max,omitempty"`
	Step   float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values []float64 `firestore:"values,omitempty" json:"values,omitempty"`
	Label  string    `firestore:"label,omitempty" json:"label,omitempty"`
}
type GuaranteFieldConfig struct {
	Min            float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max            float64   `firestore:"max,omitempty" json:"max,omitempty"`
	DurationMin    *Duration `firestore:"durationMin,omitempty" json:"durationMin,omitempty"`
	DurationMax    *Duration `firestore:"durationMax,omitempty" json:"durationMax,omitempty"`
	Step           float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values         []float64 `firestore:"values,omitempty" json:"values,omitempty"`
	Label          string    `firestore:"label,omitempty" json:"label,omitempty"`
	DefaultValue   float64   `firestore:"defaultValue,omitempty" json:"defaultValue,omitempty"`
	Description    string    `firestore:"description,omitempty" json:"description,omitempty"`
	TextFieldLabel string    `firestore:"textFieldLabel,omitempty" json:"textFieldLabel,omitempty"`
}
type DiscountConfig struct {
	Min          float64 `firestore:"min,omitempty" json:"min,omitempty"`
	Max          float64 `firestore:"max,omitempty" json:"max,omitempty"`
	Step         float64 `firestore:"step,omitempty" json:"step,omitempty"`
	Label        string  `firestore:"label,omitempty" json:"label,omitempty"`
	DefaultValue float64 `firestore:"defaultValue" json:"defaultValue"`
	Description  string  `firestore:"description,omitempty" json:"description,omitempty"`
}
type DurationFieldValue struct {
	Min   int    `firestore:"min,omitempty" json:"min,omitempty"`
	Max   int    `firestore:"max,omitempty" json:"max,omitempty"`
	Step  int    `firestore:"step,omitempty" json:"step,omitempty"`
	Label string `firestore:"label,omitempty" json:"label,omitempty"`
}

type Tax struct {
	Tax        float64 `firestore:"tax,omitempty" json:"tax,omitempty"`
	Percentage float64 `firestore:"percentage,omitempty" json:"percentage,omitempty"`
}

type Duration struct {
	Year    int `firestore:"year,omitempty" json:"year,omitempty"`
	Day     int `firestore:"day,omitempty" json:"day,omitempty"`
	Monthly int `firestore:"monthly,omitempty" json:"monthly,omitempty"`
}

type ReservedLimits struct {
	MaxAge int64 `json:"maxAge" firestore:"maxAge" bigquery:"-"`
	MinAge int64 `json:"minAge" firestore:"minAge" bigquery:"-"`
}

func SetGuaranteBigquery(policy Policy, status string, origin string) error {
	var e error
	for _, asset := range policy.Assets {
		for _, g := range asset.Guarantees {
			g.Status = status
			g.PolicyUid = policy.Uid
			e = lib.InsertRowsBigQuery(WoptaDataset, origin, g)
			if e != nil {
				log.Error(e)
			}
		}

	}
	return e
}
