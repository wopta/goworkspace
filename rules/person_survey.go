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
		groule []byte
		policy models.Policy
	)
	const (
		rulesFileName = "person_survey.json"
	)

	log.Println("Person Survey")

	b, err := io.ReadAll(r.Body)
	lib.CheckError(err)
	err = json.Unmarshal(b, &policy)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	surveys := &Surveys{Surveys: make([]*models.Survey, 0), Text: ""}

	switch os.Getenv("env") {
	case "local":
		groule = lib.ErrorByte(os.ReadFile("../function-data/dev/grules/" + rulesFileName))
	case "dev":
		groule = lib.GetFromStorage("function-data", "grules/"+rulesFileName, "")
	case "prod":
		groule = lib.GetFromStorage("core-350507-function-data", "grules/"+rulesFileName, "")
	}

	_, ruleOutput := rulesFromJson(groule, surveys, policyJson, []byte(getCoherenceData()))

	ruleOutputJson, err := json.Marshal(ruleOutput)
	lib.CheckError(err)

	return string(ruleOutputJson), ruleOutput, nil
}

type Surveys struct {
	Surveys []*models.Survey `json:"surveys"`
	Text    string           `json:"text,omitempty"`
}

type FxSurvey struct{}

func (fx *FxSurvey) AppendStatement(statements []*models.Statement, title string, subtitle string, hasMultipleAnswers bool, hasAnswer bool, expectedAnswer bool) []*models.Statement {
	statement := &models.Statement{
		Title:              title,
		Subtitle:           subtitle,
		HasMultipleAnswers: nil,
		Questions:          make([]*models.Question, 0),
		Answer:             nil,
		HasAnswer:          hasAnswer,
		ExpectedAnswer:     nil,
	}
	if hasAnswer {
		statement.ExpectedAnswer = &expectedAnswer
	}
	if hasMultipleAnswers {
		statement.HasMultipleAnswers = &hasMultipleAnswers
	}
	return append(statements, statement)
}

func (fx *FxSurvey) AppendSurvey(surveys []*models.Survey, title string, subtitle string, hasMultipleAnswers bool, hasAnswer bool, expectedAnswer bool) []*models.Survey {
	survey := &models.Survey{
		Title:              title,
		Subtitle:           subtitle,
		HasMultipleAnswers: nil,
		Questions:          make([]*models.Question, 0),
		Answer:             nil,
		HasAnswer:          hasAnswer,
		ExpectedAnswer:     nil,
	}
	if hasAnswer {
		survey.ExpectedAnswer = &expectedAnswer
	}
	if hasMultipleAnswers {
		survey.HasMultipleAnswers = &hasMultipleAnswers
	}
	return append(surveys, survey)
}

func (fx *FxSurvey) AppendQuestion(questions []*models.Question, text string, isBold bool, indent bool, hasAnswer bool, expectedAnswer bool) []*models.Question {
	question := &models.Question{
		Question:       text,
		IsBold:         isBold,
		Indent:         indent,
		Answer:         nil,
		HasAnswer:      hasAnswer,
		ExpectedAnswer: nil,
	}
	if hasAnswer {
		question.ExpectedAnswer = &expectedAnswer
	}

	return append(questions, question)
}

func (fx *FxSurvey) HasGuaranteePolicy(input map[string]interface{}, guaranteeSlug string) bool {
	j, err := json.Marshal(input)
	lib.CheckError(err)
	var policy models.Policy
	err = json.Unmarshal(j, &policy)
	lib.CheckError(err)
	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == guaranteeSlug {
				return true
			}
		}
	}
	return false
}

func (fx *FxSurvey) GetGuaranteeIndex(input map[string]interface{}, guranteeSlug string) int {
	j, _ := json.Marshal(input)
	var policy models.Policy
	_ = json.Unmarshal(j, &policy)
	for _, asset := range policy.Assets {
		for i, guarantee := range asset.Guarantees {
			if guarantee.Slug == guranteeSlug {
				return i
			}
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
