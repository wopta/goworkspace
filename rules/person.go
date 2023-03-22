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
	"net/http"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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
	coveragesJson, coverages := rulesFromJson(rulesFile, initCoverageP(), quotingInputData, []byte(getQuotingData()))

	return coveragesJson, coverages, nil
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

func initCoverageP() *models.RuleOut {
	var coverages = make(map[string]*models.CoverageOut)
	offerPrice := make(map[string]map[string]*models.Price)

	coverages = map[string]*models.CoverageOut{
		"IPI": {
			Slug:                       "Invalidità Permanente Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"D": {
			Slug:                       "Decesso Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"ITI": {
			Slug:                       "Inabilità Totale Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"DRG": {
			Slug:                       "Diaria Ricovero / Gessatura Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"DC": {
			Slug:                       "Diaria Convalescenza Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"RSC": {
			Slug:                       "Rimborso spese di cura Infortunio",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"IPM": {
			Slug:                       "Invalidità Permanente Malattia IPM",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(2.5),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"ASS": {
			Slug:                       "Assistenza",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(10.00),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
		"TL": {
			Slug:                       "Tutela Legale",
			Deductible:                 "0",
			Tax:                        decimal.NewFromFloat(21.25),
			SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
			Offer: map[string]*models.CoverageValueOut{
				base: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				your: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
				premium: {
					Deductible:                 "0",
					DeductibleType:             "",
					SumInsuredLimitOfIndemnity: decimal.NewFromFloat(0.0),
					PremiumNetYearly:           decimal.NewFromFloat(0.0),
					PremiumTaxAmountYearly:     decimal.NewFromFloat(0.0),
					PremiumGrossYearly:         decimal.NewFromFloat(0.0),
					PremiumNetMonthly:          decimal.NewFromFloat(0.0),
					PremiumTaxAmountMonthly:    decimal.NewFromFloat(0.0),
					PremiumGrossMonthly:        decimal.NewFromFloat(0.0),
					SelfInsurance:              "0",
				},
			},
			IsBase:    false,
			IsYour:    false,
			IsPremium: false,
		},
	}

	offerPrice[base] = map[string]*models.Price{
		monthly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
		yearly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
	}
	offerPrice[your] = map[string]*models.Price{
		monthly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
		yearly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
	}
	offerPrice[premium] = map[string]*models.Price{
		monthly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
		yearly: {
			Net:      decimal.NewFromFloat(0.0),
			Tax:      decimal.NewFromFloat(0.0),
			Gross:    decimal.NewFromFloat(0.0),
			Delta:    decimal.NewFromFloat(0.0),
			Discount: decimal.NewFromFloat(0.0),
		},
	}

	return &models.RuleOut{
		Coverages:  coverages,
		OfferPrice: offerPrice,
	}
}
func getQuotingData() string {
	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
