package AppcheckProxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	pay "github.com/wopta/goworkspace/payment"

	doc "github.com/wopta/goworkspace/document"
	models "github.com/wopta/goworkspace/models"
	//lib "github.com/wopta/goworkspace/lib"
)

func init() {
	log.Println("INIT AppcheckProxy")
	functions.HTTP("Bpmn", Bpmn)
}

func Bpmn(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback")
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
	)
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	e := json.Unmarshal([]byte(req), &policy)
	j, e := policy.Marshal()
	log.Println("Proposal request proposal: ", string(j))
	defer r.Body.Close()
	BpmnEngine(policy)
	return "", nil, e
}
func BpmnEngine(policy models.Policy) string {
	// Init workflow with a name, and max concurrent tasks
	log.Println("--------------------------BpmnEngine-------------------------------------------")
	bpmnEngine := bpmn_engine.New("a name")
	// basic example loading a BPMN from file,
	filePath := "./serverless_function_source_code/test.bpmn"
	process, err := bpmnEngine.LoadFromFile(filePath)
	//lib.C
	if err != nil {
		panic("file \"simple_task.bpmn\" can't be read.")
	}
	// register a handler for a service task by defined task type
	bpmnEngine.AddTaskHandler("test", test)
	bpmnEngine.AddTaskHandler("fabrickPayment", fabrickPayment)
	bpmnEngine.AddTaskHandler("contract", contract)
	bpmnEngine.AddTaskHandler("namirialSign", namirialSign)
	bpmnEngine.AddTaskHandler("sendMailSign", sendMailSign)

	// setup some variables
	variables := map[string]interface{}{}
	//variables["policy"] = policy
	variables["policy"] = "test"
	// and execute the process
	bpmnEngine.CreateAndRunInstance(process.ProcessKey, variables)
	return ""
}
func test(job bpmn_engine.ActivatedJob) {
	log.Println("--------------------------Test-------------------------------------------")
	fmt.Printf("job.GetState(): %v\n", job.GetState())
	fmt.Printf(" job.GetVariable(policy): %v\n", job.GetVariable("policy"))
}
func contract(job bpmn_engine.ActivatedJob) {
	policy := job.GetVariable("policy")
	doc.ContractObj(policy.(models.Policy))

}
func fabrickPayment(job bpmn_engine.ActivatedJob) {
	var payRes pay.FabrickPaymentResponse
	policye := job.GetVariable("policy")
	policy := policye.(models.Policy)
	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = pay.FabbrickYearPay(policy, "")
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = pay.FabbrickMontlyPay(policy, "")

	}
	state := job.GetState()
	fmt.Printf("state: %v\n", state)
	fmt.Printf("payRes: %v\n", payRes)

}

func namirialSign(job bpmn_engine.ActivatedJob) {
	policy := job.GetVariable("policy")
	doc.NamirialOtpV6(policy.(models.Policy), "")

}
func sendMailSign(job bpmn_engine.ActivatedJob) {
	policy := job.GetVariable("policy")
	mail.SendMailSign(policy.(models.Policy))

}
