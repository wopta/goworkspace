package question

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-chi/chi"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

const (
	statements = "statements"
	surveys    = "surveys"
)

var questionRoutes []lib.Route = []lib.Route{
	{
		Route:   "/v1/{questionType}",
		Handler: lib.ResponseLoggerWrapper(GetQuestionsFx),
		Method:  http.MethodPost,
		Roles:   []string{lib.UserRoleAll},
	},
}

func init() {
	log.Println("INIT Question")
	functions.HTTP("Question", Question)
}

func Question(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	router := lib.GetRouter("question", questionRoutes)
	router.ServeHTTP(w, r)
}

func GetQuestionsFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result interface{}
		policy *models.Policy
	)

	log.SetPrefix("[GetQuestionsFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	questionType := chi.URLParam(r, "questionType")
	log.Printf("questions: %s", questionType)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("error unmarshaling request body: %s", err.Error())
		return "", nil, err
	}

	policy.Normalize()

	switch questionType {
	case statements:
		log.Printf("loading statements for %s product", policy.Name)
		result, err = GetStatements(policy)
	case surveys:
		log.Printf("loading surveys for %s product", policy.Name)
		result, err = GetSurveys(policy)
	default:
		log.Printf("questionType %s not allowed", questionType)
		return "", nil, fmt.Errorf("questionType %s not allowed", questionType)
	}

	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", nil, err
	}

	jsonOut, err := json.Marshal(result)

	log.Println("handler end -------------------------------------------------")

	return `{"` + questionType + `": ` + string(jsonOut) + `}`, result, err
}

func loadExternalData(productName, productVersion string) []byte {
	var data []byte

	log.Println("[loadExternalData] function start ------------")

	switch productName {
	case models.PersonaProduct:
		data = lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/statements_coherence_data.json", models.ProductsFolder,
			productName, productVersion))
	}

	log.Printf("[loadExternalData] response: %s", string(data))

	log.Println("[loadExternalData] function end ------------")

	return data
}
