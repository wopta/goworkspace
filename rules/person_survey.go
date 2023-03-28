package rules

import (
	"encoding/json"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
)

func PersonSurvey(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy    models.Policy
		groule    []byte
		e         error
		questions []*Statement
	)
	const (
		rulesFileName = "person_survey.json"
	)

	log.Println("Person Survey")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	policy, e = models.UnmarshalPolicy(req)
	lib.CheckError(e)

	statements := &Statements{Statements: questions}
	dynamicTitle := &DynamicTitle{Text: ""}

	switch os.Getenv("env") {
	case "local":
		groule = lib.ErrorByte(os.ReadFile("../function-data/dev/grules/" + rulesFileName))
	case "dev":
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFileName, "")
	case "prod":
		groule = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFileName, "")
	}

	fx := &Fx{}
	fxSurvey := &FxSurvey{}
	// create new instance of DataContext
	dataContext := ast.NewDataContext()
	// add your JSON Fact into data context using AddJSON() function.
	err := dataContext.Add("in", policy)
	log.Println("RulesFromJson in")
	lib.CheckError(err)

	err = dataContext.Add("title", dynamicTitle)
	log.Println("RulesFromJson title")
	lib.CheckError(err)

	err = dataContext.Add("out", statements)
	log.Println("RulesFromJson out")
	lib.CheckError(err)

	err = dataContext.AddJSON("data", []byte(getCoherenceData()))
	log.Println("RulesFromJson data loaded")
	lib.CheckError(err)

	err = dataContext.Add("fx", fx)
	log.Println("RulesFromJson fx loaded")
	lib.CheckError(err)

	err = dataContext.Add("fxSurvey", fxSurvey)
	log.Println("RulesFromJson fx loaded")
	lib.CheckError(err)

	underlying := pkg.NewBytesResource(groule)
	resource := pkg.NewJSONResourceFromResource(underlying)
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//bs := pkg.NewBytesResource([]byte(fileRes))

	err = ruleBuilder.BuildRuleFromResource("rules", "0.0.1", resource)
	lib.CheckError(err)
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("rules", "0.0.1")
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataContext, knowledgeBase)
	lib.CheckError(err)

	b, err := json.Marshal(statements)
	lib.CheckError(err)

	return string(b), statements, nil
}

type Statements struct {
	Statements []*Statement `json:"statements"`
}

type Statement struct {
	Title              string      `json:"title"`
	HasMultipleAnswers *bool       `json:"hasMultipleAnswers,omitempty"`
	Questions          []*Question `json:"questions"`
	Answer             *bool       `json:"answer,omitempty"`
}

type Question struct {
	Question string `json:"question"`
	IsBold   bool   `json:"isBold"`
	Indent   bool   `json:"indent"`
	Answer   *bool  `json:"answer,omitempty"`
}

type DynamicTitle struct {
	Text string
}

type FxSurvey struct{}

func (fx *FxSurvey) AppendStatement(statements []*Statement, title string, hasMultipleAnswers bool, answer bool) []*Statement {
	statement := &Statement{
		Title:              title,
		HasMultipleAnswers: nil,
		Questions:          make([]*Question, 0),
		Answer:             nil,
	}
	if answer {
		statement.Answer = &answer
	}
	if hasMultipleAnswers {
		statement.HasMultipleAnswers = &hasMultipleAnswers
	}

	return append(statements, statement)
}

func (fx *FxSurvey) AppendQuestion(questions []*Question, text string, isBold bool, indent bool, answer bool) []*Question {
	question := &Question{
		Question: text,
		IsBold:   isBold,
		Indent:   indent,
		Answer:   nil,
	}
	if answer {
		question.Answer = &answer
	}
	return append(questions, question)
}

func (fx *FxSurvey) HasGuaranteePolicy(policy models.Policy, guaranteeName string) bool {
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Name == guaranteeName {
			return true
		}
	}
	return false
}

func (fx *FxSurvey) GetGuaranteeIndex(policy models.Policy, guaranteeName string) int {
	for i, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Name == guaranteeName {
			return i
		}
	}
	return -1
}

func getCoherenceData() string {
	return `{
		"AA": {
			"extra": "nel tempo libero,",
			"professionale": "nello svolgimento dell'attivtà lavorativa,",
			"24ore": "sia al lavoro che nel tempo libero,"
		},
		"BB": {
			"dipendente": "lavoratore dipendente",
			"autonomo": "lavoratore autonomo",
			"non lavoratore": "non lavoratore"
		},
		"CC1": "un capitale all'Assicurato, in caso di Invalidità Permanente,",
		"FR": {
			"3": "oltre il 3% di invalidità,",
			"5": "oltre il 5% di invalidità,",
			"10": "oltre il 10% di invalidità,"
		},
		"CC2": "a copertura della minore capacità di reddito;",
		"DD": "un capitale in caso di Decesso, avente finalità previdenziale, a copertura delle minori disponibilità che risulterebbero, in seguito al decesso dell'Assicurato, a favore dei beneficiari designati;",
		"EE": "una Diaria (importo giornaliero) in caso di ricovero o gessatura;",
		"FF": "una Diaria in caso di convalescenza post ricovero;",
		"GG": "un indennizzo per ogni giorno in cui l'Assicurato è inabile, in tutto o in parte, allo svolgimento delle proprie attività lavorative;",
		"HH": "inoltre l'assicurazione, sempre in caso di infortunio, risponde al bisogno di",
		"II": "rimborsare le spese mediche sostenute;",
		"JJ": "assicurare la difesa legale, per fatti illeciti di terzi o malpractice sanitaria;",
		"KK": "aiutare l'Assicurato con servizi di assistenza utili in momenti di bisogno (es. invio di un medico, consulto medico, trasporto in ambulanza, …);",
		"LL": "in caso di malattia, infine, consente di disporre di un capitale, a integrazione del reddito, qualora all'Assicurato derivi dalla malattia stessa una riduzione della capacità lavorativa (invalidità permanente) oltre il 24%"
	}`
}
