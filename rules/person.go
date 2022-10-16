package rules

/*

{
	age
	work
	worktype
	coverageType
	childrenScool
	issue1500
	riskInLifeIs
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
)

func Person(w http.ResponseWriter, r *http.Request) {
	var (
		result map[string]interface{}
		groule []byte
	)
	const (
		rulesFile = "person.json"
	)
	log.Println("PmiAllrisk")
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

func initCoverageP() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["IPI"] = &Coverage{
		Slug:                       "third-party-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["D"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["ITI"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["DR"] = &Coverage{
		Slug:                       "defect-liability-12-months",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["DC"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["RSC"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["IPM"] = &Coverage{
		Slug:                       "damage-to-goods-course-of-works",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["ASS"] = &Coverage{
		Slug:                       "employers-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["TL"] = &Coverage{
		Slug:                       "product-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}

	return res
}
