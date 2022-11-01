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
)

func Person(w http.ResponseWriter, r *http.Request) (string, interface{}) {
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
	s, i := lib.RulesFromJson(groule, initCoverageP(), request, []byte(getData()))
	return s, i
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
			Deductible:                 "0",
			DeductibleType:             "",
			SumInsuredLimitOfIndemnity: 0.0,
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
			SelfInsurance:              "0"},
		IsBase:    false,
		IsYour:    false,
		IsPremium: false,
	}
	res["TL"] = &Coverage{
		Slug: "third-party-liability",

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

	return res
}
func getData() string {
	return `{
		"IPI": {
			"extra": {
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
			"professional": {
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
			"extraprof": {
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
						"10": 1.65
					}
				}
			}
		},
		"DRC": {
			"extra": {
				"1": 0.81
			},
			"professional": {
				"1": 0.55,
				"2": 0.66,
				"3": 0.88,
				"4": 1.04
			},
			"extraprof": {
				"1": 0.87,
				"2": 1.04,
				"3": 1.40,
				"4": 1.65
			}
		},
		"MI": {
			"extra": {
				"1": 0.62
			},
			"professional": {
				"1": 0.38,
				"2": 0.50,
				"3": 0.70,
				"4": 0.89
			},
			"extraprof": {
				"1": 0.64,
				"2": 0.83,
				"3": 1.16,
				"4": 1.48
			}
		},
		"IT": {
			"1": {
				"7": 4.17,
				"15": 3.32
			},
			"2": {
				"7": 4.98,
				"15": 3.98
			},
			"3": {
				"7": 5.80,
				"15": 4.64
			},
			"4": {
				"7": 7.46,
				"15": 5.97
			}
		},
		"IPM":  {
			"18": 0.83,
			"19": 0.83,
			"20": 0.83,
			"21": 0.83,
			"22": 0.83,
			"23": 0.83,
			"24": 0.83,
			"25": 0.83,
			"26": 0.83,
			"27": 0.83,
			"28": 0.83,
			"29": 0.83,
			"30": 0.83,
			"31": 1.00,
			"32": 1.00,
			"33": 1.00,
			"34": 1.00,
			"35": 1.00,
			"36": 1.50,
			"37": 1.50,
			"38": 1.50,
			"39": 1.50,
			"40": 1.50,
			"41": 2.06,
			"42": 2.06,
			"43": 2.06,
			"44": 2.06,
			"45": 2.06,
			"46": 2.73,
			"47": 2.73,
			"48": 2.73,
			"49": 2.73,
			"50": 2.73,
			"51": 3.32,
			"52": 3.32,
			"53": 3.32,
			"54": 3.32,
			"55": 3.32,
			"56": 4.15,
			"57": 4.15,
			"58": 4.15,
			"59": 4.15,
			"60": 4.15,
			"61": 4.15,
			"62": 4.15,
			"63": 4.15,
			"64": 4.15,
			"65": 4.15
		},
		"DC": {
			"extra": {
				"1": 0.33
			},
			"professional": {
				"1": 0.20,
				"2": 0.25,
				"3": 0.33,
				"4": 0.40
			},
			"extraprof": {
				"1": 0.32,
				"2": 0.40,
				"3": 0.53,
				"4": 0.63
			}
		}
		,
		"RSC": {
			"extra": {
				"1": {
					"2500":31.59,
					"5000":56.95,
					"10000":104.70,
					"15000":148.91
				}
			},
			"professional": {
				"1": {
					"2500": 19.49,
					"5000":35.13,
					"10000":64.60,
					"15000":91.87
				},   
				"2": {
					"2500": 25.39,
					"5000":45.77,
					"10000":84.15,
					"15000":119.67
				},   
				"3": {
					"2500": 36.02,
					"5000":64.92,
					"10000":119.36,
					"15000":169.76
				},   
				"4": {
					"2500": 45.47,
					"5000":81.97,
					"10000":150.70,
					"15000":214.34
				}
			},
			"extraprof": { 
			"1": {
				"2500": 32.49,
				"5000":58.55,
				"10000":107.67,
				"15000":153.12
			},   
			"2": {
				"2500": 42.32,
				"5000":76.28,
				"10000":140.25,
				"15000":199.45
			},   
			"3": {
				"2500": 60.03,
				"5000":108.20,
				"10000":198.94,
				"15000":282.93
			},   
			"4": {
				"2500": 75.79,
				"5000":136.62,
				"10000":251.17,
				"15000":357.23
			}
			}
		}
	}`
}
