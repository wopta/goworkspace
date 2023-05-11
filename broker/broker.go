package broker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	dag "github.com/heimdalr/dag"
	lib "github.com/wopta/goworkspace/lib"
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
			},
			{
				Route:   "/v1/policy/:uid",
				Handler: GetPolicyFx,
				Method:  "GET",
			},

			{
				Route:   "/v1/policy/proposal",
				Handler: Proposal,
				Method:  "POST",
			},

			{
				Route:   "/v1/policy/emit",
				Handler: Emit,
				Method:  "POST",
			},
			{
				Route:   "/v1/policy/reserved",
				Handler: reserved,
				Method:  "POST",
			},
			{
				Route:   "policy/v1/:uid",
				Handler: UpdatePolicy,
				Method:  http.MethodPatch,
			},
			{
				Route:   "attachment/v1/:policyUid",
				Handler: GetPolicyAttachmentFx,
				Method:  http.MethodGet,
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
