package broker

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/bpmn"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine"
	"gitlab.dev.wopta.it/goworkspace/bpmn/bpmnEngine/flow"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/mail"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"

	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

func DraftAcceptanceFx(w http.ResponseWriter, r *http.Request) (string, any, error) {
	var (
		err     error
		payload AcceptancePayload
		policy  models.Policy
	)

	log.AddPrefix("AcceptanceFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	log.Println("loading authToken from idToken...")

	token := r.Header.Get("Authorization")
	authToken, err := lib.GetAuthTokenFromIdToken(token)
	if err != nil {
		log.ErrorF("error getting authToken: %s", err.Error())
		return "", nil, err
	}
	log.Printf(
		"authToken - type: '%s' role: '%s' uid: '%s' email: '%s'",
		authToken.Type,
		authToken.Role,
		authToken.UserID,
		authToken.Email,
	)

	policyUid := chi.URLParam(r, "policyUid")

	log.Printf("Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = lib.CheckPayload(body, &payload, []string{"action"})
	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	policy, err = plc.GetPolicy(policyUid)
	if err != nil {
		log.ErrorF("error retrieving policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}
	if !lib.SliceContains(models.GetWaitForApprovalStatusList(), policy.Status) {
		log.Printf("policy Uid %s: wrong status %s", policy.Uid, policy.Status)
		return "", nil, fmt.Errorf("policy uid '%s': wrong status '%s'", policy.Uid, policy.Status)
	}
	addresses := &flow.Addresses{}
	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	switch policy.Channel {
	case models.MgaChannel:
		addresses.ToAddress = mail.Address{
			Address: authToken.Email,
		}
	case models.NetworkChannel:
		addresses.ToAddress = mail.GetNetworkNodeEmail(networkNode)
	default:
		addresses.ToAddress = mail.GetContractorEmail(&policy)
	}

	storage := bpmnEngine.NewStorageBpnm()
	storage.AddGlobal("addresses", addresses)
	storage.AddGlobal("action", &flow.String{String: payload.Action})

	flow, err := bpmn.GetFlow(&policy, storage)
	if err != nil {
		return "", nil, err
	}
	err = flow.Run("acceptance")
	if err != nil {
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}
