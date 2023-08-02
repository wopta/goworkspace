package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	tr "github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

var (
	origin    string
	authToken models.AuthToken
)

//var policy *models.Policy

func EmitV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFxV2] Handler start ----------------------------------------")

	var (
		result     EmitRequest
		e          error
		firePolicy string
		policy     models.Policy
	)
	authToken, e = models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	origin = r.Header.Get("origin")
	firePolicy = lib.GetDatasetByEnv(origin, models.PolicyCollection)
	request := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFxV2] Request: %s", string(request))
	json.Unmarshal([]byte(request), &result)

	uid := result.Uid
	log.Printf("[EmitFxV2] Uid: %s", uid)

	docsnap := lib.GetFirestore(firePolicy, string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFxV2] Policy %s JSON: %s", uid, string(policyJsonLog))

	emitUpdatePolicy(&policy, result)
	responseEmit := EmitV2(&policy, result, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFxV2] Response: ", string(b))

	return string(b), responseEmit, e
}

func EmitV2(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var (
		responseEmit EmitResponse
	)

	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")
	policy.Uid = request.Uid // we should enforce the setting of the ID on proposal

	if policy.IsReserved && policy.Status != models.PolicyStatusWaitForApproval {
		emitApproval(policy)
	} else {

		log.Println("[EmitFxV2] AgencyUid: ", policy.AgencyUid)

		if policy.AgencyUid != "" {
			state := runBpmn(policy, "agency")
			log.Println("[EmitV2] state.Data Policy:", state.Data)
			policy = state.Data
		} else if policy.AgentUid != "" {
			runBpmn(policy, "agent")
		} else {
			log.Printf("[EmitV2] Policy Uid %s", request.Uid)
			ecommerceFlow(policy, origin)
		}

	}
	responseEmit = EmitResponse{UrlPay: policy.PayUrl, UrlSign: policy.SignUrl}
	policyJson, _ := policy.Marshal()
	log.Printf("[EmitV2] Policy %s: %s", request.Uid, string(policyJson))
	policy.Updated = time.Now().UTC()
	lib.SetFirestore(firePolicy, request.Uid, policy)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", guaranteFire)

	return responseEmit
}
func ecommerceFlow(policy *models.Policy, origin string) string {
	emitBase(policy, origin)

	emitSign(policy, origin)

	emitPay(policy, origin)
	return ""
}
func GetFlow[F any](policy models.Policy, funtions map[string]F) F {

	if policy.AgencyUid != "" {
		return funtions["agency"]
	} else if policy.AgentUid != "" {
		return funtions["agent"]
	} else {
		return funtions["ecommerce"]
	}

}
func setAdvice(policy *models.Policy, origin string) {

	policy.Payment = "manual"
	policy.StatusHistory = append(policy.StatusHistory, string(models.PolicyStatusToPay))
	policy.StatusHistory = append(policy.StatusHistory, string(models.PolicyStatusPay))
	policy.StatusHistory = append(policy.StatusHistory, string(models.PolicyStatusToSign))
	policy.Status = string(models.PolicyStatusToSign)

	policy.PaymentSplit = string(models.PaySingleInstallment)
	policy.IsPay = true
	tr.PutByPolicy(*policy, "", origin, "", "", policy.PriceGross, policy.PriceNett, "", true, authToken.Role)

}
func setAdviceBpm(state *bpmn.State) error {

	p := state.Data
	setAdvice(p, origin)
	return nil
}
func setData(state *bpmn.State) error {
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	p := state.Data
	emitBase(p, origin)
	return lib.SetFirestoreErr(firePolicy, p.Uid, p)
}

func sendMailSign(state *bpmn.State) error {
	policy := state.Data
	mail.SendMailSign(*policy)
	return nil
}
func sign(state *bpmn.State) error {
	policy := state.Data
	emitSign(policy, origin)
	return nil
}
func putUser(state *bpmn.State) error {
	policy := state.Data
	user.SetUserIntoPolicyContractor(policy, origin)
	return nil
}

func runBpmn(policy *models.Policy, channel string) *bpmn.State {
	settingByte, _ := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+channel+"/setting.json")
	//var prod models.Product

	var setting models.NodeSetting

	//Parsing/Unmarshalling JSON encoding/json
	json.Unmarshal(settingByte, &setting)
	state := bpmn.NewBpmn(*policy)
	state.AddTaskHandler("emitData", setData)
	state.AddTaskHandler("sendMailSign", sendMailSign)
	state.AddTaskHandler("sign", sign)
	state.AddTaskHandler("setAdvice", setAdviceBpm)
	state.AddTaskHandler("putUser", putUser)
	log.Println(state.Handlers)
	log.Println(state.Processes)
	state.RunBpmn(setting.EmitFlow)
	return state
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
        "name": "sign",
        "type": "TASK",
        "id": 1,
        "outProcess": [2],
        "inProcess": [0],
        "status": "READY"

    },

	{
        "name": "sendMailSign",
        "type": "TASK",
        "id": 2,
        "outProcess": [3],
        "inProcess": [1],
        "status": "READY"

    },

	{
        "name": "setAdvice",
        "type": "TASK",
        "id": 3,
        "outProcess": [4],
        "inProcess": [2],
        "status": "READY"

    },

	{
        "name": "putUser",
        "type": "TASK",
        "id": 4,
        "outProcess": [],
        "inProcess": [3],
        "status": "READY"

    }
]`
}
