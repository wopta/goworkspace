package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/lib"
	plc "github.com/wopta/goworkspace/policy"
	tr "github.com/wopta/goworkspace/transaction"
)

const fabrickBillPaid string = "PAID"

func PaymentV2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PaymentV2Fx] Handler start -----------------------------------")

	var (
		responseFormat  string = `{"result":%t,"requestPayload":%s,"locale": "it"}`
		err             error
		fabrickCallback FabrickCallback
	)

	policyUid := r.URL.Query().Get("uid")
	trSchedule = r.URL.Query().Get("schedule")
	origin = r.URL.Query().Get("origin")
	log.Printf("[PaymentV2Fx] uid %s, schedule %s", policyUid, trSchedule)

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
		trSchedule = ext[1]
		origin = ext[2]
	}

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		paymentMethod = strings.ToLower(*fabrickCallback.Bill.Transactions[0].PaymentMethod)
		err = fabrickPayment(origin, policyUid, trSchedule)
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

func fabrickPayment(origin, policyUid, trSchedule string) error {
	log.Printf("[fabrickPayment] Policy %s", policyUid)

	policy := plc.GetPolicyByUid(policyUid, origin)

	transaction, err := tr.GetTransactionByPolicyUidAndScheduleDate(policy.Uid, trSchedule, origin)
	if err != nil {
		log.Printf("[fabrickPayment] ERROR getting transaction: %s", err.Error())
		return err
	}

	if transaction.IsPay {
		log.Printf("[fabrickPayment] ERROR Policy %s with transaction %s already paid", policy.Uid, transaction.Uid)
		return errors.New("transaction already paid")
	}

	state := runCallbackBpmn(&policy, payFlowKey)
	if state == nil || state.Data == nil {
		log.Println("[fabrickPayment] error bpmn - state not set")
		return nil
	}
	if state.IsFailed {
		log.Println("[fabrickPayment] error bpmn - state failed")
		return nil
	}

	return nil
}
