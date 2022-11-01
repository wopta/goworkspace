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
		Routes: []lib.Route{{
			Route:   "/risk/person",
			Hendler: Person,
		},
			{
				Route:   "/risk/pmi",
				Hendler: PmiAllrisk,
			},
		},
	}
	route.Router(w, r)

}

type Coverage struct {
	DailyAllowance             string
	LegalDefence               string
	Assistance                 string
	IsYuor                     bool
	Type                       string
	TypeOfSumInsured           string
	Description                string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	Price                      float64
	Value                      *CoverageValue
	Offer                      map[string]*CoverageValue
	Slug                       string
	SelfInsurance              string
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
	PremiumNet                 float64
	PremiumTaxAmount           float64
	PremiumGross               float64
}
