package broker

import (
	"encoding/json"
	"fmt"
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

	log.SetPrefix("[EmitFx]")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
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

	log.Printf("Request: %s", string(body))
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

	policyJsonLog, _ := policy.Marshal()
	log.Printf("Policy %s JSON: %s", uid, string(policyJsonLog))

	if request.SendEmail == nil {
		sendEmail = true
	} else {
		sendEmail = *request.SendEmail
	}

	emitUpdatePolicy(&policy, request)

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if lib.GetBoolEnv("PROPOSAL_V2") {
		if policy.IsReserved && policy.Status != models.PolicyStatusApproved {
			log.Printf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
			return "", nil, fmt.Errorf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
		}
		responseEmit = emitV2(authToken, &policy, request, origin)
	} else {
		responseEmit = emit(authToken, &policy, request, origin)
	}
	b, e := json.Marshal(responseEmit)

	log.Println("Response: ", string(b))
	log.Println("Handler end ----------------------------------------")

	return string(b), responseEmit, e
}

func emit(authToken models.AuthToken, policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

	emitType := getEmitTypeFromPolicy(policy)
	switch emitType {
	case typeApprove:
		log.Printf("[Emit] Wait for approval - Policy Uid %s", policy.Uid)
		emitApproval(policy)
		mail.SendMailReserved(
			*policy,
			mail.AddressAnna,
			mail.GetContractorEmail(policy),
			mail.GetAgentEmail(policy),
			models.ProviderMgaFlow, // With PROPOSAL_V2 turned off, the only flow that should get here is the old agent
			[]string{models.ProposalAttachmentName},
		)
	case typeEmit:
		log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
		log.Println("[Emit] starting bpmn flow...")
		state := runBrokerBpmn(policy, emitFlowKey)
		if state == nil || state.Data == nil {
			log.Println("[Emit] error bpmn - state not set")
			return responseEmit
		}
		*policy = *state.Data
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

	log.Println("[Emit] end --------------------------------------------------")
	return responseEmit
}

func emitV2(authToken models.AuthToken, policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)

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
	if policy.PaymentSplit == string(models.PaySplitYear) {
		log.Println("[brokerUpdatePolicy] rectify paysplit year into yearly")
		policy.PaymentSplit = string(models.PaySplitYearly)
	}
	if policy.Payment == "" {
		log.Println("[brokerUpdatePolicy] inject policy payment provider from request")
		policy.Payment = request.Payment
	}
	if policy.Payment == "" || policy.Payment == "fabrik" {
		policy.Payment = models.FabrickPaymentProvider
	}
	if policy.PaymentMode == "" {
		log.Println("[brokerUpdatePolicy] inject policy payment mode from request")
		policy.PaymentMode = request.PaymentMode
	}
	if policy.PaymentMode == "" {
		if policy.PaymentSplit == string(models.PaySplitYearly) {
			log.Println("[brokerUpdatePolicy] inject policy payment mode from fallback")
			policy.PaymentMode = models.PaymentModeSingle
		}
		if policy.PaymentSplit == string(models.PaySplitMonthly) {
			log.Println("[brokerUpdatePolicy] inject policy payment mode from fallback")
			policy.PaymentMode = models.PaymentModeRecurrent
		}
	}
	log.Println("[brokerUpdatePolicy] end --------------------------------------")
}

func getEmitTypeFromPolicy(policy *models.Policy) string {
	if !policy.IsReserved || policy.Status == models.PolicyStatusApproved {
		return typeEmit
	}

	deniedStatuses := []string{models.PolicyStatusDeleted, models.PolicyStatusRejected}

	if policy.IsReserved && !lib.SliceContains(deniedStatuses, policy.Status) {
		return typeApprove
	}

	log.Printf("[getEmitTypeFromPolicy] error no type found for isReserved '%t' and status '%s'", policy.IsReserved, policy.Status)
	return ""
}

func emitApproval(policy *models.Policy) {
	log.Printf("[emitApproval] Policy Uid %s: Reserved Flow", policy.Uid)
	policy.Status = models.PolicyStatusWaitForApproval
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
}

func emitBase(policy *models.Policy, origin string) {
	log.Printf("[emitBase] Policy Uid %s", policy.Uid)
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)
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
	policy.PayUrl, _ = payment.PaymentController(origin, policy, product, mgaProduct)
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

	tr := transaction.PutByPolicy(*policy, "", origin, "", "", policy.PriceGross, policy.PriceNett, "", models.PayMethodRemittance, true, mgaProduct, policy.StartDate)

	transaction.CreateNetworkTransactions(policy, tr, networkNode, mgaProduct)
}
