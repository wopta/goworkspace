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
}

type EmitRequest struct {
	Uid          string              `json:"uid,omitempty"`
	Payment      string              `json:"payment,omitempty"`
	PaymentType  string              `json:"paymentType,omitempty"`
	PaymentSplit string              `json:"paymentSplit,omitempty"`
	Statements   *[]models.Statement `json:"statements,omitempty"`
	SendEmail    *bool               `json:"sendEmail"`
}

func EmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		request      EmitRequest
		err          error
		policy       models.Policy
		responseEmit EmitResponse
	)

	log.Println("[EmitFx] Handler start --------------------------------------")

	log.Println("[EmitFx] loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := models.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.Printf("[EmitFx] error getting authToken")
		return "", nil, err
	}
	log.Printf(
		"[EmitFx] authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	origin = r.Header.Get("origin")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[EmitFx] Request: %s", string(body))
	err = json.Unmarshal([]byte(body), &request)
	if err != nil {
		log.Printf("[EmitFx] error unmarshaling policy: %s", err.Error())
		return "", nil, err
	}

	uid := request.Uid
	log.Printf("[EmitFx] Uid: %s", uid)

	policy, err = plc.GetPolicy(uid, origin)
	lib.CheckError(err)

	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFx] Policy %s JSON: %s", uid, string(policyJsonLog))

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
			log.Printf("[EmitFx] cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
			return "", nil, fmt.Errorf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
		}
		responseEmit = emitV2(authToken, &policy, request, origin)
	} else {
		responseEmit = emit(authToken, &policy, request, origin)
	}
	b, e := json.Marshal(responseEmit)

	models.CreateAuditLog(r, string(body))

	log.Println("[EmitFx] Response: ", string(b))
	log.Println("[EmitFx] Handler end ----------------------------------------")

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
	if state == nil || state.Data == nil {
		log.Println("[Emit] error bpmn - state not set")
		return responseEmit
	}
	*policy = *state.Data

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
	if policy.PaymentSplit == "" {
		policy.PaymentSplit = request.PaymentSplit
	}
	log.Println("[emitUpdatePolicy] end --------------------------------------")
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
	policy.PayUrl, _ = payment.PaymentController(origin, policy, product)
}

func setAdvance(policy *models.Policy, origin string) {
	policy.Payment = models.ManualPaymentProvider
	policy.IsPay = true
	policy.Status = models.PolicyStatusPay
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusToPay, models.PolicyStatusPay)
	policy.PaymentSplit = string(models.PaySplitSingleInstallment)

	tr := transaction.PutByPolicy(*policy, "", origin, "", "", policy.PriceGross, policy.PriceNett, "", models.PayMethodRemittance, true)

	transaction.CreateNetworkTransactions(policy, tr, networkNode, mgaProduct)
}
