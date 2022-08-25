package rules

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	df := lib.CsvToDataframe("data/Riclassificazione _Ateco.csv")

	err := json.NewDecoder(r.Body).Decode(&profileAllriskJson)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fil := df.Filter(
		dataframe.F{Colidx: 4, Colname: "Codice Ateco 2007", Comparator: series.Eq, Comparando: profileAllriskJson.Ateco},
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
		prof.AtecoMacro = fil.Elem(0, 0).String()
		prof.AtecoSub = fil.Elem(0, 1).String()
		prof.AtecoDesc = fil.Elem(0, 2).String()
		prof.Fire = fil.Elem(0, 4).String()
		prof.FireLow500k = fil.Elem(0, 5).String()
		prof.FireUp500k = fil.Elem(0, 6).String()
		prof.Theft = fil.Elem(0, 7).String()
		prof.ThefteLow500k = fil.Elem(0, 8).String()
		prof.TheftUp500k = fil.Elem(0, 9).String()
		prof.Rct = fil.Elem(0, 10).String()
		prof.Rco = fil.Elem(0, 11).String()
		prof.RcoProd = fil.Elem(0, 12).String()

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
	}
	resp, _ := json.Marshal(m)
	fmt.Fprintf(w, string(resp))
	fmt.Println(string(resp))
	fmt.Println(prof.Result)

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
	Coverages        map[string]*Coverage
}
type Coverage struct {
	TypeOfSumInsured           string
	Deductible                 string
	SumInsuredLimitOfIndemnity int64
	Slug                       string
	isBase                     bool
	isYuor                     bool
	isPremium                  bool
}

func initCoverage() map[string]*Coverage {

	var res = make(map[string]*Coverage)
	res["third-party-liability"] = &Coverage{
		Slug:                       "third-party-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["damage-to-goods-in-custody"] = &Coverage{
		Slug:                       "damage-to-goods-in-custody",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["defect-liability-workmanships"] = &Coverage{
		Slug:                       "defect-liability-workmanships",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["defect-liability-12-months"] = &Coverage{
		Slug:                       "defect-liability-12-months",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["defect-liability-dm-37-2008"] = &Coverage{
		Slug:                       "defect-liability-dm-37-2008",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["property-damage-due-to-theft"] = &Coverage{
		Slug:                       "property-damage-due-to-theft",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["damage-to-goods-course-of-works"] = &Coverage{
		Slug:                       "damage-to-goods-course-of-works",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["employers-liability"] = &Coverage{
		Slug:                       "employers-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["product-liability"] = &Coverage{
		Slug:                       "product-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["third-party-liability-construction-company"] = &Coverage{
		Slug:                       "third-party-liability-construction-company",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["legal-defence"] = &Coverage{
		Slug:                       "legal-defence",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["cyber"] = &Coverage{
		Slug:                       "cyber",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["building"] = &Coverage{
		Slug:                       "building",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["content"] = &Coverage{
		Slug:                       "content",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["lease-holders-interest"] = &Coverage{
		Slug:                       "lease-holders-interest",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["burst-pipe"] = &Coverage{
		Slug:                       "burst-pipe",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["power-surge"] = &Coverage{
		Slug:                       "power-surge",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["atmospheric-event"] = &Coverage{
		Slug:                       "atmospheric-event",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["sociopolitical-event"] = &Coverage{
		Slug:                       "sociopolitical-event",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["terrorism"] = &Coverage{
		Slug:                       "terrorism",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["earthquake"] = &Coverage{
		Slug:                       "earthquake",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["river-flood"] = &Coverage{
		Slug:                       "river-flood",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["water-damage"] = &Coverage{
		Slug:                       "water-damage",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["glass"] = &Coverage{
		Slug:                       "glass",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["machinery-breakdown"] = &Coverage{
		Slug:                       "machinery-breakdown",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["third-party-recourse"] = &Coverage{
		Slug:                       "third-party-recourse",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["theft"] = &Coverage{
		Slug:                       "theft",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["valuables-in-safe-strongrooms"] = &Coverage{
		Slug:                       "valuables-in-safe-strongrooms",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["valuables"] = &Coverage{
		Slug:                       "valuables",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["electronic-equipment"] = &Coverage{
		Slug:                       "electronic-equipment",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["increased-cost-of-working"] = &Coverage{
		Slug:                       "increased-cost-of-working",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["restoration-of-data"] = &Coverage{
		Slug:                       "restoration-of-data",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["business-interruption"] = &Coverage{
		Slug:                       "business-interruption",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["property-owners-liability"] = &Coverage{
		Slug:                       "property-owners-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["software-under-license"] = &Coverage{
		Slug:                       "software-under-license",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["environmental-liability"] = &Coverage{
		Slug:                       "environmental-liability",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	res["assistance"] = &Coverage{
		Slug:                       "assistance",
		TypeOfSumInsured:           "namedPerils",
		Deductible:                 "0",
		SumInsuredLimitOfIndemnity: 0,
		isBase:                     false,
		isYuor:                     false,
		isPremium:                  false,
	}
	return res
}
