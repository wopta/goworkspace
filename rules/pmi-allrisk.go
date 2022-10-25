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
		log.Println("rules", string(groule))
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
		"rcCostruction":"` + strings.ToUpper(fil.Elem(0, 26).String()) + `",
		"eletronic":"` + strings.ToUpper(fil.Elem(0, 27).String()) + `",
		"machineFaliure":"` + strings.ToUpper(fil.Elem(0, 28).String()) + `"}`)
	} else {
		enrichByte = []byte(`{}`)
	}

	log.Println(string(enrichByte))
	lib.RulesFromJson(groule, initCoverage(), request, enrichByte, w)

}

func out() []byte {

	return []byte(`{
		"assistance": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "assistance",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"atmospheric-event": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "atmospheric-event",
			"selfInsurance":"0.00",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"building": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "building",
			"IsBase": true,
			"IsYuor": true,
			"IsPremium": true
		},
		"burst-pipe": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "burst-pipe",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"business-interruption": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "business-interruption",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"content": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"cyber": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "cyber",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": true
		},
		"damage-to-goods-course-of-works": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "damage-to-goods-course-of-works",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"damage-to-goods-in-custody": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "damage-to-goods-in-custody",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"defect-liability-12-months": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"selfInsurance":"0.00",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "defect-liability-12-months",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"defect-liability-dm-37-2008": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"selfInsurance":"0.00",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "defect-liability-dm-37-2008",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"defect-liability-workmanships": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"selfInsurance":"0.00",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "defect-liability-workmanships",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"earthquake": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "earthquake",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"electronic-equipment": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "electronic-equipment",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"employers-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "employers-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"environmental-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "environmental-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"glass": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "glass",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"increased-cost-of-working": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "increased-cost-of-working",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"lease-holders-interest": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "lease-holders-interest",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"legal-defence": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity":0.0,
			"Slug": "legal-defence",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": true
		},
		"machinery-breakdown": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "machinery-breakdown",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"power-surge": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "power-surge",
			"selfInsurance":"0.00",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"product-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "product-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"property-damage-due-to-theft": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "property-damage-due-to-theft",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"property-owners-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "property-owners-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"restoration-of-data": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "restoration-of-data",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"river-flood": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "river-flood",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"sociopolitical-event": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "sociopolitical-event",
			"selfInsurance":"0.00",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"software-under-license": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "software-under-license",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"terrorism": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"selfInsurance":"0.00",
			"Slug": "terrorism",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"theft": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "theft",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"third-party-liability": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "third-party-liability",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"third-party-liability-construction-company": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "third-party-liability-construction-company",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"third-party-recourse": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"selfInsurance":"0.00",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "third-party-recourse",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"valuables": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "valuables",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"valuables-in-safe-strongrooms": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"SumInsuredLimitOfIndemnity": 0.0,
			"Slug": "valuables-in-safe-strongrooms",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"water-damage": {
			"Type": "",
			"TypeOfSumInsured": "namedPerils",
			"Deductible": "0",
			"selfInsurance":"0.00",
			"SumInsuredLimitOfIndemnity":0.0,
			"Slug": "water-damage",
			"IsBase": false,
			"IsYuor": false,
			"IsPremium": false
		},
		"error": {
			"message": ""
		}
	}`)
}

type Coverage struct {
	Type                       string
	TypeOfSumInsured           string
	Deductible                 string
	SumInsuredLimitOfIndemnity float64
	Slug                       string
	SelfInsurance              string
	IsBase                     bool
	IsYuor                     bool
	IsPremium                  bool
}

func initCoverage() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["third-party-liability"] = &Coverage{
		Slug:                       "third-party-liability",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-in-custody"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-workmanships"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-12-months"] = &Coverage{
		Slug:                       "defect-liability-12-months",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-dm-37-2008"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-damage-due-to-theft"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-course-of-works"] = &Coverage{
		Slug:                       "damage-to-goods-course-of-works",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["employers-liability"] = &Coverage{
		Slug:                       "employers-liability",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["product-liability"] = &Coverage{
		Slug:                       "product-liability",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-liability-construction-company"] = &Coverage{
		Slug:                       "third-party-liability-construction-company",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["legal-defence"] = &Coverage{
		Slug:                       "legal-defence",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["cyber"] = &Coverage{
		Slug:                       "cyber",
		Type:                       "company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["building"] = &Coverage{
		Slug:                       "building",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["content"] = &Coverage{
		Slug:                       "content",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["lease-holders-interest"] = &Coverage{
		Slug:                       "lease-holders-interest",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["burst-pipe"] = &Coverage{
		Slug:                       "burst-pipe",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["power-surge"] = &Coverage{
		Slug:                       "power-surge",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["atmospheric-event"] = &Coverage{
		Slug:                       "atmospheric-event",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["sociopolitical-event"] = &Coverage{
		Slug:                       "sociopolitical-event",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["terrorism"] = &Coverage{
		Slug:                       "terrorism",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["earthquake"] = &Coverage{
		Slug:                       "earthquake",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["river-flood"] = &Coverage{
		Slug:                       "river-flood",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["water-damage"] = &Coverage{
		Slug:                       "water-damage",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["glass"] = &Coverage{
		Slug:                       "glass",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["machinery-breakdown"] = &Coverage{
		Slug:                       "machinery-breakdown",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-recourse"] = &Coverage{
		Slug:                       "third-party-recourse",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["theft"] = &Coverage{
		Slug:                       "theft",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables-in-safe-strongrooms"] = &Coverage{
		Slug:                       "valuables-in-safe-strongrooms",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables"] = &Coverage{
		Slug:                       "valuables",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["electronic-equipment"] = &Coverage{
		Slug:                       "electronic-equipment",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["increased-cost-of-working"] = &Coverage{
		Slug:                       "increased-cost-of-working",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["restoration-of-data"] = &Coverage{
		Slug:                       "restoration-of-data",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["business-interruption"] = &Coverage{
		Slug:                       "business-interruption",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["environmental-liability"] = &Coverage{
		Slug:                       "environmental-liability",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["assistance"] = &Coverage{
		Slug:                       "assistance",
		Type:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	return res
}
