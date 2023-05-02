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
			/*{
				Route:   "/risk/person",
				Handler: Person,
				Method:  http.MethodPost,
			},
			{
				Route:   "/survey/person",
				Handler: PersonSurvey,
				Method:  http.MethodPost,
			},
			{
				Route:   "/risk/pmi",
				Handler: PmiAllrisk,
				Method:  http.MethodPost,
			},
			{
				Route:   "/sales/life",
				Handler: Life,
				Method:  http.MethodPost,
			},*/
			{
				Route:   "/v1/survey/life",
				Handler: LifeSurvey,
				Method:  http.MethodPost,
			},
			/*{
				Route:   "/v1/statements/life",
				Handler: LifeStatements,
				Method:  http.MethodPost,
			},*/
		},
	}
	route.Router(w, r)

}

type Surveys struct {
	Surveys []*models.Survey `json:"surveys"`
	Text    string           `json:"text,omitempty"`
}
