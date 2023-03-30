package rules

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/mohae/deepcopy"
	lib "github.com/wopta/goworkspace/lib"
	q "github.com/wopta/goworkspace/quote"
)

func PmiAllrisk(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
		groule = lib.ErrorByte(ioutil.ReadFile("../function-data/dev/grules/" + rulesFile))
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
		enrichByte = []byte(`{
		"RcpRD":"` + strings.ToUpper(fil.Elem(0, 13).String()) + `",
		"RcoRD":"` + strings.ToUpper(fil.Elem(0, 12).String()) + `",
		"RctRD":"` + strings.ToUpper(fil.Elem(0, 11).String()) + `",
		"theft500":"` + strings.ToUpper(fil.Elem(0, 8).String()) + `",
		"fire500":"` + strings.ToUpper(fil.Elem(0, 6).String()) + `",
		"atecoMacro":"` + strings.ToUpper(fil.Elem(0, 0).String()) + `",
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
	log.Println("enrichByte:" + string(enrichByte))
	_, i := lib.RulesFromJson(groule, initCoverage(), request, enrichByte)

	var g map[string]*Coverage
	p, e := json.Marshal(i)
	log.Println("rules result::" + string(p))
	g = i.(map[string]*Coverage)
	ateco := strings.ReplaceAll(result["ateco"].(string), ".", "")
	var alarm string
	if result["isAllarm"].(bool) {
		alarm = "yes"
	} else {
		alarm = "no"
	}
	q0 := q.MunichReQuoteRequest{
		SME: q.SME{
			SubproductID: 35.0,
			UWRole:       "Agent",
			Ateco:        ateco,
			Company: q.Company{
				Country:   "Italy",
				Vatnumber: q.Vatnumber{Value: result["vat"].(string)},
				OpreEur: q.Employees{
					Value: int64(result["revenue"].(float64)),
				},
				Employees: q.Employees{
					Value: int64(result["employer"].(float64)),
				},
			},
			Answers: q.Answers{
				Step1: []q.Step1{},
				Step2: []q.Step2{{
					BuildingID: "2",
					Value: q.Step2Value{
						BuildingType:     result["constructionMaterial"].(string),
						NumberOfFloors:   result["floor"].(string),
						ConstructionYear: result["buildingYear"].(string),
						Alarm:            alarm,
						TypeOfInsurance:  "namedPerils",
						Ateco:            ateco,
						Postcode:         result["postcode"].(string),
						Province:         result["province"].(string),
					},
				}},
			},
		},
	}

	q1 := deepcopy.Copy(q0).(q.MunichReQuoteRequest)
	q2 := deepcopy.Copy(q0).(q.MunichReQuoteRequest)
	q3 := deepcopy.Copy(q0).(q.MunichReQuoteRequest)

	for _, v := range g {
		var typeOfSumInsured string

		if v.TypeOfSumInsured == "" {

			if result["isPra"].(bool) {
				typeOfSumInsured = "firstLoss"
			} else {
				typeOfSumInsured = "replacementValue"
			}
		} else {

			typeOfSumInsured = v.TypeOfSumInsured
		}
		st1value := q.Value{}

		if v.Slug == "legal-defence" {
			st1value.LegalDefence = &v.LegalDefence
		} else if v.Slug == "assistance" {

			yes := "yes"
			st1value.Assistance = &yes
		} else {
			st1value.TypeOfSumInsured = &typeOfSumInsured
			st1value.Deductible = &v.Deductible
			st1value.SumInsuredLimitOfIndemnity = &v.SumInsuredLimitOfIndemnity
			st1value.Deductible = &v.Deductible
			if v.Slug == "business-interruption" {
				st1value.DailyAllowance = &v.DailyAllowance
			}
			if v.SelfInsurance != "" {
				st1value.SelfInsurance = &v.SelfInsurance
			}
		}
		if v.Type == "company" {
			if v.IsBase {

				q1.SME.Answers.Step1 = append(q1.SME.Answers.Step1, q.Step1{Slug: v.Slug, Value: st1value})
			}
			if v.IsYour {
				q2.SME.Answers.Step1 = append(q2.SME.Answers.Step1, q.Step1{Slug: v.Slug, Value: st1value})
			}
			if v.IsPremium {
				q3.SME.Answers.Step1 = append(q3.SME.Answers.Step1, q.Step1{Slug: v.Slug, Value: st1value})
			}
		}
		if v.Type == "building" {

			if v.IsBase {
				q1.SME.Answers.Step2[0].Value.Answer = append(q1.SME.Answers.Step2[0].Value.Answer, q.Answer{Slug: v.Slug, Value: st1value})
			}
			if v.IsYour {
				q2.SME.Answers.Step2[0].Value.Answer = append(q2.SME.Answers.Step2[0].Value.Answer, q.Answer{Slug: v.Slug, Value: st1value})
			}
			if v.IsPremium {
				q3.SME.Answers.Step2[0].Value.Answer = append(q3.SME.Answers.Step2[0].Value.Answer, q.Answer{Slug: v.Slug, Value: st1value})
			}
		}

	}
	//m1, _ := q1.Marshal()
	m3, e := q3.Marshal()
	log.Println(string(m3))
	q3resbyte := <-q.PmiMunich([]byte(m3))
	log.Println(q3resbyte)
	var q3res q.MunichReQuoteResponse
	q3res, e = q3res.Unmarshal([]byte(q3resbyte))
	log.Println(q3res)
	log.Println(e)
	var (
		sumBase int
		sumY    int
		sumP    int
	)
	if !lib.StructIsEmpty(q3res) {
		for _, r := range q3res.Result.Answers.Step1 {
			g[r.Slug].PriceNett = r.Value.PremiumNet
			if g[r.Slug].Tax == 0 {
				var taxsum float64
				for _, t := range g[r.Slug].Taxes {
					premiumPerc := ((r.Value.PremiumNet * t.Percentage) / 100)
					taxsum = taxsum + ((premiumPerc * t.Tax) / 100)
				}

				g[r.Slug].PriceGross = taxsum
			} else {
				g[r.Slug].PriceGross = r.Value.PremiumNet + ((r.Value.PremiumNet * g[r.Slug].Tax) / 100)
			}

		}
		for _, r := range q3res.Result.Answers.Step2[0].Value {
			g[r.Slug].PriceNett = r.Value.PremiumNet
			if g[r.Slug].Tax == 0 {
				var taxsum float64
				for _, t := range g[r.Slug].Taxes {
					premiumPerc := ((r.Value.PremiumNet * t.Percentage) / 100)
					taxsum = taxsum + ((premiumPerc * t.Tax) / 100)
				}
				g[r.Slug].PriceGross = taxsum
			} else {
				g[r.Slug].PriceGross = r.Value.PremiumNet + ((r.Value.PremiumNet * g[r.Slug].Tax) / 100)
			}

		}

		for _, t := range g {

			if t.IsBase {
				sumBase = sumBase + int(t.PriceGross)
			}
			if t.IsYuor || t.IsYour {
				sumY = sumY + int(t.PriceGross)
			}
			if t.IsPremium {
				sumP = sumP + int(t.PriceGross)
			}
		}
	}
	b, e := json.Marshal(g)
	requestQuote, e := q3.Marshal()
	responseQuote, e := q3res.Marshal()
	var status int64
	if lib.StructIsEmpty(q3res) {
		status = 500
	} else {
		status = 200
	}
	//println(civil.DateTimeOf(time.Now()))
	//l:="2016-12-17T17:16:27"
	dwh := q.MunichReQuotePmiDWHCall{
		Status:            status,
		CreationDate:      civil.DateTimeOf(time.Now()),
		RequestRules:      string(request),
		ResponseRules:     string(b),
		RequestQuote:      string(requestQuote),
		ResponseQuote:     string(responseQuote),
		RequestRulesJson:  string(request),
		ResponseRulesJson: string(b),
		RequestQuoteJson:  string(requestQuote),
		ResponseQuoteJson: string(responseQuote),
		Base:              int64(sumBase),
		Your:              int64(sumY),
		Premium:           int64(sumP),
	}
	e = lib.InsertRowsBigQuery("wopta", "policy-rules-log", dwh)

	println("end")

	return string(b), i, e
}

func initCoverage() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["third-party-liability"] = &Coverage{
		Slug:                       "third-party-liability",
		Type:                       "company",
		CompanyCodec:               "CT",
		Group:                      "RCT",
		IsExtension:                false,
		CompanyName:                "Responsabilità Civile Terzi",
		TypeOfSumInsured:           "",
		Tax:                        22.25,
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0.0,
		SelfInsuranceDesc:          "",
		Description:                "Danni involontariamente causati a terzi per danni a cose o persone di cui sia responsabile a termini di legge (R.C.T.). La garanzia include, ma non si limita a questi, i  danni: a veicoli di terzi e prestatori di lavoro; a cose in consegna e custodia; a cose nell'ambito di esecuzione dei lavori; a cose di terzi sollevate, caricate, scaricate, movimentate, trasportate o rimorchiate; a mezzi di trasporto sotto carico e scarico; da interruzione o sospensione di attività di terzi; da smercio; da committenza autoveicoli; da responsabilità civile personale addetti; da attività di commercio ambulante; da lavori presso terzi",
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-in-custody"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                true,
		CompanyCodec:               "DV",
		CompanyName:                "Danni ai Veicoli in consegna e custodia",
		TypeOfSumInsured:           "",
		SelfInsurance:              "0",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Per responsabilità relative a danni a veicoli in consegna, custodia o comunque detenuti, compresa errata erogazione di carburante, riparazione e manutenzione o movimentazione/caduta ponte sollevamento",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-workmanships"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                true,
		TypeOfSumInsured:           "",
		CompanyCodec:               "CP",
		CompanyName:                "Responsabilità civile postuma officine",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Per danni dopo l'ultimazione dei lavori, subiti o causati da veicoli a motore riparati, revisionati o manutenuti, compresi interventi su pneumatici",
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
		Group:                      "RCT",
		IsExtension:                true,
		CompanyCodec:               "12",
		CompanyName:                "Responsabilità civile postuma 12 mesi",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Per danni  in qualità di installatore, manutentore o riparatore, causati dopo ultimazione lavori dalle cose installate, riparate o manutenute (NO Ateco Edili 41. 42. 43.)",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-dm-37-2008"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                true,
		CompanyCodec:               "DM",
		CompanyName:                "Responsabilità civile D.M.37/2008",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		Tax:                        22.25,
		SelfInsuranceDesc:          "",
		Description:                "Per danni  in qualità di installatore, manutentore o riparatore, causati dopo ultimazione lavori dalle cose installate, riparate o manutenute (Solo per Ateco Edili e attività soggette a DM 37/2008)",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-damage-due-to-theft"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                true,
		CompanyCodec:               "DF",
		CompanyName:                "Danni da furto",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Per danni da furto commesso da soggetti che abbiano utilizzato impalcature, attrezzature fisse o ponteggi eretti o fatti erigere dall'assicurato",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-course-of-works"] = &Coverage{
		Slug: "damage-to-goods-course-of-works",

		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                true,
		CompanyCodec:               "DC",
		CompanyName:                "Danni alle cose su cui si eseguono lavori",
		TypeOfSumInsured:           "",
		SelfInsuranceDesc:          "",
		Description:                "Per danni alle cose di terzi sulle quali si eseguono i lavori",
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
		Group:                      "RCO",
		IsExtension:                false,
		CompanyCodec:               "CO",
		CompanyName:                "Resp. Civile verso Prestatori di Lavoro",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		SelfInsuranceDesc:          "",
		Description:                "Responsabilità per: la rivalsa INAIL per gli infortuni sul lavoro subiti dai prestatori di lavoro;  morte; e lesioni personali dalle quali sia derivata un'invalidità permanente ai sensi del codice civile, incluse le malattie professionali",
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["product-liability"] = &Coverage{
		Slug:                       "product-liability",
		Type:                       "company",
		TypeOfSumInsured:           "",
		Group:                      "RCT",
		IsExtension:                false,
		CompanyCodec:               "RP",
		CompanyName:                "RC Prodotti",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Danni involontariamente cagionati a terzi da difetto dei prodotti, venduti o distribuiti, per i quali rivesta in Italia la qualifica di produttore dopo la loro consegna a terzi, per danni a persone o cose",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-liability-construction-company"] = &Coverage{
		Slug:                       "third-party-liability-construction-company",
		Type:                       "company",
		Group:                      "RCT",
		IsExtension:                false,
		CompanyCodec:               "RE",
		CompanyName:                "RC Edile",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Pacchetto di garanzie che includei danni: 1) a condutture o impianti sotterranei; 2) da cedimento e franamento del terreno; 3) da bagnamento e spargimento d'acqua",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["legal-defence"] = &Coverage{
		Slug:              "legal-defence",
		Type:              "company",
		CompanyCodec:      "DP",
		Group:             "LEGAL",
		IsExtension:       false,
		CompanyName:       "Tutela Legale",
		LegalDefence:      "basic",
		SelfInsuranceDesc: "",
		Description:       "Difesa penale per reati di natura colposa o contravvenzionale, inclusi i casi di sicurezza aziendale da D. Lgs. 81/08 e D. Lgs. 106/09, D. Lgs. 193/07, D. Lgs. 152/06, D. Lgs. 101/18, D. Lgs. 231/01",
		Tax:               21.25,
		IsBase:            false,
		IsYuor:            false,
		IsPremium:         false,
	}
	res["cyber"] = &Coverage{
		Slug:                       "cyber",
		Type:                       "company",
		Group:                      "CYBER",
		IsExtension:                false,
		CompanyCodec:               "CR",
		CompanyName:                "Cyber Risk",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Indennizzo delle spese, a seguito di un attacco informatico, per: Ripristino dei dati; Violazione della privacy e violazione di dati confidenziali; Estorsione cyber; Cyber crime; Danno reputazionale; Danni su carte di pagamento/credito (PCI-DSS), nonché danni a terzi da: Violazioni della sicurezza della rete; Danni da interruzione di attività; Danni da responsabilità multimediale",
		Tax:                        0,
		Taxes:                      []Tax{{Tax: 22.25, Percentage: 40.0}, {Tax: 21.25, Percentage: 60.0}},
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["building"] = &Coverage{
		Slug:                       "building",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                false,
		TypeOfSumInsured:           "",
		CompanyCodec:               "IF",
		CompanyName:                "Incendio Fabbricato",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni diretti a Fabbricato, causati da eventi quali: incendio, esplosione, scoppio, fulmine, conseguenti fumi, gas e vapori",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0.0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["content"] = &Coverage{
		Slug:                       "content",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                false,
		CompanyCodec:               "IC",
		CompanyName:                "Incendio Contenuto",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni diretti a Contenuto (Merci, macchinari, attrezzature, arredamento), causati da eventi quali: incendio, esplosione, scoppio, fulmine, conseguenti fumi, gas e vapori",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["lease-holders-interest"] = &Coverage{
		Slug:                       "lease-holders-interest",
		Type:                       "building",
		Group:                      "RL",
		IsExtension:                false,
		CompanyCodec:               "RL",
		CompanyName:                "Rischio Locativo",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Danni materiali e diretti, cagionati ai locali tenuti in locazione, da Incendio, esplosio, scoppio, fumo nei casi di responsabilità ai termini degli artt. 1588, 1589 e 1611 del Codice Civile",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["burst-pipe"] = &Coverage{
		Slug:                       "burst-pipe",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "AC",
		CompanyName:                "Danni d’Acqua",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Danni da allagamento (eccesso o accumulo d’acqua in luogo normalmente asciutto) verificatosi all'interno del Fabbricato a seguito di formazione di ruscelli o accumulo esterno di acqua",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["power-surge"] = &Coverage{
		Slug:                       "power-surge",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "FE",
		CompanyName:                "Fenomeno Elettrico",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni originati da scariche, correnti, corto circuito ed altri fenomeni elettrici",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["atmospheric-event"] = &Coverage{
		Slug:                       "atmospheric-event",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "EA",
		CompanyName:                "Eventi Atmosferici",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da eventi atmosferici quali uragano, bufera, tempesta, grandine, vento e cose trascinate da esso, tromba d’aria, gelo, sovraccarico di neve",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["sociopolitical-event"] = &Coverage{
		Slug:                       "sociopolitical-event",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "ES",
		CompanyName:                "Eventi Sociopolitici",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da eventi sociopolitici, quali tumulti popolari, scioperi, sommosse, atti vandalici o dolosi e atti di sabotaggio",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["terrorism"] = &Coverage{
		Slug:                       "terrorism",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "AT",
		CompanyName:                "Atti di Terrorismo",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da atti di terrorismo",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["earthquake"] = &Coverage{
		Slug:                       "earthquake",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "TR",
		CompanyName:                "Terremoto",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da terremoto ai beni assicurati  compresi quelli di Incendio, Esplosione, Scoppio",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["river-flood"] = &Coverage{
		Slug:                       "river-flood",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "AL",
		CompanyName:                "Alluvione/Inondazione",
		TypeOfSumInsured:           "replacementValue",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da alluvione/inondazione (allagamento di un territorio causato da straripamento, esondazione, tracimazione o fuoriuscita dagli argini di corsi d’acqua, da laghi e bacini, sia naturali sia artificiali, anche se derivanti da eventi atmosferici.",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["water-damage"] = &Coverage{
		Slug:                       "water-damage",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "AG",
		CompanyName:                "Allagamento",
		TypeOfSumInsured:           "replacementValue",
		SelfInsurance:              "0",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni da fuoriuscita di acqua condotta (ingorghi, trabocchi, rotture accedentali), e, se assicurato il fabbricato, ricerca ripristino e riparazione del danno",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["glass"] = &Coverage{
		Slug:                       "glass",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "RO",
		CompanyName:                "Rottura Lastre",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Spese sostenute per la sostituzione di Lastre e insegne con altre nuove eguali o equivalenti per caratteristiche, compresi i costi di trasporto ed installazione, la cui rottura sia avvenuta per cause Accidentali o imputabili a fatti di terzi.",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["machinery-breakdown"] = &Coverage{
		Slug:                       "machinery-breakdown",
		Type:                       "building",
		Group:                      "FIRE",
		IsExtension:                true,
		CompanyCodec:               "GU",
		CompanyName:                "Guasti Macchine",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni materiali e diretti, al Macchinario causati o dovuti a guasti Accidentali meccanici in genere.",
		SumInsuredLimitOfIndemnity: 0,
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-recourse"] = &Coverage{
		Slug:                       "third-party-recourse",
		Type:                       "building",
		Group:                      "RT",
		IsExtension:                false,
		CompanyCodec:               "RT",
		CompanyName:                "Ricorso Terzi da Incendio",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Responsabilità per danni materiali e diretti arrecati alle cose di terzi in seguito a Incendio, Esplosione o Scoppio del Fabbricato e/o Contenuto, qualora assicurati, anche quando il Fabbricato lo è nella forma di Rischio Locativo",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["theft"] = &Coverage{
		Slug:                       "theft",
		Type:                       "building",
		Group:                      "THEFT",
		IsExtension:                false,
		CompanyCodec:               "FU",
		CompanyName:                "Garanzia Furto, Rapina ed Estorsione",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni subiti da furto, rapina o estorsione, inclusi i guasti e gli atti vandalici commessi dai ladri",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables-in-safe-strongrooms"] = &Coverage{
		Slug:                       "valuables-in-safe-strongrooms",
		Type:                       "building",
		Group:                      "THEFT",
		IsExtension:                true,
		CompanyCodec:               "FV",
		CompanyName:                "Furto Valori e Preziosi in Cassaforte",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Furto di valori e preziosi inerenti l’Attività assicurata riposti in cassforte a muro o di peso non inferiore a 200Kg",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables"] = &Coverage{
		Slug:                       "valuables",
		Type:                       "building",
		Group:                      "THEFT",
		IsExtension:                true,
		CompanyCodec:               "FP",
		CompanyName:                "Portavalori",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Perdita di valori a seguito di Furto o di Rapina commessi nel corso del loro trasporto al di fuori dei locali dell’Attività, compiuti nei confronti dell’Assicurato, dei soci o familiari coadiuvanti o dei prestatori di lavoro.",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["electronic-equipment"] = &Coverage{
		Slug:                       "electronic-equipment",
		Type:                       "building",
		Group:                      "ELETRONIC",
		IsExtension:                false,
		CompanyCodec:               "EL",
		CompanyName:                "Garanzia Apparecchiature Elettroniche",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Danni materiali e diretti ad apparecchiature elettroniche fisse e ad impiego mobile, causati da qualsiasi evento accidentale",
		Tax:                        21.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["increased-cost-of-working"] = &Coverage{
		Slug:                       "increased-cost-of-working",
		Type:                       "building",
		Group:                      "ELETRONIC",
		IsExtension:                true,
		CompanyCodec:               "MG",
		CompanyName:                "Maggiori Costi",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "In caso di sinistro indennizzabile con garanzia Elettronica, che provochi l’interruzione parziale o totale del funzionamento dei beni assicurati, si indennizzano i maggiori costi necessari alla prosecuzione delle funzioni svolte dall’apparecchio danneggiato o distrutto",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["restoration-of-data"] = &Coverage{
		Slug:                       "restoration-of-data",
		Type:                       "building",
		Group:                      "ELETRONIC",
		IsExtension:                true,
		CompanyCodec:               "SD",
		CompanyName:                "",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "In caso di danno indennizzabile, i costi necessari ed effettivamente sostenuti per il riacquisto dei supporti di Dati distrutti, danneggiati o sottratti.",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		Type:                       "building",
		Group:                      "ELETRONIC",
		IsExtension:                true,
		CompanyCodec:               "LU",
		CompanyName:                "Programmi in licenza d'uso",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "In caso di danno indennizzabile, i costi necessari ed effettivamente sostenuti per la duplicazione o per il riacquisto dei programmi in licenza d'uso distrutti, danneggiati o sottratti",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["business-interruption"] = &Coverage{
		Slug:                       "business-interruption",
		Type:                       "building",
		Group:                      "BUSINESS INTERRUPTTION",
		IsExtension:                false,
		CompanyCodec:               "BI",
		CompanyName:                "Business Interruption",
		TypeOfSumInsured:           "firstLoss",
		Deductible:                 "0",
		SelfInsuranceDesc:          "",
		Description:                "Indennizzo per il periodo di documentata inattività forzata, a seguito di un Sinistro avvenuto nel Fabbricato, che abbia danneggiato i locali e/o i Macchinar e/o le Apparecchiature Elettroniche funzionali all’attività",
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
		Group:                      "RCF",
		IsExtension:                false,
		CompanyCodec:               "RF",
		CompanyName:                "RC Fabbricato",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		SelfInsuranceDesc:          "",
		Description:                "Danni a persone o cose, verificatosi in relazione alla proprietà e conduzione del Fabbricato e delle eventuali parti comuni a esso riferite/collegate",
		Tax:                        22.25,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}

	res["environmental-liability"] = &Coverage{
		Slug:                       "environmental-liability",
		Type:                       "building",
		Group:                      "RCI",
		IsExtension:                false,
		CompanyCodec:               "RI",
		CompanyName:                "RC Inquinamento",
		TypeOfSumInsured:           "",
		Deductible:                 "0",
		SelfInsurance:              "0",
		SelfInsuranceDesc:          "",
		Description:                "Danni a persone o cose in conseguenza di contaminazione dell’acqua, dell’aria o del suolo, provocati dalla fuoriuscita di sostanze di qualunque natura a seguito di fatto improvviso, imprevedibile e dovuto a rottura accidentale di impianti, macchinari e condutture",
		Tax:                        22.25,
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["assistance"] = &Coverage{
		Slug:              "assistance",
		Type:              "building",
		Group:             "ASSSISTANCE",
		IsExtension:       false,
		CompanyCodec:      "AS",
		CompanyName:       "Assistenza al Fabbricato",
		Assistance:        "yes",
		SelfInsuranceDesc: "",
		Description:       "Prestazioni di assistenza e servizio 24/7 al Fabbricato quali invio di artigiani come: idraulico, elettricista, fabbro, serrandista, vetraio, sorvegliante, nei casi indicati in polizza di necessità (la compagnia eroga direttamente la prestazione non il rimborso delle spese)",
		Tax:               10.00,
		IsBase:            false,
		IsYuor:            false,
		IsPremium:         false,
	}
	return res
}
