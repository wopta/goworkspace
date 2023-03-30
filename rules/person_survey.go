package rules

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
	"os"
)

func PersonSurvey(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		//policy    models.Policy
		groule []byte
		//e         error
		questions []*models.Statement
	)
	const (
		rulesFileName = "person_survey.json"
	)

	log.Println("Person Survey")
	policyJson := lib.ErrorByte(io.ReadAll(r.Body))

	statements := &Statements{Statements: questions, Text: nil}
	//dynamicTitle := &DynamicTitle{Text: ""}

	switch os.Getenv("env") {
	case "local":
		groule = lib.ErrorByte(os.ReadFile("../function-data/dev/grules/" + rulesFileName))
	case "dev":
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFileName, "")
	case "prod":
		groule = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFileName, "")
	}

	_, statementsOut := rulesFromJson(groule, statements, policyJson, []byte(getCoherenceData()))

	b, err := json.Marshal(statementsOut.(models.Policy))

	return string(b), statementsOut, err
}

type Statements struct {
	Statements []*models.Statement `json:"statements"`
	Text       *string             `json:"text,omitempty"`
}

type DynamicTitle struct {
	Text string
}

type FxSurvey struct{}

func (fx *FxSurvey) AppendStatement(statements []*models.Statement, title string, hasMultipleAnswers bool, answer bool) []*models.Statement {
	statement := &models.Statement{
		Title:              title,
		HasMultipleAnswers: nil,
		Questions:          make([]*models.Question, 0),
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

func (fx *FxSurvey) AppendQuestion(questions []*models.Question, text string, isBold bool, indent bool, answer bool) []*models.Question {
	question := &models.Question{
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

func (fx *FxSurvey) HasGuaranteePolicy(input []byte, guaranteeName string) bool {
	var policy models.Policy
	err := json.Unmarshal(input, &policy)
	lib.CheckError(err)
	for _, guarantee := range policy.Assets[0].Guarantees {
		if guarantee.Name == guaranteeName {
			return true
		}
	}
	return false
}

func (fx *FxSurvey) GetGuaranteeIndex(input []byte, guaranteeName string) int {
	var policy models.Policy
	err := json.Unmarshal(input, &policy)
	lib.CheckError(err)
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
