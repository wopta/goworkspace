package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	draftbpnm "github.com/wopta/goworkspace/broker/draftBpnm"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	plc "github.com/wopta/goworkspace/policy"
	prd "github.com/wopta/goworkspace/product"
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

func DraftEmitFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
	//TODO: to remove after test
	//	if policy.Channel == models.NetworkChannel && policy.ProducerUid != authToken.UserID {
	//		log.Printf("user %s cannot emit policy %s because producer not equal to request user", authToken.UserID, policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

	policyJsonLog, _ := policy.Marshal()
	log.Printf("Policy %s JSON: %s", uid, string(policyJsonLog))
	//TODO: to remove after test
	//	if policy.IsPay || policy.IsSign || policy.CompanyEmit || policy.CompanyEmitted || policy.IsDeleted {
	//		log.Printf("cannot emit policy %s because state is not correct", policy.Uid)
	//		return "", nil, errors.New("operation not allowed")
	//	}

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
	responseEmit, err = emitDraft(&policy, request, origin)
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(responseEmit)

	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, err
}

func emitDraft(policy *models.Policy, request EmitRequest, origin string) (EmitResponse, error) {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	fireGuarantee := lib.GetDatasetByEnv(origin, lib.GuaranteeCollection)

	log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
	log.Println("[Emit] starting bpmn flow...")
	flow, err := getFlow(policy, nil, draftbpnm.NewStorageBpnm())
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
	log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))

	log.Println("[Emit] saving policy to firestore...")
	err = lib.SetFirestoreErr(firePolicy, request.Uid, policy)
	if err != nil {
		return responseEmit, err
	}

	log.Println("[Emit] saving policy to bigquery...")
	policy.BigquerySave(origin)

	log.Println("[Emit] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	log.Println("[Emit] end --------------------------------------------------")
	return responseEmit, nil
}
