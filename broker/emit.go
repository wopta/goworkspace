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

	log.Println("Emit")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(request))
	json.Unmarshal([]byte(request), &result)
	uid := result.Uid
	log.Println(uid)

	var policy models.Policy
	docsnap := lib.GetFirestore("policy", string(uid))
	docsnap.DataTo(&policy)
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
	p := <-doc.ContractObj(policy)
	log.Println(p.LinkGcs)
	policy.DocumentName = p.LinkGcs
	_, res, _ := doc.NamirialOtpV6(policy)
	policy.ContractFileId = res.FileId
	company, numb := GetSequenceByProduct("global")
	policy.Number = numb
	policy.NumberCompany = company
	policy.IdSign = res.EnvelopeId
	var payRes pay.FabrickPaymentResponse
	if policy.PaymentSplit == string(models.PaySplitYear) {
		payRes = pay.FabbrickYearPay(policy)
	}
	if policy.PaymentSplit == string(models.PaySplitMonthly) {
		payRes = pay.FabbrickMontlyPay(policy)
	}
	responseEmit := EmitResponse{UrlPay: *payRes.Payload.PaymentPageURL, UrlSign: res.Url}
	lib.SetFirestore("policy", uid, policy)
	policyJson, e := policy.Marshal()
	policy.Data = string(policyJson)
	//e = lib.InsertRowsBigQuery("wopta", "policy", policy)
	mail.SendMail(getEmitMailObj(policy, responseEmit))
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

func getEmitMailObj(policy models.Policy, emitResponse EmitResponse) mail.MailRequest {
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>Ciao ` + policy.Contractor.Name + `` + policy.Contractor.Surname + ` </p>
	<p>Polizza n° ` + policy.NumberCompany + `</p> 
	<p>Grazie per aver scelto uno dei nostri prodotti Wopta per te</p> 
	<p>Puoi ora procedere alla firma della polizza in oggetto. Qui trovi il link per
	 accedere alla procedura semplice e guidata di firma elettronica avanzata tramite utilizzo di
	  un codice usa e getta che verrà inviato via sms sul tuo cellulare a noi comunicato. 
	Ti verrà richiesta l’adesione al servizio che è fornito in maniera gratuita da Wopta. 
	Potrai prendere visione delle condizioni generali di servizio e delle caratteristiche tecniche.</p> 
	<p><a class="button" href='` + emitResponse.UrlSign + `'>Firma la tua polizza:</a></p>
	<p>Ultimata la procedura di firma potrai procedere al pagamento. nella prossima mail  </p> 
	<p>Grazie per aver scelto Wopta </p> 
	<p>Proteggiamo chi sei</p> 
	`
	obj.Subject = " Wopta Contratto e pagamento"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
