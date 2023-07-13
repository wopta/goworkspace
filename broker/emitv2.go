package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
)

var origin string

func EmitV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFxV2] Handler start ----------------------------------------")

	var (
		result     EmitRequest
		e          error
		firePolicy string
		policy     models.Policy
	)

	origin = r.Header.Get("origin")
	firePolicy = lib.GetDatasetByEnv(origin, "policy")
	request := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFxV2] Request: %s", string(request))
	json.Unmarshal([]byte(request), &result)

	uid := result.Uid
	log.Printf("[EmitFxV2] Uid: %s", uid)

	docsnap := lib.GetFirestore(firePolicy, string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFxV2] Policy %s JSON: %s", uid, string(policyJsonLog))

	responseEmit := EmitV2(&policy, result, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFxV2] Response: ", string(b))

	return string(b), responseEmit, e
}

func EmitV2(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")
	policy.Uid = request.Uid // we should enforce the setting of the ID on proposal

	if policy.IsReserved && policy.Status != models.PolicyStatusWaitForApproval {
		emitApproval(policy)
	} else {
		log.Println("[EmitFxV2] AgencyUid: ", policy.AgencyUid)
		if policy.AgencyUid != "" {

			runBpmn(policy, getTest())
		} else if policy.AgencyUid != "" {

		} else {
			log.Printf("[Emit] Policy Uid %s", request.Uid)

			emitBase(policy, origin)

			emitSign(policy, origin)

			emitPay(policy, origin)

			responseEmit = EmitResponse{UrlPay: policy.PayUrl, UrlSign: policy.SignUrl}
			policyJson, _ := policy.Marshal()
			log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))
		}
	}

	policy.Updated = time.Now().UTC()
	lib.SetFirestore(firePolicy, request.Uid, policy)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", guaranteFire)

	return responseEmit
}

func emitData(state *bpmn.State) error {
	p := state.Data.(models.Policy)

	emitBase(&p, origin)
	log.Println(p)
	log.Println(state.Data.(models.Policy))
	return nil
}

func sendMailSign(state *bpmn.State) error {
	policy := state.Data
	mail.SendMailSign(policy.(models.Policy))
	return nil
}

func runBpmn(policy *models.Policy, processByte string) {
	state := bpmn.NewBpmn(policy)
	state.AddTaskHandler("emitData", emitData)

	state.AddTaskHandler("sendMailSign", sendMailSign)
	log.Println(state.Handlers)
	log.Println(state.Processes)
	process, _ := state.LoadProcesses(processByte)
	state.RunBpmn(process)

}
func getTest() string {
	return `
	[{
        "name": "emitData",
        "type": "TASK",
        "id": 0,
        "outProcess": [1],
        "inProcess": [],
        "status": "READY"

    },

	{
        "name": "sendMailSign",
        "type": "TASK",
        "id": 1,
        "outProcess": [],
        "inProcess": [1],
        "status": "READY"

    },
]
	`
}
