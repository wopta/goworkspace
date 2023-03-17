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
	coveragesJson, coverages := lib.RulesFromJson(rulesFile, initCoverageP(), quotingInputData, []byte(getQuotingData()))

	return coveragesJson, coverages, nil
}

func getRulesFile(rulesFile []byte, rulesFileName string) []byte {
	switch os.Getenv("env") {
	case "local":
		rulesFile = lib.ErrorByte(ioutil.ReadFile("../function-data/grules/" + rulesFileName))
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

	var res = make(map[string]*Coverage)
	res["IPI"] = &Coverage{
		Slug:                       "Invalidità Permanente Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["D"] = &Coverage{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["ITI"] = &Coverage{
		Slug: "Inabilità Totale Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "7",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["DRG"] = &Coverage{
		Slug: "Diaria Ricovero / Gessatura Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["DC"] = &Coverage{
		Slug:                       "Diaria Convalescenza Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["RSC"] = &Coverage{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["IPM"] = &Coverage{
		Slug: "Invalidità Permanente Malattia IPM",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["ASS"] = &Coverage{
		Slug: "Assistenza",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["TL"] = &Coverage{
		Slug: "Tutela Legale",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Your: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		Premium: &CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			PremiumNet:                 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}

	return res
}
func getQuotingData() string {

	return string(lib.GetByteByEnv("quote/persona-tassi.json", false))
}
