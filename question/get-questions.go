package question

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

func GetQuestionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		out    interface{}
		policy models.Policy
	)
	log.Println("[GetQuestionsFx]")

	questionType := r.Header.Get("questionType")
	log.Println("[GetQuestionFx] questionType " + questionType)

	body, err := io.ReadAll(r.Body)
	lib.CheckError(err)
	err = json.Unmarshal(body, &policy)

	switch questionType {
	case "statements":
		log.Printf("[GetQuestionFx] loading statements for %s product", policy.Name)
		out = GetStatements(policy)
	case "surveys":
		log.Printf("[GetQuestionFx] loading surveys for %s product", policy.Name)
		out = GetSurveys(policy)
	}

	jsonOut, err := json.Marshal(out)

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, out, err
}

func GetStatements(policy models.Policy) []models.Statement {
	const (
		rulesFilenameSuffix = "_statements.json"
	)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	fx := new(models.Fx)
	statements := &Statements{
		Statements: make([]*models.Statement, 0),
		Text:       "",
	}

	rulesFile := lib.GetRulesFile(policy.Name + rulesFilenameSuffix)
	data := loadExternalData(policy.Name)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, statements, policyJson, data)

	out := make([]models.Statement, 0)
	for _, statement := range ruleOutput.(*Statements).Statements {
		out = append(out, *statement)
	}

	return out
}

func GetSurveys(policy models.Policy) []models.Survey {
	const (
		rulesFilenameSuffix = "_survey.json"
	)

	policyJson, err := policy.Marshal()
	lib.CheckError(err)

	fx := new(models.Fx)
	surveys := &Surveys{
		Surveys: make([]*models.Survey, 0),
		Text:    "",
	}

	rulesFile := lib.GetRulesFile(policy.Name + rulesFilenameSuffix)
	data := loadExternalData(policy.Name)

	_, ruleOutput := lib.RulesFromJsonV2(fx, rulesFile, surveys, policyJson, data)

	out := make([]models.Survey, 0)
	for _, survey := range ruleOutput.(*Surveys).Surveys {
		out = append(out, *survey)
	}

	return out
}

func loadExternalData(productName string) []byte {
	var data []byte
	switch productName {
	case "persona":
		data = []byte(getCoherenceData())
	}
	return data
}

func getCoherenceData() string {
	return `{
		"AA": {
			"extra": "nel tempo libero,",
			"professionale": "nello svolgimento dell'attivtà lavorativa,",
			"24h": "sia al lavoro che nel tempo libero,"
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
