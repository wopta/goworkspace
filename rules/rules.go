package rules

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT Rules")

	functions.HTTP("Rules", Rules)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/risk/person",
				Handler: Person,
				Method:  "POST",
			},
			{
				Route:   "/survey/person",
				Handler: PersonSurvey,
				Method:  "POST",
			},
			{
				Route:   "/risk/pmi",
				Handler: PmiAllrisk,
				Method:  "POST",
			},
		},
	}
	route.Router(w, r)

}

type Price struct {
	Net      float64 `json:"net"`
	Tax      float64 `json:"tax"`
	Gross    float64 `json:"gross"`
	Delta    float64 `json:"delta"`
	Discount float64 `json:"discount"`
}

type Out struct {
	Coverages  map[string]*Coverage         `json:"coverages"`
	OfferPrice map[string]map[string]*Price `json:"offerPrice"`
}

type Coverage struct {
	DailyAllowance             string                    `json:"dailyAllowance"`
	Name                       string                    `json:"name"`
	LegalDefence               string                    `json:"legalDefence"`
	Assistance                 string                    `json:"assistance"`
	Group                      string                    `json:"group"`
	CompanyCodec               string                    `json:"companyCodec"`
	CompanyName                string                    `json:"companyName"`
	IsExtension                bool                      `json:"isExtension"`
	IsSellable                 bool                      `json:"isSellable"`
	IsYuor                     bool                      `json:"isYuor"`
	Type                       string                    `json:"type"`
	TypeOfSumInsured           string                    `json:"typeOfSumInsured"`
	Description                string                    `json:"description"`
	Deductible                 string                    `json:"deductible"`
	Tax                        float64                   `json:"tax"`
	Taxes                      []Tax                     `json:"taxes"`
	SumInsuredLimitOfIndemnity float64                   `json:"sumInsuredLimitOfIndemnity"`
	Price                      float64                   `json:"price"`
	PriceNett                  float64                   `json:"priceNett"`
	PriceGross                 float64                   `json:"priceGross"`
	Value                      *CoverageValue            `json:"value"`
	Offer                      map[string]*CoverageValue `json:"offer"`
	Slug                       string                    `json:"slug"`
	SelfInsurance              string                    `json:"selfInsurance"`
	SelfInsuranceDesc          string                    `json:"selfInsuranceDesc"`
	Config                     *GuaranteValue            `json:"config"`
	IsBase                     bool                      `json:"isBase"`
	IsYour                     bool                      `json:"isYour"`
	IsPremium                  bool                      `json:"isPremium"`
}

type GuaranteValue struct {
	TypeOfSumInsured           string             `firestore:"typeOfSumInsured,omitempty" json:"typeOfSumInsured,omitempty"`
	Deductible                 string             `firestore:"deductible,omitempty" json:"deductible,omitempty"`
	DeductibleValues           GuaranteFieldValue `firestore:"deductibleValues,omitempty" json:"deductibleValues,omitempty"`
	DeductibleType             string             `firestore:"deductibleType,omitempty" json:"deductibleType,omitempty"`
	SumInsuredLimitOfIndemnity float64            `json:"sumInsuredLimitOfIndemnity,omitempty" json:"sumInsuredLimitOfIndemnity,omitempty"`
	SumInsured                 float64            `json:"sumInsured,omitempty" json:"sumInsured,omitempty"`
	LimitOfIndemnity           float64            `json:"limitOfIndemnity,omitempty" json:"limitOfIndemnity,omitempty"`
	SelfInsurance              string             `firestore:"selfInsurance,omitempty" json:"selfInsurance,omitempty"`
	SumInsuredValues           GuaranteFieldValue `firestore:"sumInsuredValues,omitempty" json:"sumInsuredValues,omitempty"`
	DeductibleDesc             string             `firestore:"deductibleDesc,omitempty" json:"deductibleDesc,omitempty"`
	SelfInsuranceValues        GuaranteFieldValue `firestore:"selfInsuranceValues,omitempty" json:"selfInsuranceValues,omitempty"`
	SelfInsuranceDesc          string             `firestore:"selfInsuranceDesc,omitempty" json:"selfInsuranceDesc,omitempty"`
	Duration                   GuaranteFieldValue `firestore:"duration,omitempty" json:"duration,omitempty"`
}
type CoverageValue struct {
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

type GuaranteFieldValue struct {
	Min    float64   `firestore:"min,omitempty" json:"min,omitempty"`
	Max    float64   `firestore:"max,omitempty" json:"max,omitempty"`
	Step   float64   `firestore:"step,omitempty" json:"step,omitempty"`
	Values []float64 `firestore:"values,omitempty" json:"values,omitempty"`
}
type Tax struct {
	Tax        float64
	Percentage float64
}
