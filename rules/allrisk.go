package rules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	log.Println("Allrisk")
	///log.Println(os.Getenv("SA_KEY"))
	ricAteco := lib.GetFromStorage("function-data", "data/rules/Riclassificazione_Ateco.csv", "")
	groule := lib.GetFromStorage("function-data", "grules/allrisk.grl", "")
	log.Println("GetFromStorage")
	var profileAllriskJson models.ProfileAllriskJson
	//var profileAllrisk ProfileAllrisk
	df := lib.CsvToDataframe(ricAteco)

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
	//fmt.Println("filtered", fil)

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

	//fmt.Println("profileAllrisk post copy:", prof)
	dataCtx := ast.NewDataContext()
	err = dataCtx.Add("In", prof)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	// lets prepare a rule definitionS

	var fileRes pkg.Resource

	if os.Getenv("env") == "dev" {
		fileRes = pkg.NewFileResource("rules/allrisk.grl")
	} else {
		fileRes = pkg.NewBytesResource(groule)
	}
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//bs := pkg.NewBytesResource([]byte(fileRes))
	err = ruleBuilder.BuildRuleFromResource("rules", "0.0.1", fileRes)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("rules", "0.0.1")
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	//resp := strings.ReplaceAll(prof.Result, "'", "\"")
	b, err := json.Marshal(prof.Coverages)
	log.Println(string(b))

	m := make([]*Coverage, 0, len(prof.Coverages))
	for _, val := range prof.Coverages {
		m = append(m, val)
		fmt.Println(val)
	}
	resp, _ := json.Marshal(m)
	fmt.Fprintf(w, string(resp))
	//fmt.Println(string(resp))

	//fmt.Println(prof.Coverages)

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
