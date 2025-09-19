package question

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

const (
	statements = "statements"
	surveys    = "surveys"
)

var questionRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/{questionType}",
		Handler: lib.ResponseLoggerWrapper(getQuestionsFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Question")
	functions.HTTP("Question", question)
}

func question(w http.ResponseWriter, r *http.Request) {
	router := lib.GetRouter("question", questionRoutes)
	router.ServeHTTP(w, r)
}

func getQuestionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result interface{}
		policy *models.Policy
	)

	log.AddPrefix("GetQuestionsFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	questionType := chi.URLParam(r, "questionType")
	log.Printf("questions: %s", questionType)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.ErrorF("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	switch questionType {
	case statements:
		log.Printf("loading statements for %s product", policy.Name)
		result, err = GetStatements(policy, true)
	case surveys:
		log.Printf("loading surveys for %s product", policy.Name)
		result, err = GetSurveys(policy)
	default:
		log.Printf("questionType %s not allowed", questionType)
		return "", nil, fmt.Errorf("questionType %s not allowed", questionType)
	}

	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := json.Marshal(result)

	log.Println("handler end -------------------------------------------------")

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, result, err
}

func loadExternalData(productName, productVersion string) []byte {
	var data []byte
	log.AddPrefix("LoadExternalData")
	defer log.PopPrefix()
	log.Println("function start ------------")

	switch productName {
	case models.PersonaProduct:
		data = lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/statements_coherence_data.json", models.ProductsFolder,
			productName, productVersion))
	}

	log.Printf("response: %s", string(data))

	log.Println("function end ------------")

	return data
}
