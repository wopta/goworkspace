package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	prd "github.com/wopta/goworkspace/product"

	"cloud.google.com/go/civil"
	"github.com/wopta/goworkspace/callback_out"
	"github.com/wopta/goworkspace/document"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/payment"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/question"
	"github.com/wopta/goworkspace/transaction"
)

const (
	typeEmit    string = "emit"
	typeApprove string = "approve"
)

type EmitResponse struct {
	UrlPay       string               `json:"urlPay,omitempty"`
	UrlSign      string               `json:"urlSign,omitempty"`
	Uid          string               `json:"uid,omitempty"`
	ReservedInfo *models.ReservedInfo `json:"reservedInfo,omitempty"`
	CodeCompany  string               `json:"codeCompany,omitempty"`
}

type EmitRequest struct {
	BrokerBaseRequest
	Uid         string              `json:"uid,omitempty"`
	PaymentType string              `json:"paymentType,omitempty"`
	Statements  *[]models.Statement `json:"statements,omitempty"`
	SendEmail   *bool               `json:"sendEmail"`
}

func EmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request      EmitRequest
		err          error
		policy       models.Policy
		responseEmit EmitResponse
	)

	log.SetPrefix("[EmitFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin = r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal([]byte(body), &request)
	if err != nil {
		log.Printf("error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	uid := request.Uid
	log.Printf("Uid: %s", uid)

	paymentSplit = request.PaymentSplit
	log.Printf("paymentSplit: %s", paymentSplit)

	paymentMode = request.PaymentMode
	log.Printf("paymentMode: %s", paymentMode)

	policy, err = plc.GetPolicy(uid, origin)
	lib.CheckError(err)

	if policy.Channel == models.NetworkChannel && policy.ProducerUid != authToken.UserID {
		log.Printf("user %s cannot emit policy %s because producer not equal to request user", authToken.UserID, policy.Uid)
		return "", nil, errors.New("operation not allowed")
	}

	policyJsonLog, _ := policy.Marshal()
	log.Printf("Policy %s JSON: %s", uid, string(policyJsonLog))

	if policy.IsPay || policy.IsSign || policy.CompanyEmit || policy.CompanyEmitted || policy.IsDeleted {
		log.Printf("cannot emit policy %s because state is not correct", policy.Uid)
		return "", nil, errors.New("operation not allowed")
	}

	if request.SendEmail == nil {
		sendEmail = true
	} else {
		sendEmail = *request.SendEmail
	}

	productConfig := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	if err = policy.CheckStartDateValidity(productConfig.EmitMaxElapsedDays); err != nil {
		return "", "", err
	}

	emitUpdatePolicy(&policy, request)

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if policy.IsReserved && policy.Status != models.PolicyStatusApproved {
		log.Printf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
		return "", nil, fmt.Errorf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
	}
	responseEmit = emit(&policy, request, origin)

	b, e := json.Marshal(responseEmit)

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, e
}

func emit(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, lib.GuaranteeCollection)

	log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
	log.Println("[Emit] starting bpmn flow...")
	state := runBrokerBpmn(policy, emitFlowKey)
	if state == nil || state.Data == nil || state.IsFailed {
		log.Println("[Emit] error bpmn - state not set correctly")
		return responseEmit
	}
	*policy = *state.Data

	responseEmit = EmitResponse{
		UrlPay:       policy.PayUrl,
		UrlSign:      policy.SignUrl,
		ReservedInfo: policy.ReservedInfo,
		Uid:          policy.Uid,
		CodeCompany:  policy.CodeCompany,
	}

	policy.Updated = time.Now().UTC()
	policyJson, _ := policy.Marshal()
	log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))

	log.Println("[Emit] saving policy to firestore...")
	err := lib.SetFirestoreErr(firePolicy, request.Uid, policy)
	lib.CheckError(err)

	log.Println("[Emit] saving policy to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[Emit] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	callbackAction := callback_out.Emit
	if warrant != nil && warrant.GetFlowName(policy.Name) == models.RemittanceMgaFlow {
		callbackAction = callback_out.EmitRemittance
	}

	callback_out.Execute(networkNode, *policy, callbackAction)

	log.Println("[Emit] end --------------------------------------------------")
	return responseEmit
}

func emitUpdatePolicy(policy *models.Policy, request EmitRequest) {
	log.Println("[emitUpdatePolicy] start ------------------------------------")
	if policy.Statements == nil || len(*policy.Statements) == 0 {
		if request.Statements != nil {
			log.Println("[emitUpdatePolicy] inject policy statements from request")
			policy.Statements = request.Statements
		} else {
			log.Println("[emitUpdatePolicy] inject policy statements from question module")
			policy.Statements = new([]models.Statement)
			*policy.Statements, _ = question.GetStatements(policy)
		}
	}
	brokerUpdatePolicy(policy, request.BrokerBaseRequest)
	log.Println("[emitUpdatePolicy] end --------------------------------------")
}

