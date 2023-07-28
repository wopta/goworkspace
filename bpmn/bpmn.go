package bpmn

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
	//lib "github.com/wopta/goworkspace/lib"
)

var origin string

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("Bpmn", Bpmn)
}

func Bpmn(w http.ResponseWriter, r *http.Request) {
	log.Println("Bpmn")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1",
				Handler: BpmnFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAll},
			},
		},
	}
	route.Router(w, r)

}
func BpmnFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("--------------------------BpmnFx-------------------------------------------")
	var (
		policy models.Policy
		e      error
	)
	origin = r.Header.Get("Origin")
	jsonMap := make(map[string]interface{})
	rBody := lib.ErrorByte(io.ReadAll(r.Body))
	e = json.Unmarshal(rBody, &jsonMap)
	e = json.Unmarshal(rBody, &policy)
	j, e := policy.Marshal()
	log.Println("Proposal request proposal: ", string(j))
	defer r.Body.Close()
	return "", nil, e
}

func NewBpmn(data models.Policy) *State {
	// Init workflow with a name, and max concurrent tasks
	log.Println("--------------------------NewBpmn-------------------------------------------")
	var (
		state *State
	)
	state = &State{
		Handlers: make(map[string]func(state *State) error),
		Data:     data,
	}
	state.Handlers = make(map[string]func(state *State) error)
	return state
}
