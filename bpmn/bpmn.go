package bpmn

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	pay "github.com/wopta/goworkspace/payment"

	doc "github.com/wopta/goworkspace/document"
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
	BpmnEngine(policy, jsonMap)
	return "", nil, e
}
func BpmnEngine(policy models.Policy, mapPolicy map[string]interface{}) string {
	// Init workflow with a name, and max concurrent tasks
	log.Println("--------------------------BpmnEngine-------------------------------------------")
	var (
		state *State
	)
	state = &State{
		Handlers: make(map[string]func(state *State) error),
		Data:     policy,
	}
	state.Handlers = make(map[string]func(state *State) error)
	// basic example loading a BPMN from file,
	//filePath := "./serverless_function_source_code/test.bpmn"
	processes, err := state.LoadProcesses(getTest())
	//lib.C
	if err != nil {
		log.Println(err)

	}
	// register a handler for a service task by defined task type
	state.AddTaskHandler("test", test)
	state.AddTaskHandler("fabrickPayment", fabrickPayment)
	state.AddTaskHandler("contract", contract)
	state.AddTaskHandler("namirialSign", namirialSign)
	state.AddTaskHandler("sendMailSign", sendMailSign)

	state.RunBpmn(processes)
	return ""
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
func test(state *State) error {
	log.Println("--------------------------Test-------------------------------------------")
	return nil
}
func contract(state *State) error {
	policy := state.Data
	doc.ContractObj(origin, policy)
	return nil
}
func fabrickPayment(state *State) error {
	var payRes pay.FabrickPaymentResponse
	policye := state.Data
	policy := policye
	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = pay.FabbrickYearPay(policy, "")
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = pay.FabbrickMontlyPay(policy, "")

	}

	fmt.Printf("state: %v\n", state)
	fmt.Printf("payRes: %v\n", payRes)
	return nil

}

func namirialSign(state *State) error {
	policy := state.Data
	doc.NamirialOtpV6(policy, "")
	return nil
}
func sendMailSign(state *State) error {
	policy := state.Data
	mail.SendMailSign(&policy)
	return nil
}
func getTest() string {
	return `
	[{
        "name": "test",
        "type": "TASK",
        "id": 0,
        "outProcess": [1],
        "inProcess": [],
        "status": "READY"

    },
	{
        "name": "test",
        "type": "DECISION",
        "id": 1,
        "outTrueProcess": [2],
		"outFalseProcess": [2],
		"decision":"payment == \"fabrick\"",
        "inProcess": [0],
        "status": "READY"

    },
	{
        "name": "test",
        "type": "TASK",
        "id": 2,
        "outProcess": [],
        "inProcess": [1],
        "status": "READY"

    },
	,
	{
        "name": "test",
        "type": "TASK",
        "id": 2,
        "outProcess": [],
        "inProcess": [1],
        "status": "READY"

    }]
	`
}
