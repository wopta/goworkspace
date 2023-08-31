package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
	"github.com/wopta/goworkspace/question"
	"github.com/wopta/goworkspace/reserved"
	tr "github.com/wopta/goworkspace/transaction"
	"github.com/wopta/goworkspace/user"
)

var origin string

const (
	typeEmit    string = "emit"
	typeApprove string = "approve"
)

type EmitResponse struct {
	UrlPay       string               `firestore:"urlPay,omitempty" json:"urlPay,omitempty"`
	UrlSign      string               `firestore:"urlSign,omitempty" json:"urlSign,omitempty"`
	Uid          string               `firestore:"uid,omitempty" json:"uid,omitempty"`
	ReservedInfo *models.ReservedInfo `json:"reservedInfo,omitempty" firestore:"reservedInfo,omitempty"`
}

type EmitRequest struct {
	Uid          string              `firestore:"uid,omitempty" json:"uid,omitempty"`
	Payment      string              `firestore:"payment,omitempty" json:"payment,omitempty"`
	PaymentType  string              `firestore:"paymentType,omitempty" json:"paymentType,omitempty"`
	PaymentSplit string              `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty"`
	Statements   *[]models.Statement `firestore:"statements,omitempty" json:"statements,omitempty"`
}

func EmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFx] Handler start --------------------------------------")

	var (
		request EmitRequest
		err     error
		policy  models.Policy
	)

	origin = r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFx] Request: %s", string(body))
	json.Unmarshal([]byte(body), &request)

	uid := request.Uid
	log.Printf("[EmitFx] Uid: %s", uid)

	policy, err = GetPolicy(uid, origin)
	lib.CheckError(err)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFx] Policy %s JSON: %s", uid, string(policyJsonLog))

	emitUpdatePolicy(&policy, request)
	responseEmit := Emit(&policy, request, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFx] Response: ", string(b))

	return string(b), responseEmit, e
}

func Emit(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

	emitType := getEmitTypeFromPolicy(policy)
	switch emitType {
	case typeApprove:
		log.Printf("[Emit] Wait for approval - Policy Uid %s", policy.Uid)
		emitApproval(policy)
		reserved.GetReservedInfo(policy)
		mail.SendMailReserved(*policy)
	case typeEmit:
		log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
		if policy.AgencyUid != "" {
			log.Println("[Emit] Agency Flow")
			state := runBpmn(policy, models.AgencyChannel)
			log.Println("[Emit] state.Data Policy:", state.Data)
			policy = state.Data
		} else if policy.AgentUid != "" {
			log.Println("[Emit] Agent (E-commerce) Flow")
			ecommerceFlow(policy, origin)
		} else {
			log.Println("[Emit] E-commerce Flow")
			ecommerceFlow(policy, origin)
		}
	default:
		log.Printf("[Emit] ERROR cannot emit policy")
		return responseEmit
	}

	responseEmit = EmitResponse{
		UrlPay:       policy.PayUrl,
		UrlSign:      policy.SignUrl,
		ReservedInfo: policy.ReservedInfo,
		Uid:          policy.Uid,
	}
	policyJson, _ := policy.Marshal()
	log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))
	policy.Updated = time.Now().UTC()
	err := lib.SetFirestoreErr(firePolicy, request.Uid, policy)
	lib.CheckError(err)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	return responseEmit
}

func ecommerceFlow(policy *models.Policy, origin string) {
	emitBase(policy, origin)

	emitSign(policy, origin)

	emitPay(policy, origin)

	mail.SendMailSign(*policy)
}

func setAdvice(policy *models.Policy, origin string) {
	policy.Payment = "manual"
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay, models.PolicyStatusPay)
	policy.PaymentSplit = string(models.PaySingleInstallment)

	tr.PutByPolicy(*policy, "", origin, "", "", policy.PriceGross, policy.PriceNett, "", true)
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

func emitUpdatePolicy(policy *models.Policy, request EmitRequest) {
	if policy.Status != models.PolicyStatusInitLead {
		return
	}
	if policy.Statements == nil || len(*policy.Statements) == 0 {
		if request.Statements != nil {
			policy.Statements = request.Statements
		} else {
			*policy.Statements = question.GetStatements(*policy)
		}
	}
	policy.PaymentSplit = request.PaymentSplit
}

func getEmitTypeFromPolicy(policy *models.Policy) string {
	if !policy.IsReserved || policy.Status == models.PolicyStatusApproved {
		return typeEmit
	}

	deniedStatuses := []string{models.PolicyStatusDeleted, models.PolicyStatusRejected}

	if policy.IsReserved && !lib.SliceContains(deniedStatuses, policy.Status) {
		return typeApprove
	}

	return ""
}

func emitApproval(policy *models.Policy) {
	log.Printf("[EmitApproval] Policy Uid %s: Reserved Flow", policy.Uid)
	policy.Status = models.PolicyStatusWaitForApproval
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
}

func emitBase(policy *models.Policy, origin string) {
	log.Printf("[EmitBase] Policy Uid %s", policy.Uid)
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	company, numb, tot := GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("[EmitBase] codeCompany: %s", company)
	log.Printf("[EmitBase] numberCompany: %d", numb)
	log.Printf("[EmitBase] number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
}

func emitSign(policy *models.Policy, origin string) {
	log.Printf("[EmitSign] Policy Uid %s", policy.Uid)

	policy.IsSign = false
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact, models.PolicyStatusToSign)

	p := <-document.ContractObj(origin, *policy)
	policy.DocumentName = p.LinkGcs
	_, signResponse, _ := document.NamirialOtpV6(*policy, origin)
	policy.ContractFileId = signResponse.FileId
	policy.IdSign = signResponse.EnvelopeId
	policy.SignUrl = signResponse.Url
}

func emitPay(policy *models.Policy, origin string) {
	log.Printf("[EmitPay] Policy Uid %s", policy.Uid)
	//var payRes payment.FabrickPaymentResponse

	policy.IsPay = false

	/*if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = payment.FabbrickYearPay(*policy, origin)
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = payment.FabbrickMontlyPay(*policy, origin)
	}*/

	policy.PayUrl, _ = payment.PaymentController(origin, *policy)

}
