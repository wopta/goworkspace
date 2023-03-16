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
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
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
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	quotingInputData := getRulesInputData(policy, e, req)

	rulesFile = getRulesFile(rulesFile, rulesFileName)
	_, coverages := lib.RulesFromJson(rulesFile, initCoverageP(), quotingInputData, []byte(getQuotingData()))
	outJson, out := roundPrices(getOfferPrices(coverages))
	w.Header().Set("Content-Type", "Application/json")
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

func initCoverageP() map[string]*Coverage {

	var coverages = make(map[string]*Coverage)
	coverages["IPI"] = &Coverage{
		Slug:                       "Invalidità Permanente Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["D"] = &Coverage{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["ITI"] = &Coverage{
		Slug:                       "Inabilità Totale Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["DRG"] = &Coverage{
		Slug:                       "Diaria Ricovero / Gessatura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["DC"] = &Coverage{
		Slug:                       "Diaria Convalescenza Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["RSC"] = &Coverage{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["IPM"] = &Coverage{
		Slug:                       "Invalidità Permanente Malattia IPM",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["ASS"] = &Coverage{
		Slug:                       "Assistenza",
		Deductible:                 "0",
		Tax:                        10,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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
	coverages["TL"] = &Coverage{
		Slug:                       "Tutela Legale",
		Deductible:                 "0",
		Tax:                        21.25,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
			"Base": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Your": {
				Deductible:                 "0",
				DeductibleType:             "",
				SumInsuredLimitOfIndemnity: 0.0,
				PremiumNet:                 0.0,
				SelfInsurance:              "0",
			},
			"Premium": {
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

func getOfferPrices(coverage interface{}) Out {
	offerPrice := make(map[string]map[string]*Price)

	offerPrice["Base"] = map[string]*Price{
		"Monthly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		"Yearly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}
	offerPrice["Your"] = map[string]*Price{
		"Monthly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		"Yearly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}
	offerPrice["Premium"] = map[string]*Price{
		"Monthly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
		"Yearly": {
			Net:      0,
			Tax:      0,
			Gross:    0,
			Delta:    0,
			Discount: 0,
		},
	}

	if m, ok := coverage.(map[string]*Coverage); ok {
		for _, c := range m {
			for k, coverageValue := range c.Offer {
				offerPrice[k]["Yearly"].Net += coverageValue.PremiumNet
				offerPrice[k]["Yearly"].Tax += coverageValue.PremiumTaxAmount
				offerPrice[k]["Yearly"].Gross += coverageValue.PremiumGross
				//offerPrice[k]["Yearly"].Discount += coverageValue.
				offerPrice[k]["Monthly"].Net += coverageValue.PremiumNet / 12
				offerPrice[k]["Monthly"].Tax += coverageValue.PremiumTaxAmount / 12
				offerPrice[k]["Monthly"].Gross += coverageValue.PremiumGross / 12
				//offerPrice[k]["Monthly"].Discount += coverageValue.
			}
		}
	}

	out := Out{
		Coverages:  coverage.(map[string]*Coverage),
		OfferPrice: offerPrice,
	}

	return out
}

func roundPrices(out Out) (string, Out) {
	for typePayment, priceStruct := range out.OfferPrice {
		ceilPriceGrossYear := math.Ceil(priceStruct["Yearly"].Gross)
		priceStruct["Yearly"].Delta = ceilPriceGrossYear - priceStruct["Yearly"].Gross
		priceStruct["Yearly"].Gross = ceilPriceGrossYear
		out.Coverages["IPI"].Offer[typePayment].PremiumGross += priceStruct["Yearly"].Delta

		roundPriceGrossMonth := math.Round(priceStruct["Monthly"].Gross)
		priceStruct["Monthly"].Delta = roundPriceGrossMonth - priceStruct["Monthly"].Gross
		priceStruct["Monthly"].Gross = roundPriceGrossMonth
	}

	outJson, _ := json.Marshal(out)
	return string(outJson), out
}
