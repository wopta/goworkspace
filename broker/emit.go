package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	models "github.com/wopta/goworkspace/models"
	pay "github.com/wopta/goworkspace/payment"
)

func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result EmitRequest
		e      error
	)

	log.Println("--------------------------Emit-------------------------------------------")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println("Emit", string(request))
	json.Unmarshal([]byte(request), &result)

	uid := result.Uid
	log.Println("Emit", uid)
	var policy models.Policy
	docsnap := lib.GetFirestore("policy", string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, e := policy.Marshal()
	log.Println("Emit get policy "+uid, string(policyJsonLog))
	policy.IsSign = false
	policy.IsPay = false
	policy.Updated = time.Now()
	policy.Uid = uid
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact)
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToSign)
	policy.PaymentSplit = result.PaymentSplit
	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = time.Now()
	policy.BigEmitDate = civil.DateTimeOf(policy.Updated)

	if policy.Statements == nil {
		policy.Statements = result.Statements
	}

	company, numb, tot := GetSequenceByCompany(policy.Company)
	log.Println("Emit code "+uid+" ", company)
	log.Println("Emit code "+uid+" ", numb)
	log.Println("Emit code "+uid+" ", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
	p := <-doc.ContractObj(policy)

	policy.DocumentName = p.LinkGcs
	_, res, _ := doc.NamirialOtpV6(policy)
	policy.ContractFileId = res.FileId
	policy.IdSign = res.EnvelopeId
	var payRes pay.FabrickPaymentResponse
	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = pay.FabbrickYearPay(policy)
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = pay.FabbrickMontlyPay(policy)
	}
	responseEmit := EmitResponse{UrlPay: *payRes.Payload.PaymentPageURL, UrlSign: res.Url}
	policy.SignUrl = res.Url
	policy.PayUrl = *payRes.Payload.PaymentPageURL
	policyJson, e := policy.Marshal()
	log.Println("Emit policy "+uid, string(policyJson))
	lib.SetFirestore("policy", uid, policy)
	policy.BigquerySave()
	mail.SendMailSign(policy)
	b, e := json.Marshal(responseEmit)
	return string(b), responseEmit, e
}

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
	Survay       *[]models.Statement `firestore:"survey,omitempty" json:"survey,omitempty"`
	Statements   *[]models.Statement `firestore:"statements,omitempty" json:"statements,omitempty"`
}
