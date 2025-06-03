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
	prd "gitlab.dev.wopta.it/goworkspace/product"
)

// TO remove
func DraftEmitWithPolicyFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
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

	err = json.Unmarshal([]byte(body), &policy)
	if err != nil {
		log.Printf("error unmarshalling policy: %s", err.Error())
		return "", nil, err
	}

	productConfig := prd.GetProductV2(policy.Name, policy.ProductVersion, models.MgaChannel, nil, nil)
	if err = policy.CheckStartDateValidity(productConfig.EmitMaxElapsedDays); err != nil {
		return "", "", err
	}

	networkNode = network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	if policy.IsReserved && policy.Status != models.PolicyStatusApproved {
		log.Printf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
		return "", nil, fmt.Errorf("cannot emit policy uid %s with status %s and isReserved %t", policy.Uid, policy.Status, policy.IsReserved)
	}
	responseEmit, err = emitDraftWithPolicy(&policy, origin)
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(responseEmit)
	if err != nil {
		return string(b), responseEmit, err
	}
	log.Println("Handler end -------------------------------------------------")

	return string(b), responseEmit, err
}

func emitDraftWithPolicy(policy *models.Policy, origin string) (EmitResponse, error) {
	log.Println("[Emit] start ------------------------------------------------")
	var responseEmit EmitResponse

	fireGuarantee := lib.GetDatasetByEnv(origin, lib.GuaranteeCollection)

	log.Printf("[Emit] Emitting - Policy Uid %s", policy.Uid)
	log.Println("[Emit] starting bpmn flow...")

	paymentSplit = "monthly"
	log.Printf("paymentSplit: %s", paymentSplit)

	storage := draftbpnm.NewStorageBpnm()
	storage.AddGlobal("sendEmail", &flow.BoolBpmn{Bool: true})
	storage.AddGlobal("paymentSplit", &flow.String{String: "monthly"})
	storage.AddGlobal("paymentMode", &flow.String{String: "single"})
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
	log.Printf("[Emit] Policy %s: %s", policy.Uid, string(policyJson))

	log.Println("[Emit] saving guarantees to bigquery...")
	models.SetGuaranteBigquery(*policy, "emit", fireGuarantee)

	log.Println("[Emit] end --------------------------------------------------")
	return responseEmit, nil
}
