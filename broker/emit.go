package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	draftbpnm "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/question"
)

var (
	Proposal        string = "Proposal"
	RequestApproval string = "RequestApproval"
	Emit            string = "Emit"
	Signed          string = "Signed"
	Paid            string = "Paid"
	EmitRemittance  string = "EmitRemittance"
	Approved        string = "Approved"
	Rejected        string = "Rejected"
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

func DraftEmitFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		request      EmitRequest
		err          error
		policy       models.Policy
		responseEmit EmitResponse
	)

	log.AddPrefix("EmitFx")
	defer log.PopPrefix()

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	defer r.Body.Close()

	err = json.Unmarshal([]byte(body), &request)
	if err != nil {
		log.Printf("error unmarshalling policy: %s", err.Error())
		return "", nil, err
	}

	uid := request.Uid
	log.Printf("Uid: %s", uid)

	paymentMode = request.PaymentMode
	log.Printf("paymentMode: %s", paymentMode)

	policy, err = plc.GetPolicy(uid, origin)
	lib.CheckError(err)
	//	if policy.Channel == models.NetworkChannel && policy.ProducerUid != authToken.UserID {
	//		log.Printf("user %s cannot emit policy %s because producer not equal to request user", authToken.UserID, policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

	policyJsonLog, _ := policy.Marshal()
	log.Printf("Policy %s JSON: %s", uid, string(policyJsonLog))
	//	if policy.IsPay || policy.IsSign || policy.CompanyEmit || policy.CompanyEmitted || policy.IsDeleted {
	//		log.Printf("cannot emit policy %s because state is not correct", policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

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
	responseEmit, err = emitDraft(&policy, request, origin)
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(responseEmit)

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, err
}

func emitDraft(policy *models.Policy, request EmitRequest, origin string) (EmitResponse, error) {
	log.Println("start ------------------------------------------------")
	var responseEmit EmitResponse

	fireGuarantee := lib.GetDatasetByEnv(origin, lib.GuaranteeCollection)

	log.Printf("Emitting - Policy Uid %s", policy.Uid)
	log.Println("starting bpmn flow...")

	paymentSplit = request.PaymentSplit
	log.Printf("paymentSplit: %s", paymentSplit)

	storage := draftbpnm.NewStorageBpnm()
	if request.SendEmail == nil {
		storage.AddGlobal("sendEmail", &flow.BoolBpmn{Bool: true})
	} else {
		storage.AddGlobal("sendEmail", &flow.BoolBpmn{Bool: *request.SendEmail})
	}
	storage.AddGlobal("paymentSplit", &flow.String{String: request.PaymentSplit})
	storage.AddGlobal("paymentMode", &flow.String{String: request.PaymentMode})
	storage.AddGlobal("addresses", &flow.Addresses{FromAddress: mail.AddressAnna})

	log.Printf("paymentMode: %s", paymentMode)
	flow, err := getFlow(policy, origin, storage)
	if err != nil {
		return responseEmit, err
	}
	err = flow.Run("emit")
	if err != nil {
		return responseEmit, err
	}

	responseEmit = EmitResponse{
		UrlPay:       policy.PayUrl,
		UrlSign:      policy.SignUrl,
		ReservedInfo: policy.ReservedInfo,
		Uid:          policy.Uid,
		CodeCompany:  policy.CodeCompany,
	}

	policy.Updated = time.Now().UTC()
	policyJson, _ := policy.Marshal()
	log.Printf("Policy %s: %s", request.Uid, string(policyJson))

	log.Println("saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	log.Println("end --------------------------------------------------")
	return responseEmit, nil
}

func emitUpdatePolicy(policy *models.Policy, request EmitRequest) {
	log.AddPrefix("emitUpdatePolicy")
	defer log.PopPrefix()
	log.Println("start ------------------------------------")
	if policy.Statements == nil || len(*policy.Statements) == 0 {
		if request.Statements != nil {
			log.Println("inject policy statements from request")
			policy.Statements = request.Statements
		} else {
			log.Println("inject policy statements from question module")
			policy.Statements = new([]models.Statement)
			*policy.Statements, _ = question.GetStatements(policy)
		}
	}
	brokerUpdatePolicy(policy, request.BrokerBaseRequest)
	log.Println("end --------------------------------------")
}

func brokerUpdatePolicy(policy *models.Policy, request BrokerBaseRequest) {
	log.AddPrefix("brokerUpdatePolicy")
	defer log.PopPrefix()
	log.Println("start ------------------------------------")
	if policy.PaymentSplit == "" {
		log.Println("inject policy payment split from request")
		policy.PaymentSplit = request.PaymentSplit
	}

	if policy.Payment == "" {
		log.Println("inject policy payment provider from request")
		policy.Payment = request.Payment
	}

	if policy.PaymentMode == "" {
		log.Println("inject policy payment mode from request")
		policy.PaymentMode = request.PaymentMode
	}

	if policy.TaxAmount == 0 {
		policy.TaxAmount = lib.RoundFloat(policy.PriceGross-policy.PriceNett, 2)
	}

	if policy.TaxAmountMonthly == 0 {
		log.Println("[brokerUpdatePolicy] calculate tax amount monthly")
		policy.TaxAmountMonthly = lib.RoundFloat(policy.PriceGrossMonthly-policy.PriceNettMonthly, 2)
	}

	calculatePaymentComponents(policy)

	policy.SanitizePaymentData()

	log.Println("end --------------------------------------")
}

func emitBase(policy *models.Policy, origin string) {
	log.AddPrefix("emitBase")
	defer log.PopPrefix()
	log.Printf("Policy Uid %s", policy.Uid)
	firePolicy := lib.PolicyCollection
	now := time.Now().UTC()

	policy.CompanyEmit = true
	policy.CompanyEmitted = false
	policy.EmitDate = now
	policy.BigEmitDate = civil.DateTimeOf(now)
	company, numb, tot := utility.GetSequenceByCompany(strings.ToLower(policy.Company), firePolicy)
	log.Printf("codeCompany: %s", company)
	log.Printf("numberCompany: %d", numb)
	log.Printf("number: %d", tot)
	policy.Number = tot
	policy.NumberCompany = numb
	policy.CodeCompany = company
	policy.RenewDate = policy.StartDate.AddDate(1, 0, 0)
	policy.BigRenewDate = civil.DateTimeOf(policy.RenewDate)
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
	case models.PaySplitSemestral, models.PaySplitTrimestral:
		if policy.PaymentComponents.Rates == 0 {
			log.ErrorF("Rates should not be 0")
			return
		}
		priceSplit = models.PriceComponents{
			Gross:       policy.PriceGross / float64(policy.PaymentComponents.Rates),
			Nett:        policy.PriceGross / float64(policy.PaymentComponents.Rates),
			Tax:         policy.TaxAmount / float64(policy.PaymentComponents.Rates),
			Consultancy: 0,
			Total:       policy.PriceGross / float64(policy.PaymentComponents.Rates),
		}
		priceFirstSplit = priceSplit
		priceFirstSplit.Consultancy = policy.ConsultancyValue.Price
		priceFirstSplit.Total = lib.RoundFloat(priceFirstSplit.Gross+priceFirstSplit.Consultancy, 2)
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
