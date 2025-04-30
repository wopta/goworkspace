package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wopta/goworkspace/callback"
	"github.com/wopta/goworkspace/lib/log"

	draftbpnm "github.com/wopta/goworkspace/broker/draftBpnm"
	"github.com/wopta/goworkspace/broker/draftBpnm/flow"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	tr "github.com/wopta/goworkspace/transaction"
)

const fabrickBillPaid string = "PAID"

func PaymentDraftFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		responseFormat  string = `{"result":%t,"requestPayload":%s,"locale": "it"}`
		err             error
		fabrickCallback callback.FabrickCallback
	)

	log.AddPrefix("PaymentFx")
	defer log.PopPrefix()
	paymentInfo := flow.PaymentInfoBpmn{}
	log.Println("Handler start -----------------------------------------------")

	policyUid := r.URL.Query().Get("uid")
	origin = r.URL.Query().Get("origin")

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
	providerId := *fabrickCallback.PaymentID

	log.Printf("uid %s, providerId %s", policyUid, providerId)

	if policyUid == "" || origin == "" {
		ext := strings.Split(fabrickCallback.ExternalID, "_")
		policyUid = ext[0]
		paymentInfo.Schedule = ext[1]
		origin = ext[3]
	}

	policy := plc.GetPolicyByUid(policyUid, origin)

	switch fabrickCallback.Bill.Status {
	case fabrickBillPaid:
		paymentInfo.PaymentMethod = strings.ToLower(*fabrickCallback.Bill.Transactions[0].PaymentMethod)
		err = fabrickPayment(origin, providerId, &policy, paymentInfo)
	default:
	}

	if err != nil {
		log.ErrorF("error request (%s): %s", string(request), err.Error())
		return fmt.Sprintf(responseFormat, false, string(request)), nil, nil
	}

	//callback_out.Execute(networkNode, policy, callback_out.Paid)

	response := fmt.Sprintf(responseFormat, true, string(request))

	log.Println("Handler end -------------------------------------------------")

	return response, nil, nil
}

func fabrickPayment(origin, providerId string, policy *models.Policy, paymentInfo flow.PaymentInfoBpmn) error {
	log.AddPrefix("fabrickPayment")
	defer log.PopPrefix()
	log.Printf("Policy %s", policy.Uid)

	policy.SanitizePaymentData()

	transaction, err := tr.GetTransactionToBePaid(policy.Uid, providerId, paymentInfo.Schedule, lib.TransactionsCollection)
	if err != nil {
		log.ErrorF("error getting transaction: %s", err.Error())
		return err
	}

	if transaction.IsPay {
		log.ErrorF("error Policy %s with transaction %s already paid", policy.Uid, transaction.Uid)
		return errors.New("transaction already paid")
	}
	storage := draftbpnm.NewStorageBpnm()
	storage.AddGlobal("paymentInfo", paymentInfo)
	flowPayment, err := getFlow(policy, networkNode, storage)
	if err != nil {
		return err
	}
	err = flowPayment.Run("pay")
	if err != nil {
		return err
	}
	return nil
}
