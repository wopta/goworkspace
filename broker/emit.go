package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
	models "github.com/wopta/goworkspace/models"
	pay "github.com/wopta/goworkspace/payment"
)

func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result map[string]string
	)
	uid := result["uid"]
	log.Println("Emit")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	log.Println(string(request))
	json.Unmarshal([]byte(request), &result)
	log.Println(uid)
	log.Println(result["paymentSplit"])
	log.Println(result["payment"])
	var policy models.Policy
	docsnap := lib.GetFirestore("policy", string(result["uid"]))
	docsnap.DataTo(&policy)
	company, numb := GetSequenceByProduct("global")
	policy.NumberCompany = company
	policy.IsSign = false
	policy.IsPay = false
	policy.Number = numb
	policy.Updated = time.Now()
	policy.Uid = result["uid"]
	policy.Status = models.PolicyStatusToEmit
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToEmit)
	policy.PaymentSplit = result["paymentSplit"]
	p := <-doc.ContractObj(policy)
	log.Println(p.LinkGcs)
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
	lib.SetFirestore("policy", result["uid"], policy)
	e := lib.InsertRowsBigQuery("wopta", "policy", policy)
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
	Uid          string              `firestore:"uid,omitempty" json:"uid,omitempty" bigquery:"uid"`
	Payment      string              `firestore:"payment,omitempty" json:"payment,omitempty" bigquery:"payment"`
	PaymentType  string              `firestore:"paymentType,omitempty" json:"paymentType,omitempty" bigquery:"paymentType"`
	PaymentSplit string              `firestore:"paymentSplit,omitempty" json:"paymentSplit,omitempty" bigquery:"paymentSplit"`
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
	<p>Ultimata la procedura di firma potrai procedere al pagamento. Clicca nel link per pagare la tua polizza  </p> 
	<p><a class="button" href='` + emitResponse.UrlPay + `'>Paga la tua polizza</a></p>
	<p>Grazie per aver scelto Wopta </p> 
	<p>Proteggiamo chi sei</p> 
	`
	obj.Subject = " Wopta Contratto e pagamento"
	obj.IsHtml = true
	obj.IsAttachment = false

	return obj
}