func brokerUpdatePolicy(policy *models.Policy, request BrokerBaseRequest) {
	log.Println("[brokerUpdatePolicy] start ------------------------------------")
	if policy.PaymentSplit == "" {
		log.Println("[brokerUpdatePolicy] inject policy payment split from request")
		policy.PaymentSplit = request.PaymentSplit
	}

	if policy.Payment == "" {
		log.Println("[brokerUpdatePolicy] inject policy payment provider from request")
		policy.Payment = request.Payment
	}

	if policy.PaymentMode == "" {
		log.Println("[brokerUpdatePolicy] inject policy payment mode from request")
		policy.PaymentMode = request.PaymentMode
	}

	if policy.TaxAmount == 0 {
		log.Println("[brokerUpdatePolicy] calculate tax amount")
		policy.TaxAmount = lib.RoundFloat(policy.PriceGross - policy.PriceNett, 2)
	}

	if policy.TaxAmountMonthly == 0 {
		log.Println("[brokerUpdatePolicy] calculate tax amount monthly")
		policy.TaxAmountMonthly = lib.RoundFloat(policy.PriceGrossMonthly - policy.PriceNettMonthly, 2)
	}

	calculatePaymentComponents(policy)

	policy.SanitizePaymentData()

	log.Println("[brokerUpdatePolicy] end --------------------------------------")
}

func emitBase(policy *models.Policy, origin string) {
	log.Printf("[emitBase] Policy Uid %s", policy.Uid)
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	company, numb, tot := GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("[emitBase] codeCompany: %s", company)
	log.Printf("[emitBase] numberCompany: %d", numb)
	log.Printf("[emitBase] number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
	policy.RenewDate = policy.StartDate.AddDate(1, 0, 0)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
}

func emitSign(policy *models.Policy, origin string) {
	log.Printf("[emitSign] Policy Uid %s", policy.Uid)

	policy.IsSign = false
	policy.Status = models.PolicyStatusToSign
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusContact, models.PolicyStatusToSign)

	p := <-document.ContractObj(origin, *policy, networkNode, mgaProduct)
	policy.DocumentName = p.LinkGcs
	_, signResponse, _ := document.NamirialOtpV6(*policy, origin, sendEmail)
	policy.ContractFileId = signResponse.FileId
	policy.IdSign = signResponse.EnvelopeId
	policy.SignUrl = signResponse.Url
}

func emitPay(policy *models.Policy, origin string) {
	log.Printf("[emitPay] Policy Uid %s", policy.Uid)

	policy.IsPay = false
	payUrl, err := createPolicyTransactions(policy)
	if err != nil {
		return
	}
	policy.PayUrl = payUrl
}

func setAdvance(policy *models.Policy, origin string) {
	policy.Payment = models.ManualPaymentProvider
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay, models.PolicyStatusPay)

	//TODO: fix me someday in the future
	if paymentSplit != "" && policy.PaymentSplit == "" {
		policy.PaymentSplit = paymentSplit
	}
	if paymentMode != "" && policy.PaymentMode == "" {
		policy.PaymentMode = paymentMode
	}

	createPolicyTransactions(policy)
}

func createPolicyTransactions(policy *models.Policy) (string, error) {
	transactions := transaction.CreateTransactions(*policy, *mgaProduct, func() string { return lib.NewDoc(models.TransactionsCollection) })
	if len(transactions) == 0 {
		log.Println("no transactions created")
		return "", errors.New("no transactions created")
	}

	client := payment.NewClient(policy.Payment, *policy, *product, transactions, false, "")
	payUrl, updatedTransactions, err := client.NewBusiness()
	if err != nil {
		log.Printf("error emitPay policy %s: %s", policy.Uid, err.Error())
		return "", err
	}

	for index, tr := range updatedTransactions {
		err = lib.SetFirestoreErr(models.TransactionsCollection, tr.Uid, tr)
		if err != nil {
			log.Printf("error saving transaction %s to firestore: %s", tr.Uid, err.Error())
			return "", err
		}
		tr.BigQuerySave("")

		if tr.IsPay {
			err = transaction.CreateNetworkTransactions(policy, &updatedTransactions[index], networkNode, mgaProduct)
			if err != nil {
				log.Printf("error creating network transactions: %s", err.Error())
				return "", err
			}
		}
	}
	return payUrl, err
}

func calculatePaymentComponents(policy *models.Policy) {
	policy.PaymentComponents = models.PaymentComponents{
		Split:    models.PaySplit(policy.PaymentSplit),
		Rates:    models.PaySplitRateMap[models.PaySplit(policy.PaymentSplit)],
		Mode:     policy.PaymentMode,
		Provider: policy.Payment,
		PriceAnnuity: models.PriceComponents{
			Gross:       policy.PriceGross,
			Nett:        policy.PriceNett,
			Tax:         policy.TaxAmount,
			Consultancy: policy.ConsultancyValue.Price,
			Total:       lib.RoundFloat(policy.PriceGross+policy.ConsultancyValue.Price, 2),
		},
	}

	var priceSplit, priceFirstSplit models.PriceComponents
	switch policy.PaymentComponents.Split {
	case models.PaySplitSingleInstallment, models.PaySplitYearly, models.PaySplitYear:
		priceSplit = policy.PaymentComponents.PriceAnnuity
		priceFirstSplit = priceSplit
	case models.PaySplitSemestral:
	// TODO: unimplemented
	case models.PaySplitMonthly:
		priceSplit = models.PriceComponents{
			Gross:       policy.PriceGrossMonthly,
			Nett:        policy.PriceNettMonthly,
			Tax:         policy.TaxAmountMonthly,
			Consultancy: 0,
			Total:       policy.PriceGrossMonthly,
		}
		priceFirstSplit = priceSplit
		priceFirstSplit.Consultancy = policy.ConsultancyValue.Price
		priceFirstSplit.Total = lib.RoundFloat(priceFirstSplit.Gross+priceFirstSplit.Consultancy, 2)
	}

	policy.PaymentComponents.PriceSplit = priceSplit
	policy.PaymentComponents.PriceFirstSplit = priceFirstSplit
}
