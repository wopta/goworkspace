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
	lib.RulesFromJson(groule, initCoverageP(), request, []byte(getData()), w)

}

func initCoverageP() map[string]*models.Coverage {

	var res = make(map[string]*models.Coverage)
	res["IPI"] = &models.Coverage{
		Slug:                       "Invalidità Permanente Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["D"] = &models.Coverage{
		Slug:                       "Decesso Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["ITI"] = &models.Coverage{
		Slug: "Inabilità Totale Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["DRG"] = &models.Coverage{
		Slug: "Diaria Ricovero / Gessatura Infortunio",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["DC"] = &models.Coverage{
		Slug:                       "Diaria Convalescenza Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["RSC"] = &models.Coverage{
		Slug:                       "Rimborso spese di cura Infortunio",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["IPM"] = &models.Coverage{
		Slug: "Invalidità Permanente Malattia IPM",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["ASS"] = &models.Coverage{
		Slug: "Assistenza",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["TL"] = &models.Coverage{
		Slug: "third-party-liability",

		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		Base: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Your: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		Premium: &models.CoverageValue{
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}

	return res
}
func getData() string {
	return `{
		"IPI": {
			"extra":{
			"1": {
				"asorbable": {
					"5": 1.08,
					"10": 0.79
				},
				"absolute": {
					"3": 1.13,
					"5": 0.90,
					"10": 0.74
				}
			}
		},
		"professional":{
			"1": {
				"asorbable": {
					"5": 0.68,
					"10": 0.50
				},
				"absolute": {
					"3": 0.71,
					"5": 0.57,
					"10": 0.47
				}
			},
			"2": {
				"asorbable": {
					"5": 0.88,
					"10": 0.65
				},
				"absolute": {
					"3": 0.92,
					"5": 0.74,
					"10": 0.60
				}
			},
			"3": {
				"asorbable": {
					"5": 1.25,
					"10": 0.92
				},
				"absolute": {
					"3": 1.32,
					"5": 1.06,
					"10": 0.85
				}
			},
			"4": {
				"asorbable": {
					"5": 1.45,
					"10": 1.07
				},
				"absolute": {
					"3": 1.52,
					"5": 1.22,
					"10": 0.99
				}
			}
		},
		"extraprof":{
			"1": {
				"asorbable": {
					"5": 1.14,
					"10": 0.83
				},
				"absolute": {
					"3": 1.19,
					"5": 0.95,
					"10": 0.78
				}
			},
			"2": {
				"asorbable": {
					"5": 1.46,
					"10": 1.08
				},
				"absolute": {
					"3": 1.53,
					"5": 1.23,
					"10": 1.00
				}
			},
			"3": {
				"asorbable": {
					"5": 2.08,
					"10": 1.53
				},
				"absolute": {
					"3": 2.20,
					"5": 1.76,
					"10": 1.42
				}
			},
			"4": {
				"asorbable": {
					"5": 2.41,
					"10": 1.78
				},
				"absolute": {
					"3": 2.54,
					"5": 2.03,
					"10":1.65
				}
			}
		}
	}
	
	}`
}
