package payment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/payment/client"
	"gitlab.dev.wopta.it/goworkspace/payment/internal"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	prd "gitlab.dev.wopta.it/goworkspace/product"
	"gitlab.dev.wopta.it/goworkspace/transaction"
)

type changePaymentProviderReq struct {
	PolicyUid         string `json:"policyUid"`
	ProviderName      string `json:"providerName"`
	ScheduleFirstRate bool   `json:"scheduleFirstRate"`
}

type changePaymentProviderResp struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
}

/*
Now this function is used only to change payment provider to Fabrick (info hardcoded in frontend call) for those
policies that have been imported. When we will have multi providers we should delete transactions schedule from
old provider systems and only then send schedule new transactions to new provider systems.
*/
func changePaymentProviderFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err                  error
		payUrl               string
		policy               models.Policy
		updatedTransactions  []models.Transaction
		req                  changePaymentProviderReq
		resp                 changePaymentProviderResp
		responseTransactions = make([]models.Transaction, 0)
		unpaidTransactions   = make([]models.Transaction, 0)
	)

	log.AddPrefix("ChangePaymentProviderFx ")
	defer log.PopPrefix()
	log.Println("Handler Start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.ErrorF("error unmarshaling request body: %s", string(body))
		return "{}", nil, err
	}

	policy, err = plc.GetPolicy(req.PolicyUid)
	if err != nil {
		log.Printf("no policy found with uid %s: %s", req.PolicyUid, err.Error())
		return "{}", nil, err
	}

	policy.SanitizePaymentData()

	err = internal.UpdatePaymentProvider(&policy, req.ProviderName)
	if err != nil {
		log.Printf("provider update failed: %s", err.Error())
		return "{}", nil, err
	}

	activeTransactions := transaction.GetPolicyValidTransactions(policy.Uid, nil)
	for _, tr := range activeTransactions {
		if tr.IsPay {
			responseTransactions = append(responseTransactions, tr)
			continue
		}
		transaction.ReinitializePaymentInfo(&tr, policy.Payment)
		unpaidTransactions = append(unpaidTransactions, tr)
	}
	if len(unpaidTransactions) == 0 {
		log.Printf("no unpaid transactions found for policy %s", policy.Uid)
		return "{}", nil, errors.New("no unpaid transactions to update")
	}

	product := prd.GetProductV2(policy.Name, policy.ProductVersion, policy.Channel, nil, nil)

	client := client.NewClient(policy.Payment, policy, *product, unpaidTransactions, req.ScheduleFirstRate, "")
	payUrl, updatedTransactions, err = client.Update()
	if err != nil {
		log.ErrorF("error changing payment provider to %s: %s", req.ProviderName, err.Error())
		return "{}", nil, err
	}

	policy.PayUrl = payUrl
	responseTransactions = append(responseTransactions, updatedTransactions...)

	err = internal.SaveTransactionsToDB(updatedTransactions, lib.TransactionsCollection)
	if err != nil {
		return "{}", nil, err
	}

	err = lib.SetFirestoreErr(models.PolicyCollection, policy.Uid, policy)
	if err != nil {
		return "{}", nil, err
	}

	policy.BigquerySave()

	resp.Policy = policy
	resp.Transactions = responseTransactions
	rawResp, err := json.Marshal(resp)
	if err != nil {
		log.ErrorF("error marshaling response: %s", err.Error())
		return "{}", nil, err
	}

	log.Println("Handler End -------------------------------------------------")

	return string(rawResp), resp, err
}
