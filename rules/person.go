package rules

/*



INPUT:
{
	"age": 74,
	"work": "select work list da elenco professioni",
	"workType": "dipendente" ,
	"coverageType": "24 ore" ,
	"childrenScool":true,
	"issue1500": 1,
	"riskInLifeIs":1,
	"class":2

}
{
	age:75
	work: "select work list da elenco professioni"
	worktype: dipendente / autonomo / non lavoratore
	coverageType 24 ore / tempo libero / professionale
	childrenScool:true
	issue1500:"si, senza problemi 1 || si, ma dovrei rinunciare a qualcosa 2|| no, non ci riscirei facilmente 3"
	riskInLifeIs:da evitare; 1 da accettare; 2 da gestire 3
	class

}
*/
import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

const (
	base                = "base"
	your                = "your"
	premium             = "premium"
	monthly             = "monthly"
	yearly              = "yearly"
	yearlyPriceMinimum  = 120
	monthlyPriceMinimum = 50
)

func Person(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy    models.Policy
		rulesFile []byte
		e         error
	)
	const (
		rulesFileName = "person.json"
	)

	log.Println("Person")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	quotingInputData := getRulesInputData(policy, e, req)

	rulesFile = getRulesFile(rulesFile, rulesFileName)
	_, coverages := lib.RulesFromJson(rulesFile, initCoverageP(), quotingInputData, []byte(getQuotingData()))
	outJson, out := filterOffers(roundPrices(getOfferPrices(coverages)))
	return outJson, out, nil
}

func getRulesFile(rulesFile []byte, rulesFileName string) []byte {
	switch os.Getenv("env") {
	case "local":
		rulesFile = lib.ErrorByte(os.ReadFile("../function-data/dev/grules/" + rulesFileName))
	case "dev":
		rulesFile = lib.GetFromStorage("function-data", "grules/"+rulesFileName, "")
	case "prod":
		rulesFile = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFileName, "")
	default:

	}
	return rulesFile
}

func getRulesInputData(policy models.Policy, e error, req []byte) []byte {
	policy, e = models.UnmarshalPolicy(req)
	lib.CheckError(e)

	age, e := calculateAge(policy.Contractor.BirthDate)
	lib.CheckError(e)
	policy.QuoteQuestions["age"] = age
	policy.QuoteQuestions["work"] = policy.Contractor.Work
	policy.QuoteQuestions["workType"] = policy.Contractor.WorkType
	policy.QuoteQuestions["class"] = policy.Contractor.RiskClass

	request, e := json.Marshal(policy.QuoteQuestions)
	lib.CheckError(e)
	return request
}

func calculateAge(birthDateIsoString string) (int, error) {
	birthdate, e := time.Parse(time.RFC3339, birthDateIsoString)
	now := time.Now()
	age := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}
	return age, e
}

func initCoverageP() map[string]*CoverageOut {
	var coverages = make(map[string]*CoverageOut)

	coverages["IPI"] = &CoverageOut{
		Slug:                       "Invalidità Permanente Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["D"] = &CoverageOut{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["ITI"] = &CoverageOut{
		Slug:                       "Inabilità Totale Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["DRG"] = &CoverageOut{
		Slug:                       "Diaria Ricovero / Gessatura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["DC"] = &CoverageOut{
		Slug:                       "Diaria Convalescenza Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["RSC"] = &CoverageOut{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["IPM"] = &CoverageOut{
		Slug:                       "Invalidità Permanente Malattia IPM",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["ASS"] = &CoverageOut{
		Slug:                       "Assistenza",
		Deductible:                 "0",
		Tax:                        10,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	coverages["TL"] = &CoverageOut{
		Slug:                       "Tutela Legale",
		Deductible:                 "0",
		Tax:                        21.25,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValueOut{
			base: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			your: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			premium: {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
		},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}

	return coverages
}
func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}

func getOfferPrices(coverages interface{}) *Out {
	offerPrice := make(map[string]map[string]*Price)

	offerPrice[base] = map[string]*Price{
		monthly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		yearly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}
	offerPrice[your] = map[string]*Price{
		monthly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		yearly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}
	offerPrice[premium] = map[string]*Price{
		monthly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		yearly: {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}

	if coveragesStruct, ok := coverages.(map[string]*CoverageOut); ok {
		for _, coverage := range coveragesStruct {
			for offerKey, offerValue := range coverage.Offer {
				offerPrice[offerKey][yearly].Net += offerValue.PremiumNet
				offerPrice[offerKey][yearly].Tax += offerValue.PremiumTaxAmount
				offerPrice[offerKey][yearly].Gross += offerValue.PremiumGross
				offerPrice[offerKey][monthly].Net += offerValue.PremiumNet / 12
				offerPrice[offerKey][monthly].Tax += offerValue.PremiumTaxAmount / 12
				offerPrice[offerKey][monthly].Gross += offerValue.PremiumGross / 12
			}
		}
	}

	out := &Out{
		Coverages:  coverages.(map[string]*CoverageOut),
		OfferPrice: offerPrice,
	}

	return out
}

func roundPrices(out *Out) *Out {
	for offerType, priceStruct := range out.OfferPrice {
		ceilPriceGrossYear := math.Ceil(priceStruct[yearly].Gross)
		priceStruct[yearly].Delta = ceilPriceGrossYear - priceStruct[yearly].Gross
		priceStruct[yearly].Gross = ceilPriceGrossYear
		hasIPIGuarantee := out.Coverages["IPI"].Offer[offerType].PremiumGross > 0
		if hasIPIGuarantee {
			out.Coverages["IPI"].Offer[offerType].PremiumGross += priceStruct[yearly].Delta
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGross += priceStruct[yearly].Delta
		}

		roundPriceGrossMonth := math.Round(priceStruct[monthly].Gross)
		priceStruct[monthly].Delta = roundPriceGrossMonth - priceStruct[monthly].Gross
		priceStruct[monthly].Gross = roundPriceGrossMonth
	}

	return out
}

func filterOffers(out *Out) (string, *Out) {
	toBeDeleted := make([]string, 0)
	for offerType, priceStruct := range out.OfferPrice {
		hasNotOfferMinimumYearlyPrice := priceStruct[yearly].Gross < yearlyPriceMinimum
		hasNotOfferMinimumMonthlyPrice := priceStruct[monthly].Gross < monthlyPriceMinimum
		if hasNotOfferMinimumYearlyPrice || hasNotOfferMinimumMonthlyPrice {
			toBeDeleted = append(toBeDeleted, offerType)
		}
	}

	for _, offerType := range toBeDeleted {
		delete(out.OfferPrice, offerType)
		for _, coverage := range out.Coverages {
			delete(coverage.Offer, offerType)
		}
	}

	outJson, _ := json.Marshal(out)
	return string(outJson), out
}
