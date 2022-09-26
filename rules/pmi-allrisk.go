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

func PmiAllrisk(w http.ResponseWriter, r *http.Request) {
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
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFile, "")
		ricAteco = lib.GetFromStorage("function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	default:

	}
	df := lib.CsvToDataframe(ricAteco)
	fil := df.Filter(
		dataframe.F{Colidx: 5, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: result["ateco"]},
	)
	log.Println("filtered", fil.Nrow())
	log.Println("filtered", fil.Ncol())
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
		"damageTheft ":"` + strings.ToUpper(fil.Elem(0, 23).String()) + `",
		"damageThing":"` + strings.ToUpper(fil.Elem(0, 24).String()) + `",
		"rcCostruction":"` + strings.ToUpper(fil.Elem(0, 26).String()) + `",
		"eletronic":"` + strings.ToUpper(fil.Elem(0, 27).String()) + `",
		"machineFaliure":"` + strings.ToUpper(fil.Elem(0, 28).String()) + `"}`)
	}

	log.Println(string(enrichByte))
	lib.RulesFromJson(groule, out(), request, enrichByte, w)

}

func out() []byte {

	return []byte(`{
		"assistance": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "assistance",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"atmospheric-event": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "atmospheric-event",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"building": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "building",
			"IsBase": true,
			"IsYuor": true,
			"IsPremium": true
		},
		"burst-pipe": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "burst-pipe",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"business-interruption": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "business-interruption",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"content": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "content",
			"IsBase": true,
			"IsYuor": true,
			"IsPremium": true
		},
		"cyber": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 500000,
			"Slug": "cyber",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": true
		},
		"damage-to-goods-course-of-works": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "damage-to-goods-course-of-works",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"damage-to-goods-in-custody": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 1000000,
			"Slug": "damage-to-goods-in-custody",
			"IsBase": false,
			"IsYuor": true,
			"IsPremium": true
		},
		"defect-liability-12-months": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "defect-liability-12-months",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"defect-liability-dm-37-2008": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "defect-liability-dm-37-2008",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"defect-liability-workmanships": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "defect-liability-workmanships",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"earthquake": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "earthquake",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"electronic-equipment": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "electronic-equipment",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"employers-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 1000000,
			"Slug": "employers-liability",
			"IsBase": true,
			"IsYuor": true,
			"IsPremium": true
		},
		"environmental-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "environmental-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"glass": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "glass",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"increased-cost-of-working": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "increased-cost-of-working",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"lease-holders-interest": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "lease-holders-interest",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"legal-defence": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 500000,
			"Slug": "legal-defence",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": true
		},
		"machinery-breakdown": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "machinery-breakdown",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"power-surge": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "power-surge",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"product-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "product-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"property-damage-due-to-theft": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "property-damage-due-to-theft",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"property-owners-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "property-owners-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"restoration-of-data": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "restoration-of-data",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"river-flood": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "river-flood",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"sociopolitical-event": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "sociopolitical-event",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"software-under-license": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "software-under-license",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"terrorism": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "terrorism",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"theft": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "theft",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"third-party-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 1000000,
			"Slug": "third-party-liability",
			"IsBase": true,
			"IsYuor": true,
			"IsPremium": true
		},
		"third-party-liability-construction-company": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 1000000,
			"Slug": "third-party-liability-construction-company",
			"IsBase": false,
			"IsYuor": true,
			"IsPremium": true
		},
		"third-party-recourse": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "third-party-recourse",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"valuables": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "valuables",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"valuables-in-safe-strongrooms": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "valuables-in-safe-strongrooms",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"water-damage": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0,
			"Slug": "water-damage",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"error": {
			"message": "",
		}
	}`)
}
