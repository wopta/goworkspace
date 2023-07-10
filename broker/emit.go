package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/payment"
)

type EmitResponse struct {
	UrlPay  string `firestore:"urlPay,omitempty" json:"urlPay,omitempty"`
	UrlSign string `firestore:"urlSign,omitempty" json:"urlSign,omitempty"`
	Uid     string `firestore:"uid,omitempty" json:"uid,omitempty"`
}

type EmitRequest struct {
	Uid          string              `firestore:"uid,omitempty" json:"uid,omitempty"`
	Payment      string              `firestore:"payment,omitempty" json:"payment,omitempty"`
	PaymentType  string              `firestore:"paymentType,omitempty" json:"paymentType,omitempty"`
	PaymentSplit string              `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty"`
	Statements   *[]models.Statement `firestore:"statements,omitempty" json:"statements,omitempty"`
}

func EmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFx] Handler start ----------------------------------------")

	var (
		request EmitRequest
		e       error
		policy  models.Policy
	)

	origin := r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFx] Request body: %s", string(body))
	json.Unmarshal([]byte(body), &request)

	uid := request.Uid
	log.Printf("[EmitFx] Uid: %s", uid)
	policy, e = GetPolicy(uid, origin)
	lib.CheckError(e)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFx] Policy %s JSON: %s", uid, string(policyJsonLog))

	emitUpdatePolicy(&policy, request)
	responseEmit := Emit(&policy, origin)

	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFx] Response: ", string(b))

	return string(b), responseEmit, e
}

func Emit(policy *models.Policy, origin string) EmitResponse {
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")

	if policy.IsReserved && policy.Status != models.PolicyStatusWaitForApproval {
		emitApproval(policy)
	} else {
		log.Printf("[Emit] Policy Uid %s", policy.Uid)

		emitBase(policy, origin)

		emitSign(policy, origin)

		emitPay(policy, origin)

		responseEmit = EmitResponse{UrlPay: policy.PayUrl, UrlSign: policy.SignUrl}
		policyJson, _ := policy.Marshal()
		log.Printf("[Emit] Policy %s: %s", policy.Uid, string(policyJson))
	}

	policy.Updated = time.Now().UTC()
	lib.SetFirestore(firePolicy, policy.Uid, policy)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", guaranteFire)

	return responseEmit
}

func emitUpdatePolicy(policy *models.Policy, request EmitRequest) {
	if policy.Status == models.PolicyStatusInitLead {
		if policy.Statements == nil {
			policy.Statements = request.Statements
		}
		policy.PaymentSplit = request.PaymentSplit
	}
}

func emitApproval(policy *models.Policy) {
	log.Printf("[EmitApproval] Policy Uid %s: Reserved Flow", policy.Uid)
	policy.Status = models.PolicyStatusWaitForApproval
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
}

func emitBase(policy *models.Policy, origin string) {
	log.Printf("[EmitBase] Policy Uid %s", policy.Uid)
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
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

	p := <-document.ContractObj(*policy)
	policy.DocumentName = p.LinkGcs
	_, signResponse, _ := document.NamirialOtpV6(*policy, origin)
	policy.ContractFileId = signResponse.FileId
	policy.IdSign = signResponse.EnvelopeId
	policy.SignUrl = signResponse.Url

	mail.SendMailSign(*policy)
}

func emitPay(policy *models.Policy, origin string) {
	log.Printf("[EmitPay] Policy Uid %s", policy.Uid)
	var payRes payment.FabrickPaymentResponse

	policy.IsPay = false

	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = payment.FabbrickYearPay(*policy, origin)
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = payment.FabbrickMontlyPay(*policy, origin)
	}

	policy.PayUrl = *payRes.Payload.PaymentPageURL
}
