package rules

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
)

func PmiAllrisk(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	var (
		result   map[string]interface{}
		groule   []byte
		ricAteco []byte
	)
	const (
		rulesFile = "pmi-allrisk.json"
	)
	log.Println("PmiAllrisk")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	// Unmarshal or Decode the JSON to the interface.

	json.Unmarshal([]byte(request), &result)
	//swich source by env for retive data
	switch os.Getenv("env") {
	case "local":
		groule = lib.ErrorByte(ioutil.ReadFile("function-data/grules/" + rulesFile))

		ricAteco = lib.ErrorByte(ioutil.ReadFile("function-data/data/rules/Riclassificazione_Ateco.csv"))
	case "dev":
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFile, "")
		ricAteco = lib.GetFromStorage("function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	case "prod":
		groule = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFile, "")
		ricAteco = lib.GetFromStorage("core-350507-function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	default:

	}
	df := lib.CsvToDataframe(ricAteco)
	fil := df.Filter(
		dataframe.F{Colidx: 5, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: result["ateco"]},
	)
	log.Println("filtered row", fil.Nrow())
	log.Println("filtered col", fil.Ncol())
	var enrichByte []byte

	if fil.Nrow() > 0 {
		enrichByte = []byte(`{	"atecoMacro":"` + strings.ToUpper(fil.Elem(0, 0).String()) + `",
		"atecoSub":"` + strings.ToUpper(fil.Elem(0, 1).String()) + `",
		"atecoDesc":"` + strings.ToUpper(fil.Elem(0, 2).String()) + `",
		"businessSector":"` + strings.ToUpper(fil.Elem(0, 3).String()) + `",
		"fire":"` + strings.ToUpper(fil.Elem(0, 14).String()) + `",
		"fireLow500k":"` + strings.ToUpper(fil.Elem(0, 5).String()) + `",
		"fireUp500k":"` + strings.ToUpper(fil.Elem(0, 6).String()) + `",
		"theft":"` + strings.ToUpper(fil.Elem(0, 15).String()) + `",
		"thefteLow500k ":"` + strings.ToUpper(fil.Elem(0, 8).String()) + `",
		"theftUp500k":"` + strings.ToUpper(fil.Elem(0, 9).String()) + `",
		"rct":"` + strings.ToUpper(fil.Elem(0, 16).String()) + `",
		"rco":"` + strings.ToUpper(fil.Elem(0, 17).String()) + `",
		"rcoProd":"` + strings.ToUpper(fil.Elem(0, 18).String()) + `",
		"rcVehicle":"` + strings.ToUpper(fil.Elem(0, 19).String()) + `",
		"rcpo":"` + strings.ToUpper(fil.Elem(0, 20).String()) + `",
		"rcp12":"` + strings.ToUpper(strings.ToUpper(fil.Elem(0, 21).String())) + `",
		"rcp2008":"` + strings.ToUpper(fil.Elem(0, 22).String()) + `",
		"damageTheft":"` + strings.ToUpper(fil.Elem(0, 23).String()) + `",
		"damageThing":"` + strings.ToUpper(fil.Elem(0, 24).String()) + `",
		"rcCostruction":"` + strings.ToUpper(fil.Elem(0, 25).String()) + `",
		"eletronic":"` + strings.ToUpper(fil.Elem(0, 27).String()) + `",
		"machineFaliure":"` + strings.ToUpper(fil.Elem(0, 28).String()) + `"}`)
	} else {
		enrichByte = []byte(`{}`)
	}

	log.Println(string(enrichByte))
	s, i := lib.RulesFromJson(groule, initCoverage(), request, enrichByte)

	return s, i
}

func initCoverage() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["third-party-liability"] = &Coverage{
		Slug:                       "third-party-liability",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Tax:                        22.25,
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-in-custody"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		SelfInsurance:              "0",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-workmanships"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		SelfInsurance:              "0",
		Tax:                        22.25,
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-12-months"] = &Coverage{
		Slug:                       "defect-liability-12-months",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-dm-37-2008"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-damage-due-to-theft"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-course-of-works"] = &Coverage{
		Slug:             "damage-to-goods-course-of-works",
		Type:             "company",
		TypeOfSumInsured: "firstLoss",

		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		SelfInsurance:              "0",
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["employers-liability"] = &Coverage{
		Slug:                       "employers-liability",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["product-liability"] = &Coverage{
		Slug:                       "product-liability",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-liability-construction-company"] = &Coverage{
		Slug:                       "third-party-liability-construction-company",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["legal-defence"] = &Coverage{
		Slug:         "legal-defence",
		Type:         "company",
		LegalDefence: "basic",
		Tax:          21.25,
		IsBase:       false,
		IsYuor:       false,
		IsPremium:    false,
	}
	res["cyber"] = &Coverage{
		Slug:                       "cyber",
		Type:                       "company",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		Taxes:                      []Tax{{Tax: 22.25, Percentage: 40.0}, {Tax: 21.25, Percentage: 60.0}},
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["building"] = &Coverage{
		Slug:                       "building",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["content"] = &Coverage{
		Slug:                       "content",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["lease-holders-interest"] = &Coverage{
		Slug:                       "lease-holders-interest",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["burst-pipe"] = &Coverage{
		Slug:                       "burst-pipe",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["power-surge"] = &Coverage{
		Slug:                       "power-surge",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["atmospheric-event"] = &Coverage{
		Slug:                       "atmospheric-event",
		Type:                       "building",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["sociopolitical-event"] = &Coverage{
		Slug:                       "sociopolitical-event",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["terrorism"] = &Coverage{
		Slug:                       "terrorism",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["earthquake"] = &Coverage{
		Slug:                       "earthquake",
		Type:                       "building",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["river-flood"] = &Coverage{
		Slug:                       "river-flood",
		Type:                       "building",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["water-damage"] = &Coverage{
		Slug:                       "water-damage",
		Type:                       "building",
		TypeOfSumInsured:           "replacementValue",
		SelfInsurance:              "0",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["glass"] = &Coverage{
		Slug:                       "glass",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["machinery-breakdown"] = &Coverage{
		Slug:                       "machinery-breakdown",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-recourse"] = &Coverage{
		Slug:                       "third-party-recourse",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["theft"] = &Coverage{
		Slug:                       "theft",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables-in-safe-strongrooms"] = &Coverage{
		Slug:                       "valuables-in-safe-strongrooms",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables"] = &Coverage{
		Slug:                       "valuables",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["electronic-equipment"] = &Coverage{
		Slug:                       "electronic-equipment",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["increased-cost-of-working"] = &Coverage{
		Slug:                       "increased-cost-of-working",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["restoration-of-data"] = &Coverage{
		Slug:                       "restoration-of-data",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["business-interruption"] = &Coverage{
		Slug:                       "business-interruption",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        21.25,
		DailyAllowance:             "250",
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["environmental-liability"] = &Coverage{
		Slug:                       "environmental-liability",
		Type:                       "building",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["assistance"] = &Coverage{
		Slug:       "assistance",
		Type:       "building",
		Assistance: "yes",
		Tax:        10.00,
		IsBase:     false,
		IsYuor:     false,
		IsPremium:  false,
	}
	return res
}
