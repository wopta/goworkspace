package models

type Guarantee struct {
	Type                       string                    `firestore:"type,omitempty" json:"type,omitempty"`
	Beneficiary                User                      `firestore:"beneficiary,omitempty" json:"beneficiary,omitempty"`
	TypeOfSumInsured           string                    `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Description                string                    `firestore:"description,omitempty" json:"description,omitempty"`
	Value                      GuaranteeValue            `firestore:"value,omitempty" json:"value,omitempty"`
	Offer                      map[string]*CoverageValue `firestore:"offer,omitempty" json:"offer,omitempty"`
	Slug                       string                    `firestore:"slug,omitempty" json:"slug,omitempty"`
	IsBase                     bool                      `firestore:"isBase,omitempty" json:"isBase,omitempty"`
	IsYour                     bool                      `firestore:"isYour,omitempty" json:"isYour,omitempty"`
	IsPremium                  bool                      `firestore:"isPremium,omitempty" json:"isPremium,omitempty"`
	Base                       GuaranteeValue            `firestore:"base,omitempty" json:"base,omitempty"`
	Your                       GuaranteeValue            `firestore:"your,omitempty" json:"your,omitempty"`
	Premium                    GuaranteeValue            `firestore:"premium,omitempty" json:"premium,omitempty"`
	Name                       string                    `firestore:"name,omitempty" json:"name,omitempty"`
	SumInsuredLimitOfIndemnity float64                   `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	Deductible                 string                    `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	SelfInsurance              string                    `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	Tax                        float64                   `json:"tax,omitempty" json:"tax,omitempty"`
	Price                      int64                     `firestore:"price,omitempty" json:"price,omitempty"`
	PriceNett                  int64                     `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross                 int64                     `firestore:"priceGross,omitempty" json:"priceGross   ,omitempty"`
}
type GuaranteeValue struct {
	TypeOfSumInsured           string  `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string  `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleType             string  `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SelfInsurance              string  `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	Tax                        float64 `json:"tax,omitempty" json:"tax,omitempty"`
	Price                      int64   `firestore:"price,omitempty" json:"price,omitempty"`
	PriceNett                  int64   `firestore:"priceNett,omitempty" json:"priceNett,omitempty"`
	PriceGross                 int64   `firestore:"priceGross,omitempty" json:"priceGross   ,omitempty"`
}
