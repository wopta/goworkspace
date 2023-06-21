package broker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	dag "github.com/heimdalr/dag"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("Broker")
	lib.EnableCors(&w, r)
	route := lib.RouteData{
		Routes: []lib.Route{
			{
				Route:   "/v1/policies/fiscalcode/:fiscalcode",
				Handler: PolicyFiscalcode,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/:uid",
				Handler: GetPolicyFx,
				Method:  "GET",
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/v1/policy/proposal",
				Handler: Proposal,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},

			{
				Route:   "/v1/policy/emit",
				Handler: Emit,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "/v1/policy/reserved",
				Handler: reserved,
				Method:  "POST",
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policy/v1/:uid",
				Handler: UpdatePolicy,
				Method:  http.MethodPatch,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policy/v1/:uid",
				Handler: DeletePolicy,
				Method:  http.MethodDelete,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "attachment/v1/:policyUid",
				Handler: GetPolicyAttachmentFx,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAll},
			},
			{
				Route:   "policies/v1",
				Handler: GetPoliciesFx,
				Method:  http.MethodPost,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
			{
				Route:   "policy/transactions/v1/:policyUid",
				Handler: GetPolicyTransactions,
				Method:  http.MethodGet,
				Roles:   []string{models.UserRoleAdmin, models.UserRoleManager},
			},
		},
	}
	route.Router(w, r)

}

func GetNumberCompany(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	return "", nil
}

type BrokerResponse struct {
	EnvelopSignId string `json:"envelopSignId"`
	LinkGcs       string `json:"linkGcs"`
	Bytes         string `json:"bytes"`
}

func reserved(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	// initialize a new graph
	d := dag.NewDAG()

	// init three vertices
	v1, _ := d.AddVertex(1)
	v2, _ := d.AddVertex(2)
	v3, _ := d.AddVertex(struct {
		a string
		b string
	}{a: "foo", b: "bar"})

	// add the above vertices and connect them with two edges
	_ = d.AddEdge(v1, v2)
	_ = d.AddEdge(v1, v3)

	// describe the graph
	fmt.Print(d.String())

	return "", nil, nil
}
func test() string {
	// Init workflow with a name, and max concurrent tasks
	bpmnEngine := bpmn_engine.New("a name")
	// basic example loading a BPMN from file,
	process, err := bpmnEngine.LoadFromFile("simple_task.bpmn")
	if err != nil {
		panic("file \"simple_task.bpmn\" can't be read.")
	}
	// register a handler for a service task by defined task type
	bpmnEngine.AddTaskHandler("hello-world", printContextHandler)
	// setup some variables
	variables := map[string]interface{}{}
	variables["foo"] = "bar"
	// and execute the process
	bpmnEngine.CreateAndRunInstance(process.ProcessKey, variables)
	return ""
}
