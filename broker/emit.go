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
	log.Println("[EmitFx] handler start ----------------------------------------")

	var (
		result     EmitRequest
		e          error
		firePolicy string
		policy     models.Policy
	)

	origin := r.Header.Get("origin")
	firePolicy = lib.GetDatasetByEnv(origin, "policy")
	request := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFx] request: %s", string(request))
	json.Unmarshal([]byte(request), &result)

	uid := result.Uid
	log.Printf("[EmitFx] Uid: %s", uid)

	docsnap := lib.GetFirestore(firePolicy, string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFx] Policy %s JSON: %s", uid, string(policyJsonLog))

	responseEmit := Emit(&policy, result, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFx] response: ", string(b))

	return string(b), responseEmit, e
}

func Emit(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var payRes payment.FabrickPaymentResponse
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")
	now := time.Now().UTC()

	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = now
	policy.Uid = request.Uid
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact)
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToSign)
	policy.PaymentSplit = request.PaymentSplit
	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)

	if policy.Statements == nil {
		policy.Statements = request.Statements
	}

	company, numb, tot := GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("[Emit] Policy Uid %s", request.Uid)
	log.Printf("[Emit] codeCompany: %s", company)
	log.Printf("[Emit] numberCompany: %d", numb)
	log.Printf("[Emit] number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company

	p := <-document.ContractObj(*policy)
	policy.DocumentName = p.LinkGcs
	_, signResponse, _ := document.NamirialOtpV6(*policy, origin)
	policy.ContractFileId = signResponse.FileId
	policy.IdSign = signResponse.EnvelopeId

	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = payment.FabbrickYearPay(*policy, origin)
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = payment.FabbrickMontlyPay(*policy, origin)
	}

	responseEmit := EmitResponse{UrlPay: *payRes.Payload.PaymentPageURL, UrlSign: signResponse.Url}
	policy.SignUrl = signResponse.Url
	policy.PayUrl = *payRes.Payload.PaymentPageURL
	policyJson, _ := policy.Marshal()

	log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))
	lib.SetFirestore(firePolicy, request.Uid, policy)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", guaranteFire)
	mail.SendMailSign(*policy)

	return responseEmit
}
