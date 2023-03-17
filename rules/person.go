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
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

const (
	base    = "base"
	your    = "your"
	premium = "premium"
	monthly = "monthly"
	yearly  = "yearly"
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
	coverages["D"] = &Coverage{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["ITI"] = &Coverage{
		Slug:                       "Inabilità Totale Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["DRG"] = &Coverage{
		Slug:                       "Diaria Ricovero / Gessatura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["DC"] = &Coverage{
		Slug:                       "Diaria Convalescenza Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["RSC"] = &Coverage{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["IPM"] = &Coverage{
		Slug:                       "Invalidità Permanente Malattia IPM",
		Deductible:                 "0",
		Tax:                        2.5,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["ASS"] = &Coverage{
		Slug:                       "Assistenza",
		Deductible:                 "0",
		Tax:                        10,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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
	coverages["TL"] = &Coverage{
		Slug:                       "Tutela Legale",
		Deductible:                 "0",
		Tax:                        21.25,
		SumInsuredLimitOfIndemnity: 0.0,
		Offer: map[string]*CoverageValue{
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

	if coveragesStruct, ok := coverages.(map[string]*Coverage); ok {
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
		Coverages:  coverages.(map[string]*Coverage),
		OfferPrice: offerPrice,
	}

	return out
}

func roundPrices(out *Out) *Out {
	for offerType, priceStruct := range out.OfferPrice {
		log.Println("Offer type: " + offerType)
		log.Println("Initial IPI Price Gross: " + strconv.FormatFloat(out.Coverages["IPI"].Offer[offerType].PremiumGross, 'f', -1, 64))
		log.Println("PT: " + strconv.FormatFloat(priceStruct[yearly].Gross, 'f', -1, 64))
		log.Println("Pm: " + strconv.FormatFloat(priceStruct[monthly].Gross, 'f', -1, 64))
		ceilPriceGrossYear := math.Ceil(priceStruct[yearly].Gross)
		priceStruct[yearly].Delta = ceilPriceGrossYear - priceStruct[yearly].Gross
		priceStruct[yearly].Gross = ceilPriceGrossYear
		if out.Coverages["IPI"].Offer[offerType].PremiumGross > 0 {
			out.Coverages["IPI"].Offer[offerType].PremiumGross += priceStruct[yearly].Delta
		} else {
			out.Coverages["DRG"].Offer[offerType].PremiumGross += priceStruct[yearly].Delta
		}

		roundPriceGrossMonth := math.Round(priceStruct[monthly].Gross)
		priceStruct[monthly].Delta = roundPriceGrossMonth - priceStruct[monthly].Gross
		priceStruct[monthly].Gross = roundPriceGrossMonth

		log.Println("PGa: " + strconv.FormatFloat(priceStruct[yearly].Gross, 'f', -1, 64))
		log.Println("PGm: " + strconv.FormatFloat(priceStruct[monthly].Gross, 'f', -1, 64))
		log.Println("Monthly Delta: " + strconv.FormatFloat(priceStruct[monthly].Delta, 'f', -1, 64))
		log.Println("Yearly Delta: " + strconv.FormatFloat(priceStruct[yearly].Delta, 'f', -1, 64))
		log.Println("Final IPI Price Gross: " + strconv.FormatFloat(out.Coverages["IPI"].Offer[offerType].PremiumGross, 'f', -1, 64))
		log.Println()
	}

	return out
}

func filterOffers(out *Out) (string, *Out) {
	toBeDeleted := make([]string, 0)
	for offerType, priceStruct := range out.OfferPrice {
		if priceStruct[yearly].Gross < 120 || priceStruct[monthly].Gross < 50 {
			toBeDeleted = append(toBeDeleted, offerType)
		}
	}

	log.Println("Offers to be deleted: ", toBeDeleted)

	for _, offerType := range toBeDeleted {
		delete(out.OfferPrice, offerType)
		for _, coverage := range out.Coverages {
			delete(coverage.Offer, offerType)
		}
	}

	outJson, _ := json.Marshal(out)
	return string(outJson), out
}
