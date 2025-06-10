package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.dev.wopta.it/goworkspace/callback_out/base"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/payment/fabrick"

	"gitlab.dev.wopta.it/goworkspace/callback_out"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	tr "gitlab.dev.wopta.it/goworkspace/transaction"
)

const fabrickBillPaid string = "PAID"

func PaymentFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		responseFormat  string = `{"result":%t,"requestPayload":%s,"locale": "it"}`
		err             error
		fabrickCallback fabrick.FabrickCallback
	)

	log.AddPrefix("PaymentFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	origin = r.URL.Query().Get("origin")
	trSchedule = r.URL.Query().Get("schedule")

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal([]byte(request), &fabrickCallback)
	if err != nil {
		log.ErrorF("error unmarshaling request (%s): %s", string(request), err.Error())
		return fmt.Sprintf(responseFormat, false, string(request)), nil, nil
	}

	if fabrickCallback.PaymentID == nil {
		return "", nil, fmt.Errorf("no providerId found")
	}
	providerId = *fabrickCallback.PaymentID

	log.Printf("uid %s, providerId %s", policyUid, providerId)

	if policyUid == "" || origin == "" {
		ext := strings.Split(fabrickCallback.ExternalID, "_")
		policyUid = ext[0]
		trSchedule = ext[1]
		origin = ext[3]
	}

	policy := plc.GetPolicyByUid(policyUid, origin)

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		paymentMethod = strings.ToLower(*fabrickCallback.Bill.Transactions[0].PaymentMethod)
		err = fabrickPayment(origin, providerId, &policy)
	default:
	}

	if err != nil {
		log.ErrorF("error request (%s): %s", string(request), err.Error())
		return fmt.Sprintf(responseFormat, false, string(request)), nil, nil
	}

	callback_out.Execute(networkNode, policy, base.Paid)

	response := fmt.Sprintf(responseFormat, true, string(request))

	log.Println("Handler end -------------------------------------------------")

	return response, nil, nil
}

func fabrickPayment(origin, providerId string, policy *models.Policy) error {
	log.AddPrefix("fabrickPayment")
	defer log.PopPrefix()
	log.Printf("Policy %s", policy.Uid)

	policy.SanitizePaymentData()

	transaction, err := tr.GetTransactionToBePaid(policy.Uid, providerId, trSchedule, lib.TransactionsCollection)
	if err != nil {
		log.ErrorF("error getting transaction: %s", err.Error())
		return err
	}

	if transaction.IsPay {
		log.ErrorF("error Policy %s with transaction %s already paid", policy.Uid, transaction.Uid)
		return errors.New("transaction already paid")
	}

	state := runCallbackBpmn(policy, payFlowKey)
	if state == nil || state.Data == nil {
		log.Println("error bpmn - state not set")
		return nil
	}
	if state.IsFailed {
		log.Println("error bpmn - state failed")
		return nil
	}

	*policy = *state.Data

	return nil
}
