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

func EmitV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFxV2] Handler start ----------------------------------------")

	var (
		request    EmitRequest
		e          error
		firePolicy string
		policy     models.Policy
	)
	authToken, e = models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	origin = r.Header.Get("origin")
	firePolicy = lib.GetDatasetByEnv(origin, models.PolicyCollection)
	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFxV2] Request: %s", string(body))
	json.Unmarshal([]byte(body), &request)

	uid := request.Uid
	log.Printf("[EmitFxV2] Uid: %s", uid)

	docsnap := lib.GetFirestore(firePolicy, string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFxV2] Policy %s JSON: %s", uid, string(policyJsonLog))

	emitUpdatePolicy(&policy, request)
	responseEmit := EmitV2(&policy, request, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFxV2] Response: ", string(b))

	return string(b), responseEmit, e
}

func EmitV2(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

	if policy.IsReserved && policy.Status != models.PolicyStatusWaitForApproval {
		emitApproval(policy)
	} else {
		log.Println("[EmitFxV2] AgencyUid: ", policy.AgencyUid)

		if policy.AgencyUid != "" {
			state := runBpmn(policy, models.AgencyChannel)
			log.Println("[EmitV2] state.Data Policy:", state.Data)
			policy = state.Data
		} else if policy.AgentUid != "" {
			runBpmn(policy, models.AgentChannel)
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
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	return responseEmit
}

func ecommerceFlow(policy *models.Policy, origin string) {
	emitBase(policy, origin)

	emitSign(policy, origin)

	emitPay(policy, origin)
}

func setAdvice(policy *models.Policy, origin string) {
	policy.Payment = "manual"
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay, models.PolicyStatusPay)
	policy.PaymentSplit = string(models.PaySingleInstallment)

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

func updateUserAndAgency(state *bpmn.State) error {
	policy := state.Data
	user.SetUserIntoPolicyContractor(policy, origin)
	return models.UpdateAgencyPortfolio(policy, origin)
}

func runBpmn(policy *models.Policy, channel string) *bpmn.State {
	settingByte, _ := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+channel+"/setting.json")

	var setting models.NodeSetting

	//Parsing/Unmarshalling JSON encoding/json
	json.Unmarshal(settingByte, &setting)
	state := bpmn.NewBpmn(*policy)
	state.AddTaskHandler("emitData", setData)
	state.AddTaskHandler("sendMailSign", sendMailSign)
	state.AddTaskHandler("sign", sign)
	state.AddTaskHandler("setAdvice", setAdviceBpm)
	state.AddTaskHandler("putUser", updateUserAndAgency)
	log.Println(state.Handlers)
	log.Println(state.Processes)
	state.RunBpmn(setting.EmitFlow)
	return state
}
