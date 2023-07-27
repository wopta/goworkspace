package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wopta/goworkspace/bpmn"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/mail"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/transaction"
)

const fabrickBillPaid string = "PAID"

var (
	origin   string
	schedule string
)

func PaymentV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PaymentV2Fx] Handler start -----------------------------------")

	var (
		responseFormat  string = `{"result":%t,"requestPayload":%s,"locale": "it"}`
		err             error
		fabrickCallback FabrickCallback
	)

	policyUid := r.URL.Query().Get("uid")
	schedule = r.URL.Query().Get("schedule")
	origin = r.URL.Query().Get("origin")
	log.Printf("[PaymentV2Fx] uid %s, schedule %s", policyUid, schedule)

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	log.Printf("[PaymentV2Fx] request payload: %s", string(request))
	err = json.Unmarshal([]byte(request), &fabrickCallback)
	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR unmarshaling request (%s): %s", string(request), err.Error())
		return fmt.Sprintf(responseFormat, false, string(request)), nil, nil
	}

	if policyUid == "" || origin == "" {
		ext := strings.Split(fabrickCallback.ExternalID, "_")
		policyUid = ext[0]
		schedule = ext[1]
		origin = ext[2]
	}

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		err = fabrickPayment(origin, policyUid, schedule)
	default:
	}

	if err != nil {
		log.Printf("[PaymentV2Fx] ERROR request (%s): %s", string(request), err.Error())
		return fmt.Sprintf(responseFormat, false, string(request)), nil, nil
	}

	response := fmt.Sprintf(responseFormat, true, string(request))
	log.Printf("[PaymentV2Fx] response: %s", response)

	return response, nil, nil
}

func fabrickPayment(origin, policyUid, schedule string) error {
	log.Printf("[fabrickPayment] Policy %s", policyUid)

	policy := plc.GetPolicyByUid(policyUid, origin)

	// TODO check handling: on what state we expect the policy to be for agency flow?
	if policy.AgencyUid != "" {
		state := runAgencyFlow(policy, models.UserRoleAgency)

		if state.IsFailed {
			return errors.New("bpmn failed")
		}

		return nil
	} else if !policy.IsPay && policy.Status == models.PolicyStatusToPay {
		// promote documents from temp bucket to user and connect it to policy
		err := plc.SetUserIntoPolicyContractor(&policy, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR SetUserIntoPolicyContractor %s", err.Error())
			return err
		}

		// Add Policy contract
		err = plc.AddContract(&policy, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR AddContract %s", err.Error())
			return err
		}

		// Update Policy as paid
		err = plc.Pay(&policy, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR Policy Pay %s", err.Error())
			return err
		}

		// Get Transaction
		tr, err := transaction.GetPolicyFirstTransaction(policy.Uid, schedule, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR GetPolicyFirstTransaction %s", err.Error())
			return err
		}

		// Pay Transaction
		err = transaction.Pay(&tr, origin)
		if err != nil {
			log.Printf("[fabrickPayment] ERROR Transaction Pay %s", err.Error())
			return err
		}

		// Update agency if present
		err = models.UpdateAgencyPortfolio(&policy, origin)
		if err != nil && err.Error() != "agency not set" {
			log.Printf("[fabrickPayment] ERROR UpdateAgencyPortfolio %s", err.Error())
			return err
		}

		// Update agent if present
		err = models.UpdateAgentPortfolio(&policy, origin)
		if err != nil && err.Error() != "agent not set" {
			log.Printf("[fabrickPayment] ERROR UpdateAgentPortfolio %s", err.Error())
			return err
		}

		policy.BigquerySave(origin)
		tr.BigQuerySave(origin)

		// Send mail with the contract to the user
		mail.SendMailContract(policy, nil)

		return nil
	}

	log.Printf("[fabrickPayment] ERROR Policy %s with status %s and isPay %t cannot be paid", policyUid, policy.Status, policy.IsPay)
	return errors.New("cannot pay policy")
}

func runAgencyFlow(policy models.Policy, channel string) *bpmn.State {
	var (
		setting models.NodeSetting
		err     error
		state   *bpmn.State
	)

	settingByte, err := lib.GetFromGoogleStorage(os.Getenv("GOOGLE_STORAGE_BUCKET"), "products/"+channel+"/setting.json")
	if err != nil {
		log.Printf("[runAgencyFlow] ERROR loading setting: %s", err.Error())
	}
	err = json.Unmarshal(settingByte, &setting)
	if err != nil {
		log.Printf("[runAgencyFlow] ERROR unmarshaling setting: %s", err.Error())
	}

	state = bpmn.NewBpmn(policy)

	state.AddTaskHandler("payTransaction", payTransaction)
	state.AddTaskHandler("sendMailContract", sendMailContract)

	state.RunBpmn(setting.PayFlow)

	return state
}

func payTransaction(state *bpmn.State) error {
	log.Println("[payTransaction] Handler start ---")

	policy := state.Data
	tr, err := transaction.GetPolicyFirstTransaction(policy.Uid, schedule, origin)
	if err != nil {
		log.Printf("[payTransaction] ERROR GetPolicyFirstTransaction %s", err.Error())
		return err
	}
	err = transaction.Pay(&tr, origin)
	if err != nil {
		log.Printf("[payTransaction] ERROR Pay %s", err.Error())
		return err
	}
	tr.BigQuerySave(origin)

	return nil
}

func sendMailContract(state *bpmn.State) error {
	log.Println("[sendMailContract] Handler start ---")

	policy := state.Data
	mail.SendMailContract(policy, nil)

	return nil
}
