package question

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"log"
	"net/http"
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

type Statements struct {
	Statements []*models.Statement
	Text       string
}

type Surveys struct {
	Surveys []*models.Survey
	Text    string
}
