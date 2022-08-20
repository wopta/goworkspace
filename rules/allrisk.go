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
	}
	if fil.Nrow() > 0 {

		var compare = "OK"
		if strings.Contains(fil.Elem(0, 4).String(), compare) {
			prof.Fire = compare
		} else {
			prof.Fire = fil.Elem(0, 4).String()
		}
		if strings.Contains(fil.Elem(0, 5).String(), compare) {
			prof.FireLow500k = compare
		} else {
			prof.FireLow500k = fil.Elem(0, 5).String()
		}
		if strings.Contains(fil.Elem(0, 6).String(), compare) {
			prof.FireUp500k = compare
		} else {
			prof.FireUp500k = fil.Elem(0, 6).String()
		}
		if strings.Contains(fil.Elem(0, 7).String(), compare) {
			prof.Theft = compare
		} else {
			prof.Theft = fil.Elem(0, 7).String()
		}
		if strings.Contains(fil.Elem(0, 8).String(), compare) {
			prof.ThefteLow500k = compare
		} else {
			prof.ThefteLow500k = fil.Elem(0, 8).String()
		}
		if strings.Contains(fil.Elem(0, 9).String(), compare) {
			prof.TheftUp500k = compare
		} else {
			prof.TheftUp500k = fil.Elem(0, 9).String()
		}
		if strings.Contains(fil.Elem(0, 10).String(), compare) {
			prof.Rct = compare
		} else {
			prof.Rct = fil.Elem(0, 10).String()
		}
		if strings.Contains(fil.Elem(0, 11).String(), compare) {
			prof.Rco = compare
		} else {
			prof.Rco = fil.Elem(0, 11).String()
		}
		if strings.Contains(fil.Elem(0, 12).String(), compare) {
			prof.RcoProd = compare
		} else {
			prof.RcoProd = fil.Elem(0, 12).String()
		}

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
	profileAllriskJson.Result = prof.Result
	resp := strings.ReplaceAll(prof.Result, "'", "\"")
	//resp, _ := json.Marshal(profileAllriskJson)
	fmt.Fprintf(w, "{"+resp+"}")
	fmt.Println(prof.Revenue)
	fmt.Println(prof.Result)

}

type ProfileAllrisk struct {
	Vat              int64
	SquareMeters     int64
	IsBuildingOwner  bool
	Revenue          int64
	Address          string
	Ateco            string
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
}
