package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	draftbpnm "gitlab.dev.wopta.it/goworkspace/broker/draftBpmn"
	"gitlab.dev.wopta.it/goworkspace/broker/draftBpmn/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
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

	policy, err = plc.GetPolicy(uid)
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
	responseEmit, err = emitDraft(&policy, request)
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(responseEmit)

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, err
}

func emitDraft(policy *models.Policy, request EmitRequest) (EmitResponse, error) {
	log.Println("start ------------------------------------------------")
	var responseEmit EmitResponse

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
	flow, err := getFlow(policy, storage)
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
	models.SetGuaranteBigquery(*policy, "emit", lib.GuaranteeCollection)

	log.Println("end --------------------------------------------------")
	return responseEmit, nil
}
