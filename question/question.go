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
		},
	}
	route.Router(w, r)

}

func GetQuestionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result interface{}
		policy *models.Policy
	)

	log.Println("[GetQuestionsFx] handler start -------------------")

	questionType := r.Header.Get("questionType")
	log.Printf("[GetQuestionsFx] questions: %s", questionType)

	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[GetQuestionsFx] req body: %s", string(body))

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[GetQuestionsFx] error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	switch questionType {
	case statements:
		log.Printf("[GetQuestionsFx] loading statements for %s product", policy.Name)
		result, err = GetStatements(policy)
	case surveys:
		log.Printf("[GetQuestionsFx] loading surveys for %s product", policy.Name)
		result, err = GetSurveys(policy)
	default:
		log.Printf("[GetQuestionsFx] questionType %s not allowed", questionType)
		return "", nil, fmt.Errorf("questionType %s not allowed", questionType)
	}

	if err != nil {
		log.Printf("[GetQuestionsFx] error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := json.Marshal(result)

	log.Printf("[GetQuestionsFx] response: %s", string(jsonOut))

	log.Println("[GetQuestionsFx] handler end -------------------")

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, result, err
}

func loadExternalData(productName, productVersion string) []byte {
	var data []byte

	log.Println("[loadExternalData] function start ------------")

	switch productName {
	case models.PersonaProduct:
		data = lib.GetFilesByEnv(fmt.Sprintf("products-v2/%s/%s/statements_coherence_data.json", productName, productVersion))
	}

	log.Printf("[loadExternalData] response: %s", string(data))

	log.Println("[loadExternalData] function end ------------")

	return data
}
