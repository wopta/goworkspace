package rules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func Allrisk(w http.ResponseWriter, r *http.Request) {
	//lib.Files("")

	var profileAllriskJson models.ProfileAllriskJson
	//var profileAllrisk ProfileAllrisk
	df := lib.CsvToDataframe("data/Riclassificazione_Ateco.csv")

	err := json.NewDecoder(r.Body).Decode(&profileAllriskJson)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fil := df.Filter(
		dataframe.F{Colidx: 5, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: profileAllriskJson.Ateco},
	)
	fmt.Println("filtered", fil.Nrow())
	fmt.Println("filtered", fil.Ncol())
	fmt.Println("filtered", fil)

	prof := &ProfileAllrisk{
		Vat:              profileAllriskJson.Vat,
		SquareMeters:     profileAllriskJson.SquareMeters,
		IsBuildingOwner:  profileAllriskJson.IsBuildingOwner,
		Revenue:          profileAllriskJson.Revenue,
		Address:          profileAllriskJson.Address,
		Ateco:            profileAllriskJson.Ateco,
		BusinessSector:   profileAllriskJson.BusinessSector,
		BuildingType:     profileAllriskJson.BuildingType,
		BuildingMaterial: profileAllriskJson.BuildingMaterial,
		BuildingYear:     profileAllriskJson.BuildingYear,
		Employer:         profileAllriskJson.Employer,
		IsAllarm:         profileAllriskJson.IsAllarm,
		Floor:            profileAllriskJson.Floor,
		IsPRA:            profileAllriskJson.IsPRA,
		Costruction:      profileAllriskJson.Costruction,
		IsHolder:         profileAllriskJson.IsHolder,
		Result:           profileAllriskJson.Result,
		Coverages:        initCoverage(),
	}
	if fil.Nrow() > 0 {
		prof.AtecoMacro = strings.ToUpper(fil.Elem(0, 0).String())
		prof.AtecoSub = strings.ToUpper(fil.Elem(0, 1).String())
		prof.AtecoDesc = strings.ToUpper(fil.Elem(0, 2).String())
		prof.BusinessSector = strings.ToUpper(fil.Elem(0, 3).String())
		prof.Fire = strings.ToUpper(fil.Elem(0, 14).String())
		prof.FireLow500k = strings.ToUpper(fil.Elem(0, 5).String())
		prof.FireUp500k = strings.ToUpper(fil.Elem(0, 6).String())
		prof.Theft = strings.ToUpper(fil.Elem(0, 15).String())
		prof.ThefteLow500k = strings.ToUpper(fil.Elem(0, 8).String())
		prof.TheftUp500k = strings.ToUpper(fil.Elem(0, 9).String())
		prof.Rct = strings.ToUpper(fil.Elem(0, 16).String())
		prof.Rco = strings.ToUpper(fil.Elem(0, 17).String())
		prof.RcoProd = strings.ToUpper(fil.Elem(0, 18).String())
		prof.RcVehicle = strings.ToUpper(fil.Elem(0, 19).String())
		prof.Rcpo = strings.ToUpper(fil.Elem(0, 20).String())
		prof.Rcp12 = strings.ToUpper(strings.ToUpper(fil.Elem(0, 21).String()))
		prof.Rcp2008 = strings.ToUpper(fil.Elem(0, 22).String())
		prof.DamageTheft = strings.ToUpper(fil.Elem(0, 23).String())
		prof.DamageThing = strings.ToUpper(fil.Elem(0, 24).String())
		prof.RcCostruction = strings.ToUpper(fil.Elem(0, 26).String())
		prof.Eletronic = strings.ToUpper(fil.Elem(0, 27).String())
		prof.MachineFaliure = strings.ToUpper(fil.Elem(0, 28).String())

	}
	//copier.CopyWithOption(&profileAllrisk, &profileAllriskJson, copier.Option{IgnoreEmpty: false, DeepCopy: true})

	fmt.Println("profileAllrisk post copy:", prof)
	dataCtx := ast.NewDataContext()
	err = dataCtx.Add("In", prof)
	if err != nil {
		panic(err)
	}
	// lets prepare a rule definitionS
	fileRes := pkg.NewFileResource("rules/allrisk.grl")
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//bs := pkg.NewBytesResource([]byte(fileRes))
	err = ruleBuilder.BuildRuleFromResource("rules", "0.0.1", fileRes)
	if err != nil {
		panic(err)
	}
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("rules", "0.0.1")
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	//resp := strings.ReplaceAll(prof.Result, "'", "\"")

	m := make([]*Coverage, 0, len(prof.Coverages))
	for _, val := range prof.Coverages {
		m = append(m, val)
		fmt.Println(val)
	}
	resp, _ := json.Marshal(m)
	fmt.Fprintf(w, string(resp))
	fmt.Println(string(resp))

	fmt.Println(prof.Coverages)

}

