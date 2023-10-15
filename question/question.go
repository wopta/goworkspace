package question

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

const (
	statements = "statements"
	surveys    = "surveys"
)

func init() {
	log.Println("INIT Question")

	functions.HTTP("Question", Question)
}

func Question(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/:questionType",
				Handler: GetQuestionsFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v2/:questionType",
				Handler: GetQuestionsV2Fx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}

func GetQuestionsV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result interface{}
		policy *models.Policy
	)

	log.Println("[GetQuestionsV2Fx] handler start -------------------")

	questionType := r.Header.Get("questionType")
	log.Printf("[GetQuestionsV2Fx] questions: %s", questionType)

	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[GetQuestionsV2Fx] req body: %s", string(body))

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[GetQuestionsV2] error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	switch questionType {
	case statements:
		log.Printf("[GetQuestionsV2Fx] loading statements for %s product", policy.Name)
		result, err = GetStatementsV2(policy)
	case surveys:
		log.Printf("[GetQuestionV2Fx] loading surveys for %s product", policy.Name)
		//out = GetSurveysV2(policy)
	default:
		log.Printf("[GetQuestionV2Fx] questionType %s not allowed", questionType)
		return "", nil, fmt.Errorf("questionType %s not allowed", questionType)
	}

	if err != nil {
		log.Printf("[GetQuestionsV2Fx] error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := json.Marshal(result)

	log.Println("[GetQuestionsV2Fx] handler end -------------------")

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, result, err
}

func loadExternalData(productName string) []byte {
	var data []byte
	switch productName {
	case models.PersonaProduct:
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
