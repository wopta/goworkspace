package rules

/*



INPUT:
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

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func Person(w http.ResponseWriter, r *http.Request) {
	var (
		result map[string]interface{}
		groule []byte
	)
	const (
		rulesFile = "person.json"
	)
	log.Println("Person")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(request), &result)
	//swich source by env for retive data
	switch os.Getenv("env") {
	case "local":
		groule = lib.ErrorByte(ioutil.ReadFile("function-data/grules/" + rulesFile))

	case "dev":
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFile, "")

	case "prod":
		groule = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFile, "")

	default:

	}
	lib.RulesFromJson(groule, initCoverageP(), request, []byte(`{}`), w)

}

func initCoverageP() map[string]*models.Coverage {

	var res = make(map[string]*models.Coverage)
	res["IPI"] = &models.Coverage{
		Slug:                       "Invalidità Permanente Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,

		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["D"] = &models.Coverage{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["ITI"] = &models.Coverage{
		Slug: "Inabilità Totale Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["DRG"] = &models.Coverage{
		Slug: "Diaria Ricovero / Gessatura Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["DC"] = &models.Coverage{
		Slug: "Diaria Convalescenza Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["RSC"] = &models.Coverage{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["IPM"] = &models.Coverage{
		Slug: "Invalidità Permanente Malattia IPM",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["ASS"] = &models.Coverage{
		Slug: "Assistenza",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}
	res["TL"] = &models.Coverage{
		Slug: "third-party-liability",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYuor:    false,
		IsPremium: false,
	}

	return res
}