type ProfileAllrisk struct {
	Vat              int64
	SquareMeters     int64
	IsBuildingOwner  bool
	Revenue          int64
	Address          string
	Ateco            string
	AtecoMacro       string
	AtecoSub         string
	AtecoDesc        string
	BusinessSector   string
	BuildingType     string
	BuildingMaterial string
	BuildingYear     string
	Employer         int64
	IsAllarm         bool
	Floor            int64
	IsPRA            bool
	Costruction      string
	IsHolder         bool
	Result           string
	Fire             string
	FireLow500k      string
	FireUp500k       string
	Theft            string
	ThefteLow500k    string
	TheftUp500k      string
	Rct              string
	Rco              string
	RcoProd          string
	RcVehicle        string
	Rcpo             string
	Rcp12            string
	Rcp2008          string
	DamageTheft      string
	DamageThing      string
	RcCostruction    string
	Eletronic        string
	MachineFaliure   string
	Coverages        map[string]*Coverage
}
type Coverage struct {
	TypeOfSumInsured           string
	Deductible                 string
	SumInsuredLimitOfIndemnity int64
	Slug                       string
	IsBase                     bool
	IsYuor                     bool
	IsPremium                  bool
}

func initCoverage() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["third-party-liability"] = &Coverage{
		Slug:                       "third-party-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-in-custody"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-workmanships"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-12-months"] = &Coverage{
		Slug:                       "defect-liability-12-months",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["defect-liability-dm-37-2008"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-damage-due-to-theft"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["damage-to-goods-course-of-works"] = &Coverage{
		Slug:                       "damage-to-goods-course-of-works",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["employers-liability"] = &Coverage{
		Slug:                       "employers-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["product-liability"] = &Coverage{
		Slug:                       "product-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-liability-construction-company"] = &Coverage{
		Slug:                       "third-party-liability-construction-company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["legal-defence"] = &Coverage{
		Slug:                       "legal-defence",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["cyber"] = &Coverage{
		Slug:                       "cyber",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["building"] = &Coverage{
		Slug:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["content"] = &Coverage{
		Slug:                       "content",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["lease-holders-interest"] = &Coverage{
		Slug:                       "lease-holders-interest",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["burst-pipe"] = &Coverage{
		Slug:                       "burst-pipe",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["power-surge"] = &Coverage{
		Slug:                       "power-surge",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["atmospheric-event"] = &Coverage{
		Slug:                       "atmospheric-event",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["sociopolitical-event"] = &Coverage{
		Slug:                       "sociopolitical-event",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["terrorism"] = &Coverage{
		Slug:                       "terrorism",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["earthquake"] = &Coverage{
		Slug:                       "earthquake",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["river-flood"] = &Coverage{
		Slug:                       "river-flood",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["water-damage"] = &Coverage{
		Slug:                       "water-damage",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["glass"] = &Coverage{
		Slug:                       "glass",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["machinery-breakdown"] = &Coverage{
		Slug:                       "machinery-breakdown",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["third-party-recourse"] = &Coverage{
		Slug:                       "third-party-recourse",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["theft"] = &Coverage{
		Slug:                       "theft",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables-in-safe-strongrooms"] = &Coverage{
		Slug:                       "valuables-in-safe-strongrooms",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["valuables"] = &Coverage{
		Slug:                       "valuables",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["electronic-equipment"] = &Coverage{
		Slug:                       "electronic-equipment",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["increased-cost-of-working"] = &Coverage{
		Slug:                       "increased-cost-of-working",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["restoration-of-data"] = &Coverage{
		Slug:                       "restoration-of-data",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["business-interruption"] = &Coverage{
		Slug:                       "business-interruption",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["environmental-liability"] = &Coverage{
		Slug:                       "environmental-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	res["assistance"] = &Coverage{
		Slug:                       "assistance",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		IsBase:                     false,
		IsYuor:                     false,
		IsPremium:                  false,
	}
	return res
}
