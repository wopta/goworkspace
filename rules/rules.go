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

type Coverage struct {
	DailyAllowance             string
	Name                       string
	LegalDefence               string
	Assistance                 string
	Group                      string
	CompanyCodec               string
	CompanyName                string
	IsExtension                bool
	IsYuor                     bool
	Type                       string
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	Tax                        float64
	Taxes                      []Tax
	SumInsuredLimitOfIndemnity float64
	Price                      float64
	PriceNett                  float64
	PriceGross                 float64
	Value                      *CoverageValue
	Offer                      map[string]*CoverageValue
	Slug                       string
	SelfInsurance              string
	SelfInsuranceDesc          string
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
	Percentage                 float64
	PremiumNet                 float64
	PremiumTaxAmount           float64
	PremiumGross               float64
}
type Tax struct {
	Tax        float64
	Percentage float64
}
