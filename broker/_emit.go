package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/broker/internal/utility"
	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib/log"

	prd "gitlab.dev.wopta.it/goworkspace/product"

	"cloud.google.com/go/civil"
	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/question"
)

func EmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
		log.ErrorF("error getting authToken")
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
		log.ErrorF("error unmarshaling policy: %s", err.Error())
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
	//!!!!!TODO must be eliminated, should use either this or the new one
	//Only for test!!!!!
	if policy.Name == models.CatNatProduct {
		log.Println("Using emitCatnat")
		responseEmit, err = emitDraft(&policy, request, origin)
		if err != nil {
			return "", nil, err
		}

		b, err := json.Marshal(responseEmit)

		log.Println("Handler end -------------------------------------------------")

		return string(b), responseEmit, err
	}
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
	log.AddPrefix("Emit")
	defer log.PopPrefix()
	log.Println("start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.PolicyCollection
	fireGuarantee := lib.GuaranteeCollection

	log.Printf("Emitting - Policy Uid %s", policy.Uid)
	log.Println("starting bpmn flow...")
	state := runBrokerBpmn(policy, emitFlowKey)
	if state == nil || state.Data == nil || state.IsFailed {
		log.Println("error bpmn - state not set correctly")
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
	log.Printf("Policy %s: %s", request.Uid, string(policyJson))

	log.Println("saving policy to firestore...")
	err := lib.SetFirestoreErr(firePolicy, request.Uid, policy)
	lib.CheckError(err)

	log.Println("saving policy to bigquery...")
	policy.BigquerySave(origin)

	log.Println("saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	callbackAction := base.Emit
	if warrant != nil && warrant.GetFlowName(policy.Name) == models.RemittanceMgaFlow {
		callbackAction = base.EmitRemittance
	}

	callback_out.Execute(networkNode, *policy, callbackAction)

	log.Println("end --------------------------------------------------")
	return responseEmit
}
