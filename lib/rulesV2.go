package lib

import (
	"encoding/json"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"log"
)

func RulesFromJsonV2(fx interface{}, groule []byte, out interface{}, in []byte, data []byte) (string, interface{}) {
	log.Println("RulesFromJson")

	var err error
	// create new instance of DataContext
	dataContext := ast.NewDataContext()
	// add your JSON Fact into data context using AddJSON() function.
	if in != nil {
		err = dataContext.AddJSON("in", in)
		log.Println("RulesFromJson in")
		CheckError(err)
	}

	if out != nil {
		err = dataContext.Add("out", out)
		log.Println("RulesFromJson out")
		CheckError(err)
	}

	if data != nil {
		err = dataContext.AddJSON("data", data)
		log.Println("RulesFromJson data loaded")
		CheckError(err)
	}

	err = dataContext.Add("fx", fx)
	log.Println("RulesFromJson fx loaded")
	CheckError(err)

	underlying := pkg.NewBytesResource(groule)
	CheckError(err)

	resource := pkg.NewJSONResourceFromResource(underlying)
	CheckError(err)
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//bs := pkg.NewBytesResource([]byte(fileRes))

	err = ruleBuilder.BuildRuleFromResource("rules", "0.0.1", resource)
	CheckError(err)
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("rules", "0.0.1")
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataContext, knowledgeBase)
	CheckError(err)

	//resp := "execute"
	b, err := json.Marshal(out)
	CheckError(err)

	return string(b), out
}

func GetRulesFile(rulesFileName string) []byte {
	return GetFilesByEnv("grules/" + rulesFileName)
}
