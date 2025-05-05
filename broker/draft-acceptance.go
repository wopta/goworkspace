package broker

import (
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	bpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	draftbpmn "github.com/wopta/goworkspace/broker/draftBpmn"
	"github.com/wopta/goworkspace/broker/draftBpmn/flow"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"

	plc "github.com/wopta/goworkspace/policy"
)

func (*AcceptancePayload) GetType() string {
	return "acceptanceInfo"
}

func DraftAcceptanceFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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

	origin := r.Header.Get("origin")
	policyUid := chi.URLParam(r, "policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)

	log.Printf("Policy Uid %s", policyUid)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = lib.CheckPayload[AcceptancePayload](body, &payload, []string{"action"})
	if err != nil {
		log.ErrorF("error: %s", err.Error())
		return "", nil, err
	}

	policy, err = plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.ErrorF("error retrieving policy %s from Firestore: %s", policyUid, err.Error())
		return "", nil, err
	}
	//TODO: to remove after test
	//	if !lib.SliceContains(models.GetWaitForApprovalStatusList(), policy.Status) {
	//		log.Printf("policy Uid %s: wrong status %s", policy.Uid, policy.Status)
	//		return "", nil, fmt.Errorf("policy uid '%s': wrong status '%s'", policy.Uid, policy.Status)
	//	}
	addresses := &flow.Addresses{}
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

	storage := draftbpmn.NewStorageBpnm()
	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	storage.AddGlobal("addresses", addresses)
	storage.AddGlobal("acceptanceInfo", &payload)
	flow, err := getFlow(&policy, networkNode, storage)
	if err != nil {
		return "", nil, err
	}
	err = flow.Run("acceptance")
	if err != nil {
		return "", nil, err
	}
	policy.Updated = time.Now().UTC()

	log.Println("saving to firestore...")
	err = lib.SetFirestoreErr(firePolicy, policy.Uid, &policy)
	if err != nil {
		log.ErrorF("error saving policy to firestore: %s", err.Error())
		return "", nil, err
	}
	log.Println("firestore saved!")

	policy.BigquerySave(origin)

	policyJsonLog, err := policy.Marshal()
	if err != nil {
		log.ErrorF("error marshaling policy: %s", err.Error())
	}
	log.Printf("Policy: %s", string(policyJsonLog))

	log.Println("Handler end -------------------------------------------------")

	return "{}", nil, nil
}

func draftrejectPolicy(storage draftbpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", storage)
	if err != nil {
		return err
	}
	acceptanceInfo, err := bpmn.GetData[*AcceptancePayload]("acceptanceInfo", storage)
	if err != nil {
		return err
	}
	policy.Status = models.PolicyStatusRejected
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = acceptanceInfo.Action
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s REJECTED", policy.Uid)
	return nil
}

func draftapprovePolicy(storage draftbpmn.StorageData) error {
	policy, err := bpmn.GetData[*flow.PolicyDraft]("policy", storage)
	if err != nil {
		return err
	}
	acceptanceInfo, err := bpmn.GetData[*AcceptancePayload]("acceptanceInfo", storage)
	if err != nil {
		return err
	}
	policy.Status = models.PolicyStatusApproved
	policy.StatusHistory = append(policy.StatusHistory, policy.Status)
	policy.ReservedInfo.AcceptanceNote = acceptanceInfo.Action
	policy.ReservedInfo.AcceptanceDate = time.Now().UTC()
	log.Printf("Policy Uid %s APPROVED", policy.Uid)
	return nil
}
